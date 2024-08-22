package slips

import (
	"sailormoon/backend/database"
)

type SlipsService struct{}

func (s *SlipsService) GetSlips() ([]database.SlipsEntity, error) {
	var slips []database.SlipsEntity
	if err := database.DB.Find(&slips).Error; err != nil {
		return nil, err
	}
	return slips, nil
}
