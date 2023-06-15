package heads

import (
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/internal/storage"
)

type HeadRepoDb struct {
	client storage.Storage
}

func (h HeadRepoDb) CreateHead(data models.Head) error {
	return h.client.CreateHead(data)
}

func (h HeadRepoDb) GetHeadsByType(id uint, limit, offset int) ([]models.Head, error) {
	return h.client.GetHeadsByType(id, limit, offset)
}

func (h HeadRepoDb) GetHeadById(id uint) (*models.Head, error) {
	head, err := h.client.FindHeadByField(models.FindByField{Field: "id", Value: id})
	if err != nil {
		return nil, err
	}

	return head, nil
}

func (h HeadRepoDb) GetHeads(limit, offset int) ([]models.Head, error) {
	return h.client.GetHeads(limit, offset)
}

func NewHeadRepoDb(db storage.Storage) HeadRepo {
	return &HeadRepoDb{client: db}
}
