package services

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/domain/stats"
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/pkg/entities"
	"axis/ecommerce-backend/pkg/utils"
	"github.com/bugsnag/bugsnag-go/v2"
)

type StatService interface {
	GetTotalUsers() (int, *entities.ApiError)
	TotalStat(limit int) (*models.Stat, *entities.ApiError)
}

type DefaultStatService struct {
	repo stats.StatRepo
}

func (d DefaultStatService) TotalStat(limit int) (*models.Stat, *entities.ApiError) {
	stat, err := d.repo.GetTotalStat(limit)
	if err != nil {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError(
			"Failed to get stats",
			configs.ServerError,
			entities.E{},
		)
	}

	return stat, nil
}

func (d DefaultStatService) GetTotalUsers() (int, *entities.ApiError) {
	return 0, nil
}

func NewDefaultStatService(repo stats.StatRepo) StatService {
	return &DefaultStatService{repo: repo}
}
