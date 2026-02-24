# 🎬 tmdb-list-exporter

<div align="center">
  <a href="./README.md">English</a>｜<a href="./README.zh-CN.md">简体中文</a>
  <hr width="50%"/>
</div>

> A CLI tool for exporting TMDB lists  
> Fetch a TMDB list and export it as JSON, with optional poster/backdrop image downloading and storage (local or OSS).

**This project has been parsed by Zread. If you need a quick overview of the project, you can click here to view it：[Understand this project](https://zread.ai/zxc7563598/tmdb-list-exporter)**

---

## ✨ Overview

​`tmdb-list-exporter` is a command-line tool designed to:

- Fetch all items from a specified TMDB list
- Export the data as structured JSON
- Download poster and backdrop images
- Support multiple storage backends (Local / Alibaba Cloud OSS)

It is especially useful for:

- Personal movie walls
- Anime/show showcase pages
- Static site generation
- NAS-based media display
- Data backup and archiving

---

## 🚀 Why I Built This

I created this tool because I wanted to build a movie/anime wall for my personal website.

The problems I ran into:

- Manually collecting movie/TV metadata is time-consuming
- Finding poster and backdrop images separately is tedious
- Referencing third-party image URLs often causes network issues
- No single platform provides all the information in a consistent way

Inspired by some NAS media solutions, I decided to directly fetch **TMDB list data** and process it locally:

- All metadata exported as JSON
- All images stored locally or uploaded to OSS
- Fully self-hosted and controllable
- No dependency on external image CDNs at runtime

This makes it much easier to build and maintain a fully customizable media showcase.

---

## 🧩 Features

- Export TMDB lists as formatted JSON
- Automatically fetch all pages of list data
- Download poster images
- Download backdrop images
- Multiple storage drivers:
  - Local filesystem
  - Alibaba Cloud OSS

- Built-in retry mechanism (handles HTTP 429 rate limiting)

---

## 📦 Installation

Go to the GitHub Releases page and download the binary for your platform:

- macOS
- Linux
- Windows

Make it executable:

```bash
chmod +x tmdb
```

Then run:

```bash
./tmdb --help
```

---

## 🛠 Usage

Run:

```bash
./tmdb --help
```

Example output:

```bash
tmdb-list-exporter is a TMDB API based list export tool.

Supports:
  • Exporting lists as formatted JSON
  • Downloading poster images
  • Multiple storage drivers (Local / Alibaba Cloud OSS)

Usage:
  tmdb [flags]
  tmdb [command]

Examples:
  # Generate default configuration file
  tmdb init

  # Run with default config file
  tmdb --config=config.yaml

  # Use a custom config file
  tmdb --config=configs/prod.yaml

Available Commands:
  help        Help about any command
  init        Generate default configuration file

Flags:
      --config string   Config file (default "config.yaml")
  -h, --help            help for tmdb
```

---

## ⚙️ Configuration

Generate a default config file:

```bash
tmdb init
```

Example `config.yaml`:

```yaml
tmdb:
  access_token: "your_tmdb_access_token"
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

## 🔐 How to Get a TMDB Access Token

1. Create a TMDB account
2. Go to API settings
3. Apply for an API key
4. Use the **API Read Access Token**

---

## 📁 Output Structure Example

If using local storage:

```
output/
  12345/
    poster.jpg
    backdrop.jpg
  67890/
    poster.jpg
```

Example JSON output:

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

# ⭐ If This Project Helps You

Please consider giving it a Star ⭐

If you're also building a movie wall, anime wall, NAS showcase, or media archive site, I hope this tool makes your workflow easier.
