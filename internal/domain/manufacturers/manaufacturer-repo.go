package manufacturers

import (
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/internal/storage"
)

type ManufacturerRepoDb struct {
	client storage.Storage
}

func (d ManufacturerRepoDb) DeleteManufacturer(id uint) error {
	return d.client.DeleteManufacturer(id)
}

func (d ManufacturerRepoDb) UpdateManufacturer(data *models.Manufacturer) error {
	return d.client.UpdateManufacturer(data)
}

func (d ManufacturerRepoDb) GetManufacturers(limit, offset int) ([]models.Manufacturer, error) {
	return d.client.GetManufacturers(limit, offset)
}

func (d ManufacturerRepoDb) CreateManufacturer(data models.Manufacturer) error {
	return d.client.CreateManufacturer(data)
}

func NewManufacturerRepoDb(db storage.Storage) ManufacturerRepo {
	return &ManufacturerRepoDb{client: db}
}
