package storage

import (
	"context"
	"fmt"

	"github.com/zxc7563598/tmdb-list-exporter/internal/config"
)

// ImageStorage 定义图片存储接口。
//
// 该接口用于抽象不同的图片存储实现，例如：
//   - 本地文件系统
//   - 阿里云 OSS
//   - 未来可能支持的 S3、七牛云等
//
// Save 方法负责将图片数据持久化到对应的存储介质。
// id 通常用于标识图片所属资源（例如 TMDB 作品 ID）。
// filename 为保存时使用的文件名。
// data 为图片的原始二进制内容。
type ImageStorage interface {
	Save(ctx context.Context, id int64, filename string, data []byte) error
}

// MultiStorage 用于组合多个 ImageStorage，实现多存储同时写入。
//
// 例如可以同时：
//   - 保存到本地
//   - 上传到 OSS
//
// 只要任意一个存储返回错误，整个 Save 操作就会返回错误。
// 当前实现是“全部必须成功”策略。
type MultiStorage struct {
	stores []ImageStorage
}

// NewMultiStorage 创建一个多存储实例。
//
// 传入多个 ImageStorage 实现，返回一个统一对外的 ImageStorage。
// 调用 Save 时会按顺序依次执行所有存储。
func NewMultiStorage(stores []ImageStorage) *MultiStorage {
	return &MultiStorage{stores: stores}
}

// Save 将图片数据保存到所有注册的存储驱动中。
//
// 执行逻辑：
//   - 按顺序调用每一个存储驱动的 Save 方法
//   - 如果任意一个失败，则立即返回错误
//   - 全部成功才返回 nil
func (m *MultiStorage) Save(ctx context.Context, id int64, filename string, data []byte) error {
	for _, s := range m.stores {
		if err := s.Save(ctx, id, filename, data); err != nil {
			return err
		}
	}
	return nil
}

// NewFromConfig 根据配置创建对应的 ImageStorage 实例。
//
// 支持的存储类型包括：
//   - local  : 本地文件存储
//   - alioss : 阿里云 OSS 存储
//
// 配置中允许定义多个存储驱动：
//   - 如果只配置一个，则直接返回该存储实例
//   - 如果配置多个，则自动组合为 MultiStorage
//
// 如果出现未知类型，将返回错误。
func NewFromConfig(cfg config.StorageConfig) (ImageStorage, error) {
	var stores []ImageStorage
	for _, d := range cfg.Drivers {
		switch d.Type {
		case "local":
			stores = append(stores, NewLocalStorage(d.LocalPath))
		case "alioss":
			stores = append(stores, NewAliossStorage(
				d.AliossAccessKeyID,
				d.AliossAccessKeySecret,
				d.AliossRegion,
				d.AliossBucket,
				d.AliossPath,
			))
		default:
			return nil, fmt.Errorf("未知的存储类型: %s", d.Type)
		}
	}
	// 如果只配置了一个存储驱动，直接返回该驱动
	if len(stores) == 1 {
		return stores[0], nil
	}
	// 多个驱动则返回组合存储
	return NewMultiStorage(stores), nil
}
