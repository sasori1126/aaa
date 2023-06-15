package images

import (
	"axis/ecommerce-backend/internal/storage"
)

type ImageRepoDb struct {
	client storage.Storage
}

func (d ImageRepoDb) DeleteImage(id uint) error {
	return d.client.DeleteImage(id)
}

func NewImageRepoDb(db storage.Storage) ImageRepo {
	return &ImageRepoDb{client: db}
}
