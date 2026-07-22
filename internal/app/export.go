package app

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/zxc7563598/tmdb-list-exporter/internal/config"
	"github.com/zxc7563598/tmdb-list-exporter/internal/exporter"
	"github.com/zxc7563598/tmdb-list-exporter/internal/image"
	"github.com/zxc7563598/tmdb-list-exporter/internal/storage"
	"github.com/zxc7563598/tmdb-list-exporter/internal/tmdb"
)

// imageConcurrency 为图片下载的最大并发数。
// TMDB 图片 CDN 限流非常宽松，8 路并发可以安全使用。
const imageConcurrency = 8

// RunExport 为整个导出流程的编排入口。
//
// 流程：
//  1. 调用 tmdb 包获取片单
//  2. 根据配置并发下载海报图和背景图
//  3. 根据配置调用 storage 存储图片
//  4. 调用 exporter 包输出为 JSON 文件
func RunExport(cfg *config.Config) error {
	ctx := context.Background()

	// 创建 TMDB client
	client := tmdb.NewClient(cfg.TMDB.AccessToken)

	// 获取片单
	fmt.Println("--------------------")
	fmt.Println("开始获取片单")
	fmt.Println("--------------------")
	urlOptions := map[string]string{
		"language": "zh-CN",
	}
	list, err := client.GetAllListItems(ctx, cfg.TMDB.ListID, urlOptions)
	if err != nil {
		return fmt.Errorf("[tmdb] 片单列表获取失败: %w", err)
	}

	// 存储图片数据
	if len(cfg.Storage.Drivers) > 0 {
		fmt.Println("--------------------")
		fmt.Println("下载/存储图片信息")
		fmt.Println("--------------------")

		store, err := storage.NewFromConfig(cfg.Storage)
		if err != nil {
			return fmt.Errorf("[storage] 存储创建失败: %w", err)
		}

		imageClient := image.NewHTTPDownloader(10 * time.Second)

		if err := downloadImages(ctx, list, imageClient, store); err != nil {
			return err
		}
	}

	// 导出 JSON 数据
	fmt.Println("--------------------")
	fmt.Println("导出 JSON 数据")
	fmt.Println("--------------------")
	exportItems := exporter.FromTMDB(list)
	jsonExporter := exporter.NewJSONExporter("  ")
	if err := jsonExporter.Export(cfg.Output.File, exportItems); err != nil {
		return fmt.Errorf("[exporter] 数据导出失败: %w", err)
	}

	fmt.Println("")
	fmt.Println("")
	fmt.Printf("导出完成，已导出为: ./%s\n", cfg.Output.File)
	return nil
}

// downloadResult 记录单张图片的下载存储结果。
type downloadResult struct {
	title string
	kind  string // "海报图" or "背景图"
	err   error
}

// downloadImages 并发下载列表中所有条目的海报图与背景图。
//
// 使用有界并发（imageConcurrency 路）避免资源耗尽，
// 同时保持对 TMDB 图片 CDN 的友好访问。
func downloadImages(ctx context.Context, items []tmdb.ListItem, imageClient *image.HTTPDownloader, store storage.ImageStorage) error {
	var wg sync.WaitGroup
	sem := make(chan struct{}, imageConcurrency)
	results := make(chan downloadResult, len(items)*2)

	for _, item := range items {
		item := item // capture

		title := item.GetDisplayTitle()

		// 海报图
		if item.PosterPath != nil && *item.PosterPath != "" {
			wg.Add(1)
			go func() {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()
				results <- downloadAndSave(ctx, imageClient, store, item.ID, *item.PosterPath, title, "海报图")
			}()
		}

		// 背景图
		if item.BackdropPath != nil && *item.BackdropPath != "" {
			wg.Add(1)
			go func() {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()
				results <- downloadAndSave(ctx, imageClient, store, item.ID, *item.BackdropPath, title, "背景图")
			}()
		}
	}

	// 等待所有任务完成后关闭 results channel
	go func() {
		wg.Wait()
		close(results)
	}()

	// 收集结果并打印
	var firstErr error
	for r := range results {
		if r.err != nil {
			fmt.Fprintf(os.Stderr, "[%s] %s - %s失败: %v\n", r.kind, r.title, r.kind, r.err)
			if firstErr == nil {
				firstErr = r.err
			}
		} else {
			fmt.Printf("[storage] %s - %s存储成功\n", r.title, r.kind)
		}
	}

	// 只要有一个失败就返回错误（非致命，用户可以检查输出）
	if firstErr != nil {
		return fmt.Errorf("部分图片下载/存储失败，请检查上方输出")
	}
	return nil
}

// downloadAndSave 下载单张图片并调用存储驱动保存。
func downloadAndSave(ctx context.Context, imageClient *image.HTTPDownloader, store storage.ImageStorage, id int64, path string, title string, kind string) downloadResult {
	imageData, err := imageClient.Download(ctx, path)
	if err != nil {
		return downloadResult{title: title, kind: kind, err: fmt.Errorf("图片下载失败: %w", err)}
	}
	if err := store.Save(ctx, id, path, imageData); err != nil {
		return downloadResult{title: title, kind: kind, err: fmt.Errorf("图片存储失败: %w", err)}
	}
	return downloadResult{title: title, kind: kind}
}
