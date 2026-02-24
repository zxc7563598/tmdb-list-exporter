package exporter

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"

	"github.com/zxc7563598/tmdb-list-exporter/internal/tmdb"
)

// Item 为导出使用的标准数据结构。
//
// 该结构体是对 TMDB 原始返回数据的“扁平化 + 统一格式”封装，
// 用于最终 JSON 输出。
//
// 相比 TMDB 原始结构：
//   - 统一电影 / 电视剧字段
//   - 生成原始图片完整 URL
//   - 生成图片在存储系统中的相对路径
type Item struct {
	// ID 为 TMDB 条目 ID（电影或电视剧）
	ID int64 `json:"id"`
	// Title 为展示标题（优先取 Title，其次 Name）
	Title string `json:"title"`
	// Original 为原始语言标题
	Original string `json:"original_title"`
	// Overview 为剧情简介
	Overview string `json:"overview"`
	// PosterPath 为海报在存储系统中的相对路径
	// 格式示例：12345/abc.jpg
	PosterPath string `json:"poster_path"`
	// BackdropPath 为背景图在存储系统中的相对路径
	BackdropPath string `json:"backdrop_path"`
	// OriginalPosterURL 为 TMDB 原始海报完整 URL（可为空）
	// 例如：https://image.tmdb.org/t/p/original/xxx.jpg
	OriginalPosterURL *string `json:"original_poster_url,omitempty"`
	// OriginalBackdropURL 为 TMDB 原始背景图完整 URL（可为空）
	OriginalBackdropURL *string `json:"original_backdrop_url,omitempty"`
	// Popularity 为 TMDB 热度评分
	Popularity float64 `json:"popularity"`
	// ReleaseDate 为上映日期（电影）或首播日期（电视剧）
	ReleaseDate *string `json:"release_date,omitempty"`
	// VoteAverage 为平均评分
	VoteAverage float64 `json:"vote_average"`
	// VoteCount 为评分人数
	VoteCount int `json:"vote_count"`
}

// FromTMDB 将 TMDB API 返回的 ListItem 转换为导出结构 Item。
//
// 转换逻辑：
//   - 统一电影与电视剧字段
//   - 生成原始图片 URL
//   - 生成图片存储相对路径
//
// 参数：
//
//	items - TMDB 返回的列表数据
//
// 返回：
//
//	[]Item - 可用于导出的标准结构
func FromTMDB(items []tmdb.ListItem) []Item {
	var imageBaseURL = "https://image.tmdb.org/t/p/original"
	result := make([]Item, 0, len(items))

	for _, i := range items {
		// 标题处理
		// 电影使用 Title，电视剧使用 Name
		title := ""
		if i.Title != nil {
			title = *i.Title
		} else if i.Name != nil {
			title = *i.Name
		}
		// 原始标题处理
		// 与标题同理
		original := ""
		if i.OriginalTitle != nil {
			original = *i.OriginalTitle
		} else if i.OriginalName != nil {
			original = *i.OriginalName
		}
		// 海报处理
		var posterURL *string
		var posterPath string
		if i.PosterPath != nil {
			// 拼接原始图片 URL
			url := imageBaseURL + *i.PosterPath
			posterURL = &url
			// 存储路径：ID/文件名
			posterPath = filepath.Join(strconv.FormatInt(i.ID, 10), *i.PosterPath)
		}
		// 背景图处理
		var backdropURL *string
		var backdropPath string
		if i.BackdropPath != nil {
			// 拼接原始图片 URL
			url := imageBaseURL + *i.BackdropPath
			backdropURL = &url
			// 存储路径：ID/文件名
			backdropPath = filepath.Join(strconv.FormatInt(i.ID, 10), *i.BackdropPath)
		}
		// 上映日期处理
		// 电影使用 ReleaseDate，电视剧使用 FirstAirDate
		release := i.ReleaseDate
		if release == nil {
			release = i.FirstAirDate
		}
		result = append(result, Item{
			ID:                  i.ID,
			Title:               title,
			Original:            original,
			Overview:            i.Overview,
			PosterPath:          posterPath,
			OriginalPosterURL:   posterURL,
			BackdropPath:        backdropPath,
			OriginalBackdropURL: backdropURL,
			Popularity:          i.Popularity,
			ReleaseDate:         release,
			VoteAverage:         i.VoteAverage,
			VoteCount:           i.VoteCount,
		})
	}
	return result
}

// JSONExporter 定义 JSON 导出器。
//
// 用于将结构化数据输出为 JSON 文件。
type JSONExporter struct {
	Indent string
}

// NewJSONExporter 创建 JSONExporter 实例。
//
// 参数：
//
//	indent - 缩进字符串（为空时默认使用两个空格）
//
// 返回：
//
//	*JSONExporter
func NewJSONExporter(indent string) *JSONExporter {
	if indent == "" {
		indent = "  "
	}
	return &JSONExporter{Indent: indent}
}

// Export 将数据导出为 JSON 文件。
//
// 参数：
//
//	path - 输出文件路径
//	data - 需要序列化的数据
//
// 行为：
//   - 创建或覆盖目标文件
//   - 使用指定缩进格式输出 JSON
//
// 返回：
//
//	error - 文件创建或编码失败时返回错误
func (e *JSONExporter) Export(path string, data any) error {
	// 创建文件（若存在则覆盖）
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	// 创建 JSON 编码器
	encoder := json.NewEncoder(file)
	// 设置缩进格式
	encoder.SetIndent("", e.Indent)
	// 编码并写入文件
	return encoder.Encode(data)
}
