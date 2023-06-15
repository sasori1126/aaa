package controls

import (
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/internal/storage"
)

type ControllerRepoDb struct {
	client storage.Storage
}

func (d ControllerRepoDb) UpdateController(data *models.Controller) error {
	return d.client.UpdateController(data)
}

func (d ControllerRepoDb) DeleteController(id uint) error {
	return d.client.DeleteController(id)
}

func (d ControllerRepoDb) GetControllers(limit, offset int) ([]models.Controller, error) {
	return d.client.GetControllers(limit, offset)
}

func (d ControllerRepoDb) CreateController(data models.Controller) error {
	return d.client.CreateController(data)
}

func NewControllerRepoDb(db storage.Storage) ControllerRepo {
	return &ControllerRepoDb{client: db}
}
