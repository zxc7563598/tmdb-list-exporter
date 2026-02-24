# 🎬 tmdb-list-exporter

<div align="center">
  <a href="./README.md">English</a>｜<a href="./README.zh-CN.md">简体中文</a>
  <hr width="50%"/>
</div>

> 一个基于 TMDB API 的片单导出 CLI 工具  
> 用于快速获取指定片单信息为 JSON 文件，并支持下载封面 / 背景图进行本地或 OSS 存储。

**本项目已经经由 Zread 解析完成，如果需要快速了解项目，可以点击此处进行查看：[了解本项目](https://zread.ai/zxc7563598/tmdb-list-exporter)**

---

## ✨ 项目简介

​`tmdb-list-exporter` 是一个命令行工具，用于：

- 抓取指定 TMDB 片单的完整数据
- 导出为结构化 JSON 文件
- 下载封面图与背景图
- 支持多种存储方式（本地 / 阿里云 OSS）

它特别适合：

- 个人影视墙
- 番剧展示页
- 静态站点生成
- NAS 本地影视展示
- 数据备份与归档

---

## 🚀 为什么做这个工具？

我做这个工具的初衷，是为了在个人网站上展示一个电影 / 番剧墙。

问题在于：

- 手动收集影视信息非常麻烦
- 还要额外找封面图和背景图
- 引用第三方图片资源容易出现网络问题
- 没有一个平台可以完整统一地获取所有信息

于是参考一些 NAS 方案，我选择直接抓取 ​**TMDB 的片单数据**，并进行本地化处理：

- 所有信息 JSON 化
- 所有图片本地存储或上传 OSS
- 完全自主可控
- 展示时无需依赖外部图片 CDN

这样就可以非常方便地构建属于自己的影视展示页面。

---

## 🧩 功能特性

- 导出 TMDB 片单为格式化 JSON 文件
- 自动分页拉取全部数据
- 下载封面图（Poster）
- 下载背景图（Backdrop）
- 多存储驱动支持
  - 本地文件系统
  - 阿里云 OSS

- 自动重试机制（429 限流处理）

---

## 📦 安装方式

前往 GitHub Releases 页面下载对应平台的二进制文件：

- macOS
- Linux
- Windows

下载后赋予执行权限：

```bash
chmod +x tmdb
```

然后即可运行：

```bash
./tmdb --help
```

---

## 🛠 使用方法

执行：

```bash
./tmdb --help
```

输出如下：

```bash
tmdb-list-exporter 是一个基于 TMDB API 的片单导出工具。

支持：
  • 导出片单为格式化 JSON 文件
  • 下载封面图片
  • 多存储驱动( 本地 / 阿里云 OSS )

Usage:
  tmdb [flags]
  tmdb [command]

Examples:
  # 生成默认配置文件
  tmdb init
  # 使用默认配置文件运行
  tmdb --config=config.yaml
  # 指定其他配置文件
  tmdb --config=configs/prod.yaml

Available Commands:
  help        Help about any command
  init        生成默认配置文件

Flags:
      --config string   配置文件(默认为config.yaml)
  -h, --help            help for tmdb
```

---

## ⚙️ 配置说明

你可以通过：

```bash
tmdb init
```

生成默认配置文件。

示例 `config.yaml`：

```yaml
tmdb:
  access_token: "你的TMDB访问令牌"
  list_id: 123456

storage:
  drivers:
    - type: local
      local_path: "./output"

    - type: alioss
      alioss_region: "cn-shanghai"
      alioss_bucket: "your-bucket"
      alioss_path: "images"
      alioss_access_key_id: "your-ak"
      alioss_access_key_secret: "your-sk"
```

---

## 🔐 如何获取 TMDB Access Token

1. 注册 TMDB 账号
2. 进入 API 设置页面
3. 申请 API Key
4. 使用 **API Read Access Token**

---

## 📁 输出结构示例

假设你启用了本地存储：

```
output/
  12345/
    poster.jpg
    backdrop.jpg
  67890/
    poster.jpg
```

JSON 示例：

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

---

# ⭐ 如果这个工具对你有帮助

欢迎 Star ⭐

如果你也在做影视墙、番剧墙、NAS 展示页，希望这个工具能帮到你。
