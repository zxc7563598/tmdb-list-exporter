package tmdb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// ListItem 表示 TMDB 列表中的单个条目。
//
// 该结构同时兼容：
//   - 电影 (movie)
//   - 电视剧 (tv)
//
// 因 TMDB 不同类型字段名称不同（如 title / name），
// 所以多个字段使用指针类型以区分"字段不存在"和"空字符串"。
type ListItem struct {
	// ID 为条目唯一标识
	ID int64 `json:"id"`
	// BackdropPath 为背景图相对路径
	BackdropPath *string `json:"backdrop_path"`
	// Name 为电视剧名称
	Name *string `json:"name,omitempty"`
	// Title 为电影标题
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

// GetDisplayTitle 返回条目的展示标题。
// 优先使用 Title（电影），其次 Name（电视剧）。
func (li *ListItem) GetDisplayTitle() string {
	if li.Title != nil && *li.Title != "" {
		return *li.Title
	}
	if li.Name != nil && *li.Name != "" {
		return *li.Name
	}
	return ""
}

// GetOriginalTitle 返回条目的原始语言标题。
// 优先使用 OriginalTitle（电影），其次 OriginalName（电视剧）。
func (li *ListItem) GetOriginalTitle() string {
	if li.OriginalTitle != nil && *li.OriginalTitle != "" {
		return *li.OriginalTitle
	}
	if li.OriginalName != nil && *li.OriginalName != "" {
		return *li.OriginalName
	}
	return ""
}

// GetReleaseDate 返回条目发布日期。
// 电影使用 ReleaseDate，电视剧使用 FirstAirDate。
func (li *ListItem) GetReleaseDate() *string {
	if li.ReleaseDate != nil {
		return li.ReleaseDate
	}
	return li.FirstAirDate
}

// ListResponse 表示 TMDB 列表接口返回结构。
type ListResponse struct {
	ItemCount    int        `json:"item_count"`
	Items        []ListItem `json:"items"`
	Page         int        `json:"page"`
	TotalPages   int        `json:"total_pages"`
	TotalResults int        `json:"total_results"`
}

// GetList 获取指定列表的单页数据。
func (c *Client) GetList(ctx context.Context, listID int64, urlOptions map[string]string) (*ListResponse, error) {
	urlStr := fmt.Sprintf("%s/list/%d", c.baseURL, listID)
	query := buildQuery(urlOptions)
	if query != "" {
		urlStr += "?" + query
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.doRequestWithRetry(req, 3)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("意外的状态: %d", resp.StatusCode)
	}

	var result ListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// apiPageConcurrency 为分页请求的最大并发数。
// 设置保守以避免触发 TMDB API 限流（约 50 req/s）。
const apiPageConcurrency = 3

// apiRequestDelay 为分页并发请求之间的最小间隔。
const apiRequestDelay = 200 * time.Millisecond

// GetAllListItems 获取指定列表的全部分页数据。
//
// 使用有限并发（3 路）加速多页拉取，同时保持对 TMDB API 限流的友好。
func (c *Client) GetAllListItems(ctx context.Context, listID int64, urlOptions map[string]string) ([]ListItem, error) {
	// 先拉第一页以获取分页信息
	opts := copyMap(urlOptions)
	opts["page"] = "1"
	firstPage, err := c.GetList(ctx, listID, opts)
	if err != nil {
		return nil, err
	}

	allItems := make([]ListItem, 0, firstPage.TotalResults)
	allItems = append(allItems, firstPage.Items...)
	totalPages := firstPage.TotalResults

	fmt.Printf("获取数据完成, 共计%d页(%d条)数据, 当前已获1页(%d条)\n",
		firstPage.TotalPages, totalPages, len(allItems))

	if firstPage.TotalPages <= 1 {
		return allItems, nil
	}

	// 并发拉取剩余页（最多 3 路并发 + 请求间隔以友好对待 TMDB 限流）
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, apiPageConcurrency)
	var firstErr error
	errOnce := sync.Once{}

	for page := 2; page <= firstPage.TotalPages; page++ {
		// 若已有错误则不再发起新请求
		mu.Lock()
		hasErr := firstErr != nil
		mu.Unlock()
		if hasErr {
			break
		}

		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			// 检查 context 是否已取消
			select {
			case <-ctx.Done():
				errOnce.Do(func() { firstErr = ctx.Err() })
				return
			default:
			}

			// 交错延迟以减少限流风险
			time.Sleep(apiRequestDelay * time.Duration(page))

			pageOpts := copyMap(urlOptions)
			pageOpts["page"] = strconv.Itoa(page)
			resp, err := c.GetList(ctx, listID, pageOpts)
			if err != nil {
				errOnce.Do(func() { firstErr = fmt.Errorf("第 %d 页拉取失败: %w", page, err) })
				return
			}

			mu.Lock()
			allItems = append(allItems, resp.Items...)
			fmt.Printf("获取数据完成, 共计%d页(%d条)数据, 当前已获%d页(%d条)\n",
				firstPage.TotalPages, totalPages, page, len(allItems))
			mu.Unlock()
		}(page)
	}

	wg.Wait()
	if firstErr != nil {
		return nil, firstErr
	}

	return allItems, nil
}

// copyMap 浅拷贝一个 map[string]string。
func copyMap(src map[string]string) map[string]string {
	dst := make(map[string]string, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
