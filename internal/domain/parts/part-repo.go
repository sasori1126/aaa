package parts

import (
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/internal/storage"
	"context"
)

type PartRepoDb struct {
	client storage.Storage
}

func (p PartRepoDb) UpdatePart(data models.Part) error {
	return p.client.UpdatePart(data)
}

func (p PartRepoDb) MergeParts(mainPartId uint, delIds []uint) error {
	return p.client.MergeParts(mainPartId, delIds)
}

func (p PartRepoDb) UpdateByField(part interface{}, fv models.FindByField) error {
	return p.client.UpdateByField(part, map[string]interface{}{fv.Field: fv.Value})
}

func (p PartRepoDb) GetPartByField(fv models.FindByField) ([]models.Part, error) {
	return p.client.GetPartsByField(fv)
}

func (p PartRepoDb) SearchParts(limit, offset int, search string, returnZeroPrice bool) ([]models.Part, error) {
	return p.client.SearchParts(limit, offset, search, returnZeroPrice)
}

func (p PartRepoDb) GetDuplicates(search string) ([]models.Part, error) {
	return p.client.GetDuplicates(search)
}

func (p PartRepoDb) GetPartById(ctx context.Context, id uint) (*models.Part, error) {
	return p.client.GetPartByField(ctx, models.QueryByField{Query: "id = ?", Value: id})
}

func (p PartRepoDb) GetParts(limit, offset int, returnZeroPrice bool) ([]models.Part, error) {
	return p.client.GetParts(limit, offset, returnZeroPrice)
}

func (p PartRepoDb) CreatePart(data models.Part) error {
	return p.client.CreatePart(data)
}

func NewPartRepoDb(db storage.Storage) PartRepo {
	return &PartRepoDb{client: db}
}
