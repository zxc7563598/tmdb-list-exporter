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
//   - 实现限流重试机制
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
			Timeout: time.Second * 15, // 整个请求最多15秒
			Transport: &http.Transport{
				MaxIdleConns:    10,               // 最多保持10个空闲连接
				IdleConnTimeout: 25 * time.Second, // 空闲连接最多保持25秒
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
// 行为：
//   - 将 map 转换为标准 URL 编码格式
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
	values := url.Values{}
	for k, v := range urlOptions {
		values.Set(k, v)
	}
	return values.Encode()
}

// doRequestWithRetry 执行 HTTP 请求并支持自动重试。
//
// 参数：
//
//	req      - 已构造完成的 HTTP 请求
//	maxRetry - 最大重试次数
//
// 行为：
//   - 发起请求
//   - 若返回 429 (Too Many Requests)，执行指数退避重试
//   - 退避策略：2^i 秒 + 随机抖动
//
// 重试策略说明：
//
//	第 0 次重试：1 秒 + 随机抖动
//	第 1 次重试：2 秒 + 随机抖动
//	第 2 次重试：4 秒 + 随机抖动
//	...
//
// 抖动（jitter）用于避免多个客户端同时重试造成“雪崩”。
//
// 返回：
//
//	*http.Response - 成功响应
//	error          - 请求失败或超过最大重试次数
func (c *Client) doRequestWithRetry(req *http.Request, maxRetry int) (*http.Response, error) {
	for i := 0; i <= maxRetry; i++ {
		resp, err := c.http.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusTooManyRequests {
			return resp, nil
		}
		resp.Body.Close()
		backoff := time.Duration(1<<i) * time.Second
		jitter := time.Duration(rand.Intn(500)) * time.Millisecond
		time.Sleep(backoff + jitter)
	}
	return nil, fmt.Errorf("请求失败，达到最大重试次数")
}
