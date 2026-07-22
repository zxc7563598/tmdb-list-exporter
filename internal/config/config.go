package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// TMDBConfig 定义 TMDB 相关配置。
type TMDBConfig struct {
	// AccessToken 为 TMDB API 的读访问令牌（Read Access Token）。
	AccessToken string `mapstructure:"access_token"`
	// ListID 为需要导出的 TMDB 片单 ID。
	ListID int64 `mapstructure:"list_id"`
}

// DriverConfig 定义单个存储驱动的配置。
type DriverConfig struct {
	Type                 string `mapstructure:"type"`
	LocalPath            string `mapstructure:"local_path"`
	AliossRegion         string `mapstructure:"alioss_region"`
	AliossBucket         string `mapstructure:"alioss_bucket"`
	AliossPath           string `mapstructure:"alioss_path"`
	AliossAccessKeyID    string `mapstructure:"alioss_access_key_id"`
	AliossAccessKeySecret string `mapstructure:"alioss_access_key_secret"`
}

// StorageConfig 定义存储模块的整体配置。
type StorageConfig struct {
	Drivers []DriverConfig `mapstructure:"drivers"`
}

// OutputConfig 定义导出输出配置。
type OutputConfig struct {
	// File 为导出 JSON 文件名，默认为 "tmdbList.json"。
	File string `mapstructure:"file"`
}

// Config 为整个应用的顶层配置结构。
type Config struct {
	TMDB    TMDBConfig    `mapstructure:"tmdb"`
	Storage StorageConfig `mapstructure:"storage"`
	Output  OutputConfig  `mapstructure:"output"`
}

// DefaultOutputFile 为默认输出文件名。
const DefaultOutputFile = "tmdbList.json"

// Load 负责加载并解析配置文件，每次调用返回全新的 Config 实例。
//
// 参数：
//
//	configFile - 指定配置文件路径；为空时默认查找 ./config.yaml
func Load(configFile string) (*Config, error) {
	v := viper.New()

	if configFile == "" {
		v.AddConfigPath("./")
		v.SetConfigName("config")
		v.SetConfigType("yaml")
	} else {
		v.SetConfigFile(configFile)
	}

	v.SetEnvPrefix("TMDB")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err == nil {
		fmt.Println("[config] 使用配置文件:", v.ConfigFileUsed())
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// 设置默认值
	if cfg.Output.File == "" {
		cfg.Output.File = DefaultOutputFile
	}

	return &cfg, nil
}

// Validate 校验配置合法性。
func (c *Config) Validate() error {
	if c.TMDB.AccessToken == "" {
		return fmt.Errorf("[config] TMDB access token 不能为空")
	}
	if c.TMDB.ListID == 0 {
		return fmt.Errorf("[config] TMDB list id 不能为空")
	}

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
		default:
			return fmt.Errorf("[config] 未知的存储类型: %s", d.Type)
		}
	}
	return nil
}
