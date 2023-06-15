package categories

import "axis/ecommerce-backend/internal/models"

type CategoryRepo interface {
	GetCategories(limit, offset int) ([]models.DiagramCat, error)
	GetCategoryById(id uint) (*models.DiagramCat, error)
	DeleteCategoryById(id uint) error
	CreateDiagramCat(data models.DiagramCat) error
	UpdateDiagramCat(data *models.DiagramCat) error
	CreateDiagramSubCat(data models.DiagramSubCat) error
}
