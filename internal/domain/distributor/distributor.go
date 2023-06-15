package distributor

import (
	"axis/ecommerce-backend/internal/models"
)

type RepoDistributor interface {
	CreateDistributor(req models.Distributor) error
	GetDistributors(limit, offset int) ([]models.Distributor, error)
	GetDistributor(id uint) (*models.Distributor, error)
	DeleteDistributor(id uint) error
	UpdateDistributor(data *models.Distributor) error
}
