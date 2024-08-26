package slips

import (
	"fmt"
	"sailormoon/backend/database"
)

type SlipsService struct{}

type Filter struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

type GetSlipsParams struct {
	SortBy       string
	SortOrder    string
	Page         int
	PageSize     int
	Filters      []Filter `json:"filters"`
	GlobalFilter string   `json:"global_filter"`
}

func (s *SlipsService) GetSlips(params GetSlipsParams) ([]database.SlipsEntity, int64, error) {
	var slips []database.SlipsEntity
	var totalRecords int64

	if params.SortBy == "" {
		params.SortBy = "id"
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

	// Base query
	query := database.DB.Model(&database.SlipsEntity{}).Preload("Boats").Preload("Boats.Owners").Preload("Boats.Slips")

	// Apply global filter
	if params.GlobalFilter != "" {
		globalFilterValue := fmt.Sprintf("%%%s%%", params.GlobalFilter)
		query = query.Where("number LIKE ? OR notes LIKE ?", globalFilterValue, globalFilterValue)
	}

	// Apply individual filters
	for _, filter := range params.Filters {
		if filter.ID != "" && filter.Value != "" {
			filterValue := fmt.Sprintf("%%%s%%", filter.Value)
			query = query.Where(fmt.Sprintf("%s LIKE ?", filter.ID), filterValue)
		}
	}

	// Get total count
	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting, pagination, and execute the query
	if err := query.Order(fmt.Sprintf("%s %s", params.SortBy, params.SortOrder)).
		Limit(params.PageSize).
		Offset(offset).
		Find(&slips).Error; err != nil {
		return nil, 0, err
	}

	return slips, totalRecords, nil
}
