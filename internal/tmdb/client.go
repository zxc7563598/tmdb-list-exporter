package tmdb

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

// Client 为 TMDB API 客户端。
//
// 负责：
//   - 构造 API 请求
//   - 处理认证信息
//   - 管理 HTTP 连接池
//   - 实现限流重试机制（同时处理网络错误和 429 响应）
type Client struct {
	baseURL string
	token   string
	http    *http.Client
}

// NewClient 创建一个 TMDB API 客户端。
//
// 参数：
//
//	AccessToken - TMDB Read Access Token
//
// 行为：
//   - 设置基础 API 地址
//   - 初始化带连接池的 http.Client
//   - 设置整体请求超时时间为 15 秒
//
// 连接优化策略：
//   - 最多保持 10 个空闲连接
//   - 空闲连接最多保留 25 秒
//
// 返回：
//
//	*Client
func NewClient(AccessToken string) *Client {
	return &Client{
		baseURL: "https://api.themoviedb.org/3",
		token:   AccessToken,
		http: &http.Client{
			Timeout: time.Second * 15,
			Transport: &http.Transport{
				MaxIdleConns:    10,
				IdleConnTimeout: 25 * time.Second,
			},
		},
	}
}

// buildQuery 构建 URL Query 参数字符串。
//
// 参数：
//
//	urlOptions - key/value 形式的参数映射
//
// 示例：
//
//	map[string]string{
//	    "page": "1",
//	    "language": "zh-CN",
//	}
//
// 结果：
//
//	page=1&language=zh-CN
func buildQuery(urlOptions map[string]string) string {
	if len(urlOptions) == 0 {
		return ""
	}
	values := make(url.Values, len(urlOptions))
	for k, v := range urlOptions {
		values.Set(k, v)
	}
	return values.Encode()
}

// doRequestWithRetry 执行 HTTP 请求并支持自动重试。
//
// 对以下情况执行指数退避重试：
//   - HTTP 429 (Too Many Requests)
//   - 网络瞬时错误（连接超时、连接重置等）
//   - 5xx 服务端错误
//
// 不重试：
//   - context 取消（返回 ctx.Err()）
//   - 4xx 客户端错误（429 除外）
//
// 退避策略：2^i 秒 + 随机抖动（0~500ms）
//
// 参数：
//
//	req      - 已构造完成的 HTTP 请求
//	maxRetry - 最大重试次数
//
// 返回：
//
//	*http.Response - 成功响应
//	error          - 超过最大重试次数仍失败
func (c *Client) doRequestWithRetry(req *http.Request, maxRetry int) (*http.Response, error) {
	var lastErr error
	for i := 0; i <= maxRetry; i++ {
		resp, err := c.http.Do(req)

		// 网络层错误（DNS 解析失败、连接超时、连接重置等）→ 重试
		if err != nil {
			// context 取消则不重试
			if req.Context().Err() != nil {
				return nil, err
			}
			lastErr = err
			if i < maxRetry {
				backoff := time.Duration(1<<i) * time.Second
				jitter := time.Duration(rand.Intn(500)) * time.Millisecond
				time.Sleep(backoff + jitter)
			}
			continue
		}

		// 429 限流 → 重试
		if resp.StatusCode == http.StatusTooManyRequests {
			resp.Body.Close()
			lastErr = fmt.Errorf("HTTP 429: Too Many Requests")
			backoff := time.Duration(1<<i) * time.Second
			jitter := time.Duration(rand.Intn(500)) * time.Millisecond
			time.Sleep(backoff + jitter)
			continue
		}

		// 5xx 服务端错误 → 重试
		if resp.StatusCode >= 500 {
			resp.Body.Close()
			lastErr = fmt.Errorf("HTTP %d: server error", resp.StatusCode)
			if i < maxRetry {
				backoff := time.Duration(1<<i) * time.Second
				jitter := time.Duration(rand.Intn(500)) * time.Millisecond
				time.Sleep(backoff + jitter)
			}
			continue
		}

		// 2xx / 3xx / 4xx（非 429）→ 直接返回
		return resp, nil
	}
	return nil, fmt.Errorf("请求失败，达到最大重试次数: %w", lastErr)
}

// ImageBaseURL 为 TMDB 原始图片 CDN 地址。
const ImageBaseURL = "https://image.tmdb.org/t/p/original"
