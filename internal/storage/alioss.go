package storage

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

// AliossStorage 为阿里云 OSS 存储实现。
//
// 负责：
//   - 初始化 OSS 客户端
//   - 将图片数据上传至指定 Bucket
//
// 上传路径结构：
//
//	BasePath/<id>/<filename>
//
// 示例：
//
//	images/12345/poster.jpg
type AliossStorage struct {
	Client   *oss.Client
	Bucket   string
	BasePath string
}

// NewAliossStorage 创建一个 AliossStorage 实例。
//
// 参数：
//
//	accessKeyID     - 阿里云 AccessKey ID
//	accessKeySecret - 阿里云 AccessKey Secret
//	region          - OSS 区域，例如 cn-shanghai
//	bucket          - bucket 名称
//	path            - 上传根路径（前缀目录）
//
// 行为：
//   - 使用静态凭证构造 CredentialsProvider
//   - 初始化 OSS Client
//   - 返回封装后的存储实例
//
// 注意：
//
//	当前实现不依赖环境变量认证，完全使用显式 AK/SK。
func NewAliossStorage(accessKeyID string, accessKeySecret string, region string, bucket string, path string) *AliossStorage {
	// 创建静态凭证提供者（不使用环境变量）
	provider := credentials.NewStaticCredentialsProvider(
		accessKeyID,
		accessKeySecret,
	)
	// 加载默认配置并设置凭证提供者和区域
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(provider).
		WithRegion(region)
	client := oss.NewClient(cfg)
	return &AliossStorage{
		Client:   client,
		Bucket:   bucket,
		BasePath: path,
	}
}

// Save 将图片数据上传到阿里云 OSS。
//
// 参数：
//
//	ctx      - 上下文，用于控制取消或超时
//	id       - 资源 ID（用于构造分目录）
//	filename - 文件名，例如 poster.jpg
//	data     - 文件字节数据
//
// 行为：
//  1. 构造上传路径：BasePath/<id>/<filename>
//  2. 发起 PutObject 请求
//
// 返回：
//
//	error - 上传失败时返回错误
func (a *AliossStorage) Save(ctx context.Context, id int64, filename string, data []byte) error {
	// 构造请求
	request := &oss.PutObjectRequest{
		Bucket: oss.Ptr(a.Bucket),
		Key:    oss.Ptr(filepath.Join(a.BasePath, strconv.FormatInt(id, 10), filename)),
		Body:   bytes.NewReader(data),
	}
	// 执行上传
	_, err := a.Client.PutObject(ctx, request)
	if err != nil {
		return fmt.Errorf("阿里云OSS存储失败: %w", err)
	}
	return nil
}
