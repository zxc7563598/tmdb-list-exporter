# tmdb-list-exporter

[English](./README.md) ｜ 简体中文

基于 TMDB API 的片单导出命令行工具。自动拉取指定片单的全部数据，导出为格式化 JSON，同时下载封面图和背景图，支持本地文件系统和阿里云 OSS 等多种存储方式。

> 本项目已由 [Zread](https://zread.ai/zxc7563598/tmdb-list-exporter) 解析，可点击链接快速了解项目结构和代码逻辑。

## 使用场景

- 搭建个人影视墙或番剧展示页
- 为静态站点提供影视数据
- NAS 本地影视展示与数据归档
- 需要将 TMDB 图片转存到自有存储的场合

## 特性

- 自动分页拉取，3 路并发加速多页请求
- 8 路并发下载封面图和背景图
- 支持多存储驱动：本地文件系统 / 阿里云 OSS，可同时启用
- 指数退避重试机制，自动处理 HTTP 429 限流和网络抖动
- 支持通过环境变量覆盖配置（`TMDB_` 前缀）
- 5 平台交叉编译（macOS / Linux / Windows，amd64 + arm64）

## 技术栈

- 语言：Go 1.25
- CLI 框架：[Cobra](https://github.com/spf13/cobra) + [Viper](https://github.com/spf13/viper)
- OSS SDK：[alibabacloud-oss-go-sdk-v2](https://github.com/aliyun/alibabacloud-oss-go-sdk-v2)

## 安装

前往 [GitHub Releases](https://github.com/zxc7563598/tmdb-list-exporter/releases) 下载对应平台的二进制文件：

- macOS（amd64 / arm64）
- Linux（amd64 / arm64）
- Windows（amd64）

下载后赋予执行权限并运行：

```bash
chmod +x tmdb
./tmdb --help
```

也可以从源码编译，见[本地开发](#本地开发)。

## 快速开始

### 1. 获取 TMDB Access Token

1. 注册 [TMDB](https://www.themoviedb.org/) 账号
2. 进入 [API 设置页面](https://www.themoviedb.org/settings/api)
3. 申请 API Key，获取 **API Read Access Token**（不是 API Key）

### 2. 生成配置文件

```bash
tmdb init
```

命令会在当前目录生成 `config.yaml`，按注释填写配置即可。

### 3. 运行

```bash
tmdb --config=config.yaml
```

## 配置说明

完整配置示例：

```yaml
tmdb:
  access_token: "你的 TMDB Read Access Token"
  list_id: 8634743  # 片单 ID，从片单 URL 中获取

output:
  file: tmdbList.json  # 导出的 JSON 文件名

storage:
  drivers:  # 留空则不存储图片
    - type: local
      local_path: ./images
    - type: alioss
      alioss_region: cn-shanghai
      alioss_bucket: your-bucket
      alioss_path: images
      alioss_access_key_id: your-access-key
      alioss_access_key_secret: your-secret-key
```

> [!NOTE]
> `alioss_path` 为上传路径，**不要**以 `/` 开头。

### 片单 ID 如何获取

打开 TMDB 片单页面，URL 中的数字即为 list_id。例如 `https://www.themoviedb.org/list/8634743-anime-series` 中的 `8634743`。

### 环境变量

所有配置项支持通过环境变量覆盖，前缀为 `TMDB_`，使用下划线分隔层级：

```bash
export TMDB_ACCESS_TOKEN="your-token"
export TMDB_LIST_ID="8634743"
```

环境变量的优先级高于配置文件。

## 输出结构

运行后会在当前目录生成 JSON 文件，图片按条目 ID 分别存放：

```
.
├── tmdbList.json
└── images/
    ├── 12345/
    │   ├── poster.jpg
    │   └── backdrop.jpg
    └── 67890/
        └── poster.jpg
```

JSON 格式示例：

```json
{
  "id": 12345,
  "title": "Inception",
  "original_title": "Inception",
  "overview": "...",
  "poster_path": "12345/poster.jpg",
  "backdrop_path": "12345/backdrop.jpg",
  "vote_average": 8.8
}
```

## 本地开发

```bash
# 克隆仓库
git clone https://github.com/zxc7563598/tmdb-list-exporter.git
cd tmdb-list-exporter

# 安装依赖
go mod download

# 编译当前平台
go build -o tmdb .

# 交叉编译全部平台
make build-all
```

## 为什么做这个工具

我想在个人网站上展示一个电影/番剧墙，但手动收集影视信息、找封面图、处理图片外链都很麻烦，也没有一个平台能统一获取所有信息。于是我选择直接从 TMDB 抓取片单数据并本地化处理——所有信息 JSON 化、图片自托管，展示时不再依赖外部 CDN。如果你也在做类似的事情，希望这个工具能帮到你。

## 致谢

- 数据来源：[TMDB](https://www.themoviedb.org/)
- 项目快速概览：[Zread](https://zread.ai/zxc7563598/tmdb-list-exporter)
