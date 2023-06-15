package heads

import "axis/ecommerce-backend/internal/models"

type HeadRepo interface {
	GetHeads(limit, offset int) ([]models.Head, error)
	GetHeadsByType(id uint, limit, offset int) ([]models.Head, error)
	GetHeadById(id uint) (*models.Head, error)
	CreateHead(data models.Head) error
}
