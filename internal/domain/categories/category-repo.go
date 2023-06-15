package categories

import (
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/internal/storage"
)

type CategoryRepoDb struct {
	client storage.Storage
}

func (d CategoryRepoDb) UpdateDiagramCat(data *models.DiagramCat) error {
	return d.client.UpdateDiagramCat(data)
}

func (d CategoryRepoDb) DeleteCategoryById(id uint) error {
	return d.client.DeleteCategory(id)
}

func (d CategoryRepoDb) GetCategoryById(id uint) (*models.DiagramCat, error) {
	return d.client.GetCategoryByField(models.FindByField{Field: "id", Value: id})
}

func (d CategoryRepoDb) GetCategories(limit, offset int) ([]models.DiagramCat, error) {
	return d.client.GetCategories(limit, offset)
}

func (d CategoryRepoDb) CreateDiagramCat(data models.DiagramCat) error {
	return d.client.CreateDiagramCat(data)
}

func (d CategoryRepoDb) CreateDiagramSubCat(data models.DiagramSubCat) error {
	return d.client.CreateDiagramSubCat(data)
}

func NewCategoryRepoDb(db storage.Storage) CategoryRepo {
	return &CategoryRepoDb{client: db}
}
