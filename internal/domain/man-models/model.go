package man_models

import "axis/ecommerce-backend/internal/models"

type ModelRepo interface {
	GetModels(limit, offset int) ([]models.Model, error)
	GetModel(id uint) (*models.Model, error)
	DeleteModel(id uint) error
	CreateModel(data models.Model, controllerIds []uint) error
	UpdateModel(data *models.Model) error
	GetModelSeriesByID(id uint) (*models.Serial, error)
}
