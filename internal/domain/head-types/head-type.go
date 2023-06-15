package head_types

import "axis/ecommerce-backend/internal/models"

type HeadTypeRepo interface {
	GetHeadTypes(limit, offset int) ([]models.HeadType, error)
	CreateHeadType(headType models.HeadType) error
}
