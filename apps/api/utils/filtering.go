package utils

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Filter struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

type QueryParams struct {
	SortBy       string
	SortOrder    string
	Page         int
	PageSize     int
	Filters      []Filter `json:"filters"`
	GlobalFilter string   `json:"global_filter"`
}

type CustomFilterFunc func(query *gorm.DB, key, value string) *gorm.DB

var FilterRegistry = map[string]CustomFilterFunc{}

func ApplySortingAndPagination(query *gorm.DB, params QueryParams, defaultSortBy string) *gorm.DB {
	if params.SortBy == "" {
		params.SortBy = defaultSortBy
	}
	if params.SortOrder != "asc" && params.SortOrder != "desc" {
		params.SortOrder = "asc"
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	offset := (params.Page - 1) * params.PageSize

	query = query.Order(fmt.Sprintf("%s %s", params.SortBy, params.SortOrder)).
		Limit(params.PageSize).
		Offset(offset)

	return query
}

func ApplyFiltering(query *gorm.DB, params QueryParams) *gorm.DB {
	if params.GlobalFilter != "" {
		globalFilterValue := fmt.Sprintf("%%%s%%", params.GlobalFilter)
		query = query.Where("notes LIKE ?", globalFilterValue)

		if num, err := strconv.Atoi(params.GlobalFilter); err == nil {
			query = query.Or("number = ?", num)
		}
	}

	// Apply individual filters
	for _, filter := range params.Filters {
		if filter.ID != "" && filter.Value != "" {
			// Check if a custom filter function exists for this filter ID
			if customFilter, exists := FilterRegistry[filter.ID]; exists {
				query = customFilter(query, filter.ID, filter.Value)
			} else {
				// Default behavior if no custom filter exists
				if filter.ID == "number" {
					if num, err := strconv.Atoi(filter.Value); err == nil {
						query = query.Where("number = ?", num)
					}
				} else {
					filterValue := fmt.Sprintf("%%%s%%", filter.Value)
					query = query.Where(fmt.Sprintf("%s LIKE ?", filter.ID), filterValue)
				}
			}
		}
	}

	return query
}

func RegisterCustomFilter(filterID string, filterFunc CustomFilterFunc) {
	FilterRegistry[filterID] = filterFunc
}

func ParseQueryParams(c *fiber.Ctx) (QueryParams, map[string]interface{}, error) {
	sortBy := c.Query("sort_by", "id")
	sortOrder := c.Query("sort_order", "asc")
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 10)
	globalFilter := c.Query("global_filter", "")

	var filters []Filter
	filtersParam := c.Query("filters", "")
	if filtersParam != "" {
		if err := json.Unmarshal([]byte(filtersParam), &filters); err != nil {
			return QueryParams{}, nil, err
		}
	}

	params := QueryParams{
		SortBy:       sortBy,
		SortOrder:    sortOrder,
		Page:         page,
		PageSize:     pageSize,
		Filters:      filters,
		GlobalFilter: globalFilter,
	}

	meta := map[string]interface{}{
		"current_page": page,
		"page_size":    pageSize,
	}

	return params, meta, nil
}

func UpdateMetaWithTotal(meta map[string]interface{}, total int64, pageSize int) {
	lastPage := int(math.Ceil(float64(total) / float64(pageSize)))
	meta["total"] = total
	meta["last_page"] = lastPage
}
