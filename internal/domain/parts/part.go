package parts

import (
	"axis/ecommerce-backend/internal/models"
	"context"
)

type PartRepo interface {
	CreatePart(data models.Part) error
	UpdatePart(data models.Part) error
	UpdateByField(part interface{}, data models.FindByField) error
	MergeParts(mainPartId uint, delIds []uint) error
	GetPartById(ctx context.Context, id uint) (*models.Part, error)
	GetParts(limit, offset int, returnZeroPrice bool) ([]models.Part, error)
	GetDuplicates(search string) ([]models.Part, error)
	SearchParts(limit, offset int, q string, returnZeroPrice bool) ([]models.Part, error)
	GetPartByField(fv models.FindByField) ([]models.Part, error)
}
