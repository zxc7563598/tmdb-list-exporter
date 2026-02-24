package image

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPDownloader 基于 net/http 实现的图片下载器。
//
// 主要职责：
//   - 根据 TMDB 图片路径下载原始图片
//   - 支持超时控制
//   - 支持 context 取消
type HTTPDownloader struct {
	client       *http.Client
	imageBaseURL string
}

// NewHTTPDownloader 创建一个 HTTPDownloader 实例。
//
// 参数：
//
//	timeout - 超时时间（单位：秒）
//
// 行为：
//   - 构造带超时控制的 http.Client
//
// 示例：
//
//	downloader := NewHTTPDownloader(10)
func NewHTTPDownloader(timeout time.Duration) *HTTPDownloader {
	return &HTTPDownloader{
		client: &http.Client{
			Timeout: time.Second * timeout,
		},
		imageBaseURL: "https://image.tmdb.org/t/p/original",
	}
}

// Download 根据图片相对路径下载图片数据。
//
// 参数：
//
//	ctx  - 上下文，用于控制取消或超时
//	path - TMDB 图片相对路径，例如：/abc.jpg
//
// 行为：
//  1. 拼接 TMDB 原始图片地址
//  2. 发起 GET 请求
//  3. 校验 HTTP 状态码
//  4. 读取并返回图片字节数据
//
// 返回：
//
//	[]byte - 图片原始字节
//	error  - 请求或读取失败时返回错误
func (d *HTTPDownloader) Download(ctx context.Context, path string) ([]byte, error) {
	// 拼接完整图片 URL
	url := d.imageBaseURL + path
	// 构造带上下文的 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	// 发起请求
	resp, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// 校验响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("响应异常 Code : %d", resp.StatusCode)
	}
	// 读取响应体
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
