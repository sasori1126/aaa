package controls

import "axis/ecommerce-backend/internal/models"

type ControllerRepo interface {
	CreateController(data models.Controller) error
	UpdateController(data *models.Controller) error
	DeleteController(id uint) error
	GetControllers(limit, offset int) ([]models.Controller, error)
}
