package manufacturers

import "axis/ecommerce-backend/internal/models"

type ManufacturerRepo interface {
	GetManufacturers(limit, offset int) ([]models.Manufacturer, error)
	CreateManufacturer(data models.Manufacturer) error
	UpdateManufacturer(data *models.Manufacturer) error
	DeleteManufacturer(id uint) error
}
