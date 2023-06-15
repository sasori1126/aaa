package head_types

import (
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/internal/storage"
)

type HeadTypeRepoDb struct {
	client storage.Storage
}

func (h HeadTypeRepoDb) CreateHeadType(headType models.HeadType) error {
	return h.client.CreateHeadType(headType)
}

func (h HeadTypeRepoDb) GetHeadTypes(limit, offset int) ([]models.HeadType, error) {
	return h.client.GetHeadTypes(limit, offset)
}

func NewHeadTypeRepoDb(db storage.Storage) HeadTypeRepo {
	return &HeadTypeRepoDb{client: db}
}
