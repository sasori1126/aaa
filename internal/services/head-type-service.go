package services

import (
	"axis/ecommerce-backend/configs"
	head_types "axis/ecommerce-backend/internal/domain/head-types"
	"axis/ecommerce-backend/internal/dto"
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/pkg/entities"
	"axis/ecommerce-backend/pkg/utils"
	"github.com/bugsnag/bugsnag-go/v2"
)

type HeadTypeService interface {
	GetHeadTypes(limit, offset int) ([]dto.HeadTypeResponse, *entities.ApiError)
	CreateHeadType(hr *dto.HeadTypeRequest) *entities.ApiError
}

type DefaultHeadTypeService struct {
	repo head_types.HeadTypeRepo
}

func (d DefaultHeadTypeService) CreateHeadType(hr *dto.HeadTypeRequest) *entities.ApiError {
	data := models.HeadType{
		Name:        hr.Name,
		Description: hr.Description,
	}

	err := d.repo.CreateHeadType(data)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"failed to create head type",
			configs.ServerError,
			entities.E{"error": err.Error()},
		)
	}

	return nil
}

func (d DefaultHeadTypeService) GetHeadTypes(limit, offset int) ([]dto.HeadTypeResponse, *entities.ApiError) {
	hTypes, err := d.repo.GetHeadTypes(limit, offset)
	if err != nil {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	var hts []dto.HeadTypeResponse
	for _, hType := range hTypes {
		ht := hType.ToResponse()
		hts = append(hts, ht)
	}

	return hts, nil
}

func NewDefaultHeadTypeService(repo head_types.HeadTypeRepo) HeadTypeService {
	return &DefaultHeadTypeService{repo: repo}
}
