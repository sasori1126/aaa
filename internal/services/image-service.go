package services

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/domain/images"
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/pkg/entities"
	"axis/ecommerce-backend/pkg/utils"
	"github.com/bugsnag/bugsnag-go/v2"
)

type ImageService interface {
	DeleteImage(hid string) *entities.ApiError
}

type DefaultImageService struct {
	repo images.ImageRepo
}

func (d DefaultImageService) DeleteImage(hid string) *entities.ApiError {
	id, err := models.DecodeHashId(hid)
	if err != nil || id == 0 {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	err = d.repo.DeleteImage(id)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"model_id": "failed to delete"},
		)
	}

	return nil
}

func NewDefaultImageService(repo images.ImageRepo) ImageService {
	return &DefaultImageService{repo: repo}
}
