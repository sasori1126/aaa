package services

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/domain/manufacturers"
	"axis/ecommerce-backend/internal/dto"
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/pkg/entities"
	"axis/ecommerce-backend/pkg/utils"
	"github.com/bugsnag/bugsnag-go/v2"
)

type ManufacturerService interface {
	CreateManufacturer(mRequest *dto.ManufacturerRequest) *entities.ApiError
	UpdateManufacturerById(hid string, mRequest *dto.ManufacturerRequest) *entities.ApiError
	DeleteManufacturerById(hid string) *entities.ApiError
	GetManufacturers(limit, offset int) ([]dto.ManufacturerResponse, *entities.ApiError)
}

type DefaultManufacturerService struct {
	repo manufacturers.ManufacturerRepo
}

func (d DefaultManufacturerService) DeleteManufacturerById(hid string) *entities.ApiError {
	id, err := models.DecodeHashId(hid)
	if err != nil || id == 0 {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	err = d.repo.DeleteManufacturer(id)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"manufacturer_id": "failed to delete"},
		)
	}

	return nil
}

func (d DefaultManufacturerService) UpdateManufacturerById(hid string, req *dto.ManufacturerRequest) *entities.ApiError {
	id, err := models.DecodeHashId(hid)
	if err != nil || id == 0 {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	data := models.Manufacturer{
		Name:        req.Name,
		Description: req.Description,
	}
	data.ID = id

	err = d.repo.UpdateManufacturer(&data)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"manufacturer_id": "failed to update"},
		)
	}
	return nil
}

func (d DefaultManufacturerService) GetManufacturers(limit, offset int) ([]dto.ManufacturerResponse, *entities.ApiError) {
	getManufacturers, err := d.repo.GetManufacturers(limit, offset)
	if err != nil {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	var manufacturerResponses []dto.ManufacturerResponse
	for _, hType := range getManufacturers {
		m := hType.ToResponse()
		manufacturerResponses = append(manufacturerResponses, m)
	}

	return manufacturerResponses, nil
}

func (d DefaultManufacturerService) CreateManufacturer(mRequest *dto.ManufacturerRequest) *entities.ApiError {
	data := models.Manufacturer{
		Name:        mRequest.Name,
		Description: mRequest.Description,
	}

	err := d.repo.CreateManufacturer(data)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"failed to create manufacturer",
			configs.ServerError,
			entities.E{"error": err.Error()},
		)
	}

	return nil
}

func NewDefaultManufacturerService(repo manufacturers.ManufacturerRepo) ManufacturerService {
	return &DefaultManufacturerService{repo: repo}
}
