package slips

import (
	"fmt"
	"sailormoon/backend/database"
)

type SlipsService struct{}

type GetSlipsParams struct {
	SortBy    string
	SortOrder string
	Page      int
	PageSize  int
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

	if err := database.DB.Model(&database.SlipsEntity{}).Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	if err := database.DB.Preload("Boats").Preload("Boats.Owners").Preload("Boats.Slips").
		Order(fmt.Sprintf("%s %s", params.SortBy, params.SortOrder)).
		Limit(params.PageSize).
		Offset(offset).
		Find(&slips).Error; err != nil {
		return nil, 0, err
	}

	return slips, totalRecords, nil
}
