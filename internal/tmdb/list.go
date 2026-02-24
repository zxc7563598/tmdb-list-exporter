package tmdb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// ListItem 表示 TMDB 列表中的单个条目。
//
// 该结构同时兼容：
//   - 电影 (movie)
//   - 电视剧 (tv)
//
// 因 TMDB 不同类型字段名称不同（如 title / name），
// 所以多个字段使用指针类型以区分“字段不存在”和“空字符串”。
type ListItem struct {
	// ID 为条目唯一标识
	ID int64 `json:"id"`
	// BackdropPath 为背景图相对路径
	BackdropPath *string `json:"backdrop_path"`
	// Name 为电视剧名称（部分类型使用该字段）
	Name *string `json:"name,omitempty"`
	// Title 为电影标题（部分类型使用该字段）
	Title *string `json:"title,omitempty"`
	// OriginalName 为原始语言电视剧名称
	OriginalName *string `json:"original_name,omitempty"`
	// OriginalTitle 为原始语言电影标题
	OriginalTitle *string `json:"original_title,omitempty"`
	// Overview 为剧情简介
	Overview string `json:"overview"`
	// PosterPath 为海报相对路径
	PosterPath *string `json:"poster_path"`
	// Popularity 为热度值（TMDB 内部计算）
	Popularity float64 `json:"popularity"`
	// FirstAirDate 为首播日期（电视剧）
	FirstAirDate *string `json:"first_air_date,omitempty"`
	// ReleaseDate 为上映日期（电影）
	ReleaseDate *string `json:"release_date,omitempty"`
	// VoteAverage 为平均评分（0~10）
	VoteAverage float64 `json:"vote_average"`
	// VoteCount 为评分人数
	VoteCount int `json:"vote_count"`
}

// ListResponse 表示 TMDB 列表接口返回结构。
//
// 包含分页信息和当前页数据。
type ListResponse struct {
	// ItemCount 为列表中条目总数
	ItemCount int `json:"item_count"`
	// Items 为当前页条目数据
	Items []ListItem `json:"items"`
	// Page 为当前页码
	Page int `json:"page"`
	// TotalPages 为总页数
	TotalPages int `json:"total_pages"`
	// TotalResults 为总结果数
	TotalResults int `json:"total_results"`
}

// GetList 获取指定列表的单页数据。
//
// 参数：
//
//	listID     - TMDB 列表 ID
//	urlOptions - 查询参数，例如：page、language 等
//
// 行为：
//  1. 构造请求 URL
//  2. 添加 Bearer Token 认证头
//  3. 调用带重试机制的请求方法
//  4. 解析 JSON 响应
//
// 返回：
//
//	*ListResponse - 当前页数据
//	error         - 请求或解析失败时返回错误
func (c *Client) GetList(listID int64, urlOptions map[string]string) (*ListResponse, error) {
	// 构造基础 URL
	urlStr := fmt.Sprintf("%s/list/%d", c.baseURL, listID)
	// 构造查询参数
	query := buildQuery(urlOptions)
	if query != "" {
		urlStr += "?" + query
	}
	// 创建请求
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	// 设置请求头
	req.Header.Set("accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)
	// 执行请求（最多重试 3 次）
	resp, err := c.doRequestWithRetry(req, 3)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("意外的状态: %d", resp.StatusCode)
	}
	// 解析响应
	var result ListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetAllListItems 获取指定列表的全部分页数据。
//
// 参数：
//
//	listID     - TMDB 列表 ID
//	urlOptions - 查询参数（会自动覆盖 page 参数）
//
// 行为：
//  1. 先请求第一页
//  2. 根据 TotalPages 循环请求剩余页面
//  3. 汇总所有结果
//
// 返回：
//
//	[]ListItem - 列表全部条目
//	error      - 任意分页请求失败时返回错误
//
// 注意：
//
//	当前实现为串行分页请求，如需加速可改为并发分页。
func (c *Client) GetAllListItems(listID int64, urlOptions map[string]string) ([]ListItem, error) {
	// 先拉第一页
	urlOptions["page"] = "1"
	firstPage, err := c.GetList(listID, urlOptions)
	if err != nil {
		return nil, err
	}
	allItems := make([]ListItem, 0, firstPage.TotalResults)
	allItems = append(allItems, firstPage.Items...)
	totalPages := firstPage.TotalPages
	fmt.Printf("获取数据完成, 共计%d页(%d条)数据, 当前已获%d页(%d条)\n", totalPages, cap(allItems), 1, len(allItems))
	// 如果只有一页，直接返回
	if totalPages <= 1 {
		return allItems, nil
	}
	// 循环拉剩余页
	for page := 2; page <= totalPages; page++ {
		urlOptions["page"] = strconv.Itoa(page)
		resp, err := c.GetList(listID, urlOptions)
		if err != nil {
			return nil, err
		}
		allItems = append(allItems, resp.Items...)
		fmt.Printf("获取数据完成, 共计%d页(%d条)数据, 当前已获%d页(%d条)\n", totalPages, cap(allItems), page, len(allItems))
	}
	return allItems, nil
}
