package services

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/domain/controls"
	"axis/ecommerce-backend/internal/dto"
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/pkg/entities"
	"axis/ecommerce-backend/pkg/utils"
	"github.com/bugsnag/bugsnag-go/v2"
)

type ControllerService interface {
	CreateController(cRequest *dto.ControllerRequest) *entities.ApiError
	UpdateController(hid string, cRequest *dto.ControllerRequest) *entities.ApiError
	DeleteController(hid string) *entities.ApiError
	GetControllers(limit, offset int) ([]dto.ControllerResponse, *entities.ApiError)
}

type DefaultControllerService struct {
	repo controls.ControllerRepo
}

func (d DefaultControllerService) UpdateController(hid string, req *dto.ControllerRequest) *entities.ApiError {
	id, err := models.DecodeHashId(hid)
	if err != nil || id == 0 {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	data := models.Controller{
		Name:        req.Name,
		Description: req.Description,
	}
	data.ID = id

	err = d.repo.UpdateController(&data)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"controller_id": "failed to update"},
		)
	}
	return nil
}

func (d DefaultControllerService) DeleteController(hid string) *entities.ApiError {
	id, err := models.DecodeHashId(hid)
	if err != nil || id == 0 {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	err = d.repo.DeleteController(id)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"controller_id": "failed to delete"},
		)
	}

	return nil
}

func (d DefaultControllerService) GetControllers(limit, offset int) ([]dto.ControllerResponse, *entities.ApiError) {
	getControllers, err := d.repo.GetControllers(limit, offset)
	if err != nil {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	var controllers []dto.ControllerResponse
	for _, c := range getControllers {
		ct := c.ToResponse()
		controllers = append(controllers, ct)
	}

	return controllers, nil
}

func (d DefaultControllerService) CreateController(cRequest *dto.ControllerRequest) *entities.ApiError {
	data := models.Controller{
		Name:        cRequest.Name,
		Description: cRequest.Description,
	}

	err := d.repo.CreateController(data)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"failed to create controller",
			configs.ServerError,
			entities.E{"error": err.Error()},
		)
	}

	return nil
}

func NewDefaultControllerService(repo controls.ControllerRepo) ControllerService {
	return &DefaultControllerService{repo: repo}
}
