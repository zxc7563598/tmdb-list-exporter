package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zxc7563598/tmdb-list-exporter/internal/app"
	"github.com/zxc7563598/tmdb-list-exporter/internal/config"
)

var cfgFile string

var long = `tmdb-list-exporter 是一个基于 TMDB API 的片单导出工具。

支持：
  • 导出片单为格式化 JSON 文件
  • 下载封面图片
  • 多存储驱动( 本地 / 阿里云 OSS )`

var example = `  # 生成默认配置文件
  tmdb init
  # 使用默认配置文件运行
  tmdb --config=config.yaml
  # 指定其他配置文件
  tmdb --config=configs/prod.yaml
`

var rootCmd = cobra.Command{
	Use:     "tmdb",
	Short:   "导出 TMDB 片单及封面图片",
	Long:    long,
	Example: example,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}
		if err := cfg.Validate(); err != nil {
			return err
		}
		fmt.Printf("读取片单: %d \n", cfg.TMDB.ListID)
		if err := app.RunExport(cfg); err != nil {
			return err
		}
		return nil
	},
}

var configTemplate = `tmdb:
  access_token: # API 读访问令牌，可前往 https://www.themoviedb.org/settings/api 获取
  list_id: # 片单ID，可在片单网页链接中查看，例如 https://www.themoviedb.org/list/8634743-anime-series 的 ID 就是 8634743

output:
  file: tmdbList.json # 导出 JSON 文件名

storage:
  drivers: # 根据自己的需要进行配置，全部注释则不进行图片存储
    - type: local # 本地存储
      local_path: ./images # 本地存储路径
    # - type: alioss # 阿里云 OSS 存储
    #   alioss_region: cn-shanghai # region，详见阿里云 OSS 控制台
    #   alioss_bucket: xxxxxxx # bucket，详见阿里云 OSS 控制台
    #   alioss_path: xxxxxx # 上传文件夹路径，最前面不要带 /
    #   alioss_access_key_id: xxxxx # 阿里云 Access Key
    #   alioss_access_key_secret: xxxxx # 阿里云 Secret Key`

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "生成默认配置文件",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := os.WriteFile("config.yaml", []byte(configTemplate), 0644); err != nil {
			return fmt.Errorf("生成配置文件失败: %w", err)
		}
		fmt.Println("已生成 config.yaml, 请填写配置后重新运行。")
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "配置文件(默认为config.yaml)")
	_ = viper.BindPFlags(rootCmd.PersistentFlags())
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(initCmd)
}
