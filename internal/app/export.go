package app

import (
	"context"
	"fmt"
	"os"

	"github.com/zxc7563598/tmdb-list-exporter/internal/config"
	"github.com/zxc7563598/tmdb-list-exporter/internal/exporter"
	"github.com/zxc7563598/tmdb-list-exporter/internal/image"
	"github.com/zxc7563598/tmdb-list-exporter/internal/storage"
	"github.com/zxc7563598/tmdb-list-exporter/internal/tmdb"
)

// 业务逻辑
// ...
// 1. 调用 tmdb 包获取片单
// 2. 调用 image 包获取图片 []byte
// 3. 根据配置调用 storage 存储图片
// 4. 调用 exporter 包存储为 JSON 文件
// ...
func RunExport(cfg *config.Config) error {
	ctx := context.Background()
	// 创建 tmbd client
	client := tmdb.NewClient(cfg.TMDB.AccessToken)
	// 获取片单
	fmt.Println("--------------------")
	fmt.Println("开始获取片单")
	fmt.Println("--------------------")
	urlOptions := map[string]string{
		"language": "zh-CN",
	}
	list, err := client.GetAllListItems(cfg.TMDB.ListID, urlOptions)
	if err != nil {
		fmt.Printf("[tmdb] 片单列表获取失败: %v\n", err)
		os.Exit(1)
	}
	// 存储图片数据
	if len(cfg.Storage.Drivers) > 0 {
		fmt.Println("--------------------")
		fmt.Println("下载/存储图片信息")
		fmt.Println("--------------------")
		// 创建 storage
		store, err := storage.NewFromConfig(cfg.Storage)
		if err != nil {
			fmt.Printf("[storage] 存储创建失败: %v\n", err)
			os.Exit(1)
		}
		// 下载图片信息
		imageClient := image.NewHTTPDownloader(10)
		for _, item := range list {
			// 获取标题
			title := ""
			if item.Title != nil {
				title = *item.Title
			} else if item.Name != nil {
				title = *item.Name
			}
			// 存储海报图
			if item.PosterPath != nil && *item.PosterPath != "" {
				imageData, err := imageClient.Download(ctx, *item.PosterPath)
				if err == nil {
					if err := store.Save(ctx, item.ID, *item.PosterPath, imageData); err != nil {
						fmt.Printf("[storage] %s - 海报图存储失败: %v\n", title, err)
					} else {
						fmt.Printf("[storage] %s - 海报图存储成功\n", title)
					}
				} else {
					fmt.Printf("[image] %s - 图片下载失败: %v\n", title, err)
				}
			}
			// 存储背景图
			if item.BackdropPath != nil && *item.BackdropPath != "" {
				imageData, err := imageClient.Download(ctx, *item.BackdropPath)
				if err == nil {
					if err := store.Save(ctx, item.ID, *item.BackdropPath, imageData); err != nil {
						fmt.Printf("[storage] %s - 背景图图存储失败: %v\n", title, err)
					} else {
						fmt.Printf("[storage] %s - 海报图存储成功\n", title)
					}
				} else {
					fmt.Printf("[image] %s - 图片下载失败: %v\n", title, err)
				}
			}
		}
	}
	// 导出 JSON 数据
	fmt.Println("--------------------")
	fmt.Println("导出 JSON 数据")
	fmt.Println("--------------------")
	exportItems := exporter.FromTMDB(list)
	jsonExporter := exporter.NewJSONExporter("  ")
	errs := jsonExporter.Export("tmdbList.json", exportItems)
	if errs != nil {
		fmt.Printf("[exporter] 数据导出失败: %v\n", errs)
		os.Exit(1)
	}
	fmt.Println("")
	fmt.Println("")
	fmt.Println("导出完成，已导出为: ./tmdbList.json")
	return nil
}
