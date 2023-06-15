package services

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/domain/man-models"
	"axis/ecommerce-backend/internal/dto"
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/pkg/entities"
	"axis/ecommerce-backend/pkg/utils"
	"errors"
	"github.com/bugsnag/bugsnag-go/v2"
	"gorm.io/gorm"
)

type ModelService interface {
	GetModels(limit, offset int) ([]dto.ModelResponse, *entities.ApiError)
	GetModel(hid string) (*dto.ModelResponse, *entities.ApiError)
	CreateModel(mRequest *dto.ModelRequest) *entities.ApiError
	UpdateModel(hid string, req *dto.ModelRequest) *entities.ApiError
	DeleteModel(hid string) *entities.ApiError
}

type DefaultModelService struct {
	repo man_models.ModelRepo
}

func (d DefaultModelService) GetModel(hid string) (*dto.ModelResponse, *entities.ApiError) {
	id, err := models.DecodeHashId(hid)
	if err != nil || id == 0 {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError(
			"invalid id",
			configs.BadRequest,
			entities.E{"id": "invalid id"},
		)
	}

	getModel, err := d.repo.GetModel(id)
	if err != nil {
		bugsnag.Notify(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.FormatApiError(
				"Model not found",
				configs.NotFound,
				entities.E{},
			)
		}

		return nil, utils.FormatApiError(
			"failed to retrieve model",
			configs.ServerError,
			entities.E{"server": "server error"},
		)
	}

	mls := getModel.ToResponse()

	return &mls, nil
}

func (d DefaultModelService) UpdateModel(hid string, req *dto.ModelRequest) *entities.ApiError {
	id, err := models.DecodeHashId(hid)
	if err != nil || id == 0 {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"invalid id",
			configs.BadRequest,
			entities.E{"id": "invalid id"},
		)
	}

	manId, err := models.DecodeHashId(req.ManufacturerID)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error, could not decode id"},
		)
	}

	var serials []*models.Serial
	for _, s := range req.Serials {
		sr := models.Serial{
			SerialStart: s.SerialStart,
			SerialEnd:   s.SerialEnd,
			YearStart:   s.YearStart,
			YearEnd:     s.YearEnd,
		}

		serials = append(serials, &sr)
	}

	manf := models.Manufacturer{}
	manf.ID = manId

	data := models.Model{
		Manufacturer: manf,
		Name:         req.Name,
		Description:  req.Description,
		Serials:      serials,
	}
	data.ID = id

	var controllers []*models.Controller

	for _, id := range req.ControllerIds {
		controllerId, err := models.DecodeHashId(id)
		if err != nil {
			bugsnag.Notify(err)
			return utils.FormatApiError(
				"failed to create models",
				configs.ServerError,
				entities.E{"error": err.Error()},
			)
		}

		c := models.Controller{}
		c.ID = controllerId

		controllers = append(controllers, &c)
	}
	data.Controllers = controllers

	err = d.repo.UpdateModel(&data)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"failed to update models",
			configs.ServerError,
			entities.E{"error": err.Error()},
		)
	}

	return nil
}

func (d DefaultModelService) DeleteModel(hid string) *entities.ApiError {
	id, err := models.DecodeHashId(hid)
	if err != nil || id == 0 {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	err = d.repo.DeleteModel(id)
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

func (d DefaultModelService) CreateModel(mRequest *dto.ModelRequest) *entities.ApiError {
	manId, err := models.DecodeHashId(mRequest.ManufacturerID)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error, could not decode id"},
		)
	}

	var serials []*models.Serial
	for _, s := range mRequest.Serials {
		sr := models.Serial{
			SerialStart: s.SerialStart,
			SerialEnd:   s.SerialEnd,
			YearStart:   s.YearStart,
			YearEnd:     s.YearEnd,
		}

		serials = append(serials, &sr)
	}

	data := models.Model{
		ManufacturerID: manId,
		Name:           mRequest.Name,
		Description:    mRequest.Description,
		Serials:        serials,
	}

	var controllerIds []uint

	for _, id := range mRequest.ControllerIds {
		controllerId, err := models.DecodeHashId(id)
		if err != nil {
			bugsnag.Notify(err)
			return utils.FormatApiError(
				"failed to create models",
				configs.ServerError,
				entities.E{"error": err.Error()},
			)
		}

		controllerIds = append(controllerIds, controllerId)
	}

	err = d.repo.CreateModel(data, controllerIds)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"failed to create models",
			configs.ServerError,
			entities.E{"error": err.Error()},
		)
	}

	return nil
}

func (d DefaultModelService) GetModels(limit, offset int) ([]dto.ModelResponse, *entities.ApiError) {
	getModels, err := d.repo.GetModels(limit, offset)
	if err != nil {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError(
			"failed to retrieve getModels",
			configs.ServerError,
			entities.E{"server": "server error"},
		)
	}

	var mls []dto.ModelResponse
	mls = []dto.ModelResponse{}
	for _, ml := range getModels {
		mRes := ml.ToResponse()
		mls = append(mls, mRes)
	}

	return mls, nil
}

func NewDefaultModelService(repo man_models.ModelRepo) ModelService {
	return &DefaultModelService{repo: repo}
}
