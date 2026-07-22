# tmdb-list-exporter

English ｜ [简体中文](./README.zh-CN.md)

A command-line tool for exporting TMDB lists. Fetches all items from a specified TMDB list, exports them as formatted JSON, and downloads poster/backdrop images with support for multiple storage backends including local filesystem and Alibaba Cloud OSS.

> This project has been parsed by [Zread](https://zread.ai/zxc7563598/tmdb-list-exporter). Click the link for a quick overview of the project structure and code.

## Use Cases

- Building a personal movie wall or anime showcase
- Providing structured movie/TV data for static sites
- NAS-based media display and data archiving
- Mirroring TMDB images to your own storage

## Features

- Automatic pagination with 3-way concurrent page fetching
- 8-way concurrent poster and backdrop downloading
- Multiple storage drivers: local filesystem / Alibaba Cloud OSS, usable simultaneously
- Exponential backoff retry for HTTP 429 rate limiting and transient network errors
- Environment variable support for configuration overrides (`TMDB_` prefix)
- Cross-compilation for 5 platforms (macOS / Linux / Windows, amd64 + arm64)

## Tech Stack

- Language: Go 1.25
- CLI Framework: [Cobra](https://github.com/spf13/cobra) + [Viper](https://github.com/spf13/viper)
- OSS SDK: [alibabacloud-oss-go-sdk-v2](https://github.com/aliyun/alibabacloud-oss-go-sdk-v2)

## Installation

Download the binary for your platform from [GitHub Releases](https://github.com/zxc7563598/tmdb-list-exporter/releases):

- macOS (amd64 / arm64)
- Linux (amd64 / arm64)
- Windows (amd64)

Make it executable and run:

```bash
chmod +x tmdb
./tmdb --help
```

You can also build from source — see [Development](#development).

## Quick Start

### 1. Get a TMDB Access Token

1. Sign up for a [TMDB](https://www.themoviedb.org/) account
2. Go to [API Settings](https://www.themoviedb.org/settings/api)
3. Apply for an API key and use the **API Read Access Token** (not the API Key)

### 2. Generate a Config File

```bash
tmdb init
```

This creates a `config.yaml` in the current directory. Fill in the values as instructed by the comments.

### 3. Run

```bash
tmdb --config=config.yaml
```

## Configuration

Full configuration example:

```yaml
tmdb:
  access_token: "your TMDB Read Access Token"
  list_id: 8634743  # List ID from the TMDB list URL

output:
  file: tmdbList.json  # Output JSON filename

storage:
  drivers:  # Leave empty to skip image storage
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
> `alioss_path` is the upload path — do **not** prefix it with `/`.

### Finding the List ID

Open the TMDB list page in your browser. The numeric part of the URL is the list_id. For example, `https://www.themoviedb.org/list/8634743-anime-series` has the list ID `8634743`.

### Environment Variables

All configuration options can be overridden with environment variables using the `TMDB_` prefix and underscore-separated keys:

```bash
export TMDB_ACCESS_TOKEN="your-token"
export TMDB_LIST_ID="8634743"
```

Environment variables take precedence over the config file.

## Output Structure

After running, a JSON file is generated in the current directory, and images are organized by item ID:

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

## Development

```bash
# Clone the repository
git clone https://github.com/zxc7563598/tmdb-list-exporter.git
cd tmdb-list-exporter

# Install dependencies
go mod download

# Build for the current platform
go build -o tmdb .

# Cross-compile for all platforms
make build-all
```

## Why I Built This

I wanted to display a movie/anime wall on my personal website, but manually gathering metadata, finding poster images, and dealing with external image hotlinking was a hassle — and no single platform provided everything in one place. So I wrote this tool to pull TMDB list data and process it locally: all metadata as JSON, all images self-hosted, no external CDN dependencies at display time. If you're building something similar, I hope it helps.

## Acknowledgements

- Data sourced from [TMDB](https://www.themoviedb.org/)
- Project overview powered by [Zread](https://zread.ai/zxc7563598/tmdb-list-exporter)
