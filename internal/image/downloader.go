package image

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/zxc7563598/tmdb-list-exporter/internal/tmdb"
)

// HTTPDownloader 基于 net/http 实现的图片下载器。
type HTTPDownloader struct {
	client *http.Client
}

// NewHTTPDownloader 创建一个 HTTPDownloader 实例。
//
// 调用示例：
//
//	downloader := NewHTTPDownloader(10 * time.Second)
func NewHTTPDownloader(timeout time.Duration) *HTTPDownloader {
	return &HTTPDownloader{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Download 根据图片相对路径下载图片数据。
//
// 会校验响应的 Content-Type 是否为 image/* 类型。
func (d *HTTPDownloader) Download(ctx context.Context, path string) ([]byte, error) {
	url := tmdb.ImageBaseURL + path

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("响应异常 Code: %d", resp.StatusCode)
	}

	// 校验 Content-Type，确保下载的是图片
	contentType := resp.Header.Get("Content-Type")
	if contentType != "" && !strings.HasPrefix(contentType, "image/") {
		return nil, fmt.Errorf("非图片 Content-Type: %s", contentType)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
