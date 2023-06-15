package man_models

import (
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/internal/storage"
)

type ModelRepoDb struct {
	client storage.Storage
}

func (d ModelRepoDb) GetModelSeriesByID(id uint) (*models.Serial, error) {
	return d.client.GetSerialByID(id)
}

func (d ModelRepoDb) GetModel(id uint) (*models.Model, error) {
	return d.client.GetModel(id)
}

func (d ModelRepoDb) DeleteModel(id uint) error {
	return d.client.DeleteModel(id)
}

func (d ModelRepoDb) UpdateModel(data *models.Model) error {
	return d.client.UpdateModel(data)
}

func (d ModelRepoDb) CreateModel(data models.Model, controllerIds []uint) error {
	return d.client.CreateModel(data, controllerIds)
}

func (d ModelRepoDb) GetModels(limit, offset int) ([]models.Model, error) {
	return d.client.GetModels(limit, offset)
}

func NewModelRepoDb(db storage.Storage) ModelRepo {
	return &ModelRepoDb{client: db}
}
