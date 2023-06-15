package services

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/domain/heads"
	"axis/ecommerce-backend/internal/dto"
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/pkg/entities"
	"axis/ecommerce-backend/pkg/utils"
	"github.com/bugsnag/bugsnag-go/v2"
)

type HeadService interface {
	GetHeads(limit, offset int) ([]dto.HeadResponse, *entities.ApiError)
	GetHeadsType(hid string, limit, offset int) ([]dto.HeadResponse, *entities.ApiError)
	GetHeadById(id string) (*dto.HeadResponse, *entities.ApiError)
	CreateHead(hRequest *dto.HeadRequest, images []models.Image) *entities.ApiError
}

type DefaultHeadService struct {
	repo heads.HeadRepo
}

func (d DefaultHeadService) CreateHead(hRequest *dto.HeadRequest, images []models.Image) *entities.ApiError {
	headTypeId, err := models.DecodeHashId(hRequest.HeadTypeId)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	head := models.Head{
		Name:             hRequest.Name,
		ShortDescription: hRequest.ShortDescription,
		LongDescription:  hRequest.LongDescription,
		Overview:         hRequest.Overview,
		Specification:    hRequest.Specification,
		Images:           images,
		HeadTypeID:       headTypeId,
	}

	err = d.repo.CreateHead(head)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	return nil
}

func (d DefaultHeadService) GetHeadsType(hid string, limit, offset int) ([]dto.HeadResponse, *entities.ApiError) {
	id, err := models.DecodeHashId(hid)
	if err != nil {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	hds, err := d.repo.GetHeadsByType(id, limit, offset)
	if err != nil {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	var headResponses []dto.HeadResponse
	for _, hd := range hds {
		head := hd.ToResponse()
		headResponses = append(headResponses, head)
	}

	return headResponses, nil
}

func (d DefaultHeadService) GetHeadById(hid string) (*dto.HeadResponse, *entities.ApiError) {
	id, err := models.DecodeHashId(hid)
	if err != nil {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	head, err := d.repo.GetHeadById(id)
	if err != nil {
		bugsnag.Notify(err)
		return nil, nil
	}

	var hr dto.HeadResponse = head.ToResponse()

	return &hr, nil
}

func (d DefaultHeadService) GetHeads(limit, offset int) ([]dto.HeadResponse, *entities.ApiError) {
	hds, err := d.repo.GetHeads(limit, offset)
	if err != nil {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	var headResponses []dto.HeadResponse
	for _, hd := range hds {
		head := hd.ToResponse()
		headResponses = append(headResponses, head)
	}

	return headResponses, nil
}

func NewDefaultHeadService(repo heads.HeadRepo) HeadService {
	return DefaultHeadService{repo: repo}
}
