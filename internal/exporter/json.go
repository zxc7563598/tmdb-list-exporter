package exporter

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"

	"github.com/zxc7563598/tmdb-list-exporter/internal/tmdb"
)

// Item 为导出使用的标准数据结构。
type Item struct {
	ID                  int64    `json:"id"`
	Title               string   `json:"title"`
	Original            string   `json:"original_title"`
	Overview            string   `json:"overview"`
	PosterPath          string   `json:"poster_path"`
	BackdropPath        string   `json:"backdrop_path"`
	OriginalPosterURL   *string  `json:"original_poster_url,omitempty"`
	OriginalBackdropURL *string  `json:"original_backdrop_url,omitempty"`
	Popularity          float64  `json:"popularity"`
	ReleaseDate         *string  `json:"release_date,omitempty"`
	VoteAverage         float64  `json:"vote_average"`
	VoteCount           int      `json:"vote_count"`
}

// FromTMDB 将 TMDB API 返回的 ListItem 转换为导出结构 Item。
func FromTMDB(items []tmdb.ListItem) []Item {
	result := make([]Item, 0, len(items))

	for _, i := range items {
		// 海报处理
		var posterURL *string
		var posterPath string
		if i.PosterPath != nil && *i.PosterPath != "" {
			url := tmdb.ImageBaseURL + *i.PosterPath
			posterURL = &url
			posterPath = filepath.Join(strconv.FormatInt(i.ID, 10), *i.PosterPath)
		}

		// 背景图处理
		var backdropURL *string
		var backdropPath string
		if i.BackdropPath != nil && *i.BackdropPath != "" {
			url := tmdb.ImageBaseURL + *i.BackdropPath
			backdropURL = &url
			backdropPath = filepath.Join(strconv.FormatInt(i.ID, 10), *i.BackdropPath)
		}

		result = append(result, Item{
			ID:                  i.ID,
			Title:               i.GetDisplayTitle(),
			Original:            i.GetOriginalTitle(),
			Overview:            i.Overview,
			PosterPath:          posterPath,
			OriginalPosterURL:   posterURL,
			BackdropPath:        backdropPath,
			OriginalBackdropURL: backdropURL,
			Popularity:          i.Popularity,
			ReleaseDate:         i.GetReleaseDate(),
			VoteAverage:         i.VoteAverage,
			VoteCount:           i.VoteCount,
		})
	}
	return result
}

// JSONExporter 定义 JSON 导出器。
type JSONExporter struct {
	Indent string
}

// NewJSONExporter 创建 JSONExporter 实例。
func NewJSONExporter(indent string) *JSONExporter {
	if indent == "" {
		indent = "  "
	}
	return &JSONExporter{Indent: indent}
}

// Export 将数据导出为 JSON 文件。
func (e *JSONExporter) Export(path string, data any) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", e.Indent)
	return encoder.Encode(data)
}
