package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// LocalStorage 为本地文件系统存储实现。
//
// 负责：
//   - 将图片字节数据保存到本地磁盘
//   - 自动创建多级目录
//
// 存储结构：
//
//	BasePath/<id>/<filename>
//
// 示例：
//
//	./output/12345/poster.jpg
type LocalStorage struct {
	BasePath string
}

// NewLocalStorage 创建一个本地存储实例。
//
// 参数：
//
//	path - 本地存储根目录，例如 ./output
//
// 返回：
//
//	*LocalStorage
func NewLocalStorage(path string) *LocalStorage {
	return &LocalStorage{
		BasePath: path,
	}
}

// Save 将文件保存到本地磁盘。
//
// 参数：
//
//	ctx      - 上下文（当前实现未使用，用于接口统一和未来扩展）
//	id       - 资源 ID（用于构造分目录）
//	filename - 文件名，例如 poster.jpg
//	data     - 文件字节数据
//
// 行为：
//  1. 构造完整文件路径：BasePath/<id>/<filename>
//  2. 自动创建目录（若不存在）
//  3. 将数据写入文件
//
// 返回：
//
//	error - 创建目录或写入文件失败时返回错误
func (l *LocalStorage) Save(ctx context.Context, id int64, filename string, data []byte) error {
	// 拼接路径
	fullPath := filepath.Join(l.BasePath, strconv.FormatInt(id, 10), filename)
	// 创建文件夹
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return fmt.Errorf("本地存储文件夹创建失败: %w", err)
	}
	// 写入文件
	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return fmt.Errorf("本地存储文件写入失败: %w", err)
	}
	return nil
}
