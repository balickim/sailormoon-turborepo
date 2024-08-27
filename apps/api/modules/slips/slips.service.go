package slips

import (
	"sailormoon/backend/database"
	"sailormoon/backend/utils"
	"strings"

	"gorm.io/gorm"
)

func init() {
	utils.RegisterCustomFilter("status", func(query *gorm.DB, key, value string) *gorm.DB {
		return query.Where("LOWER(status) = ?", strings.ToLower(value))
	})

	utils.RegisterCustomFilter("date_range", func(query *gorm.DB, key, value string) *gorm.DB {
		dates := strings.Split(value, ":")
		if len(dates) == 2 {
			startDate, endDate := dates[0], dates[1]
			return query.Where("date BETWEEN ? AND ?", startDate, endDate)
		}
		return query
	})
}

type SlipsService struct{}

type Filter struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

func (s *SlipsService) GetSlips(params utils.QueryParams) ([]database.SlipsEntity, int64, error) {
	var slips []database.SlipsEntity
	var totalRecords int64

	query := database.DB.Model(&database.SlipsEntity{}).Preload("Boats").Preload("Boats.Owners").Preload("Boats.Slips")

	query = utils.ApplyFiltering(query, params)

	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	query = utils.ApplySortingAndPagination(query, params, "id")

	if err := query.Find(&slips).Error; err != nil {
		return nil, 0, err
	}

	return slips, totalRecords, nil
}
