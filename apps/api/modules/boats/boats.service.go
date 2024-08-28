package boats

import (
	"sailormoon/backend/database"
	"sailormoon/backend/utils"
)

type BoatsService struct{}

func (s *BoatsService) GetBoats(params utils.QueryParams) ([]database.BoatsEntity, int64, error) {
	var boats []database.BoatsEntity
	var totalRecords int64

	query := database.DB.Model(&database.BoatsEntity{}).Preload("Owners").Preload("Slips")

	query = utils.ApplyFiltering(query, params)

	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	query = utils.ApplySortingAndPagination(query, params, "id")

	if err := query.Find(&boats).Error; err != nil {
		return nil, 0, err
	}

	return boats, totalRecords, nil
}
