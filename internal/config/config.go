package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// TMDBConfig 定义 TMDB 相关配置。
//
// 用于存储访问 TMDB API 所需的认证信息和目标片单 ID。
type TMDBConfig struct {
	// AccessToken 为 TMDB API 的读访问令牌（Read Access Token）。
	// 可在 https://www.themoviedb.org/settings/api 中获取。
	AccessToken string `mapstructure:"access_token"`
	// ListID 为需要导出的 TMDB 片单 ID。
	// 可在片单页面 URL 中获取，规则为 `/list/ID-名称` 例如：
	// https://www.themoviedb.org/list/8634743-anime-series
	// ID为：8634743
	ListID int64 `mapstructure:"list_id"`
}

// DriverConfig 定义单个存储驱动的配置。
//
// 通过 Type 区分不同的存储类型，例如：
//   - local  本地文件系统
//   - alioss 阿里云 OSS
//
// 不同类型会使用不同的字段。
type DriverConfig struct {
	// Type 表示存储驱动类型
	Type string `mapstructure:"type"`
	// ===== local 类型配置 =====
	// LocalPath 表示本地文件保存路径
	LocalPath string `mapstructure:"local_path"`
	// ===== alioss 类型配置 =====
	// AliossRegion 为阿里云 OSS 区域，例如 cn-shanghai
	AliossRegion string `mapstructure:"alioss_region"`
	// AliossBucket 为 OSS Bucket 名称
	AliossBucket string `mapstructure:"alioss_bucket"`
	// AliossPath 为上传文件夹路径（不需要以 / 开头）
	AliossPath string `mapstructure:"alioss_path"`
	// AliossAccessKeyID 为阿里云 AccessKey ID
	AliossAccessKeyID string `mapstructure:"alioss_access_key_id"`
	// AliossAccessKeySecret 为阿里云 AccessKey Secret
	AliossAccessKeySecret string `mapstructure:"alioss_access_key_secret"`
}

// StorageConfig 定义存储模块的整体配置。
//
// 支持多个驱动同时启用，例如：
//   - 同时保存到本地
//   - 同时上传到 OSS
//
// Drivers 为空时表示不启用图片存储。
type StorageConfig struct {
	Drivers []DriverConfig `mapstructure:"drivers"`
}

// Config 为整个应用的顶层配置结构。
//
// 包含：
//   - TMDB API 配置
//   - 存储配置
type Config struct {
	TMDB    TMDBConfig    `mapstructure:"tmdb"`
	Storage StorageConfig `mapstructure:"storage"`
}

// 全局配置实例（由 Load 方法初始化）
var cfg Config

// Load 负责加载并解析配置文件。
//
// 加载流程：
//  1. 读取指定配置文件（或默认 config.yaml）
//  2. 支持环境变量覆盖
//  3. 反序列化到 Config 结构体
//
// 参数：
//
//	configFile - 指定配置文件路径；为空时默认查找 ./config.yaml
//
// 返回：
//
//	*Config - 解析后的配置对象
//	error   - 加载或解析失败时返回错误
func Load(configFile string) (*Config, error) {
	// 如果存在配置文件
	if configFile == "" {
		// 默认查找当前目录
		viper.AddConfigPath("./")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	} else {
		viper.SetConfigFile(configFile)
	}
	// 环境变量支持
	viper.SetEnvPrefix("TMDB")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	// 如果文件存在，则读取
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("[config] 使用配置文件:", viper.ConfigFileUsed())
	}
	// 绑定到结构体
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Validate 校验配置合法性。
//
// 主要检查：
//   - TMDB 必填字段是否存在
//   - 各存储驱动所需字段是否完整
//
// 若存在配置错误，将返回明确的错误信息。
func (c *Config) Validate() error {
	// 验证 TMDB 配置
	if c.TMDB.AccessToken == "" {
		return fmt.Errorf("[config] TMDB access token 不能为空")
	}
	if c.TMDB.ListID == 0 {
		return fmt.Errorf("[config] TMDB list id 不能为空")
	}
	// 验证存储配置
	for _, d := range c.Storage.Drivers {
		switch d.Type {
		case "local":
			if d.LocalPath == "" {
				return fmt.Errorf("[config] local 启用时 LocalPath 不能为空")
			}
		case "alioss":
			if d.AliossRegion == "" {
				return fmt.Errorf("[config] alioss 启用时 AliossRegion 不能为空")
			}
			if d.AliossBucket == "" {
				return fmt.Errorf("[config] alioss 启用时 AliossBucket 不能为空")
			}
			if d.AliossPath == "" {
				return fmt.Errorf("[config] alioss 启用时 AliossPath 不能为空")
			}
			if d.AliossAccessKeyID == "" {
				return fmt.Errorf("[config] alioss 启用时 AliossAccessKeyID 不能为空")
			}
			if d.AliossAccessKeySecret == "" {
				return fmt.Errorf("[config] alioss 启用时 AliossAccessKeySecret 不能为空")
			}
		}
	}
	return nil
}
