package services

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/domain/equipments"
	"axis/ecommerce-backend/internal/dto"
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/pkg/entities"
	"axis/ecommerce-backend/pkg/utils"
	"github.com/bugsnag/bugsnag-go/v2"
)

type EquipmentService interface {
	AddEquipment(userId string, r *dto.EquipmentRequest) *entities.ApiError
	GetUserEquipments(userId string) ([]dto.EquipmentResponse, *entities.ApiError)
}

type DefaultEquipmentService struct {
	repo equipments.EquipmentRepo
}

func (d DefaultEquipmentService) AddEquipment(userHid string, r *dto.EquipmentRequest) *entities.ApiError {
	userId, err := models.DecodeHashId(userHid)
	if err != nil || userId == 0 {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	modelId, err := models.DecodeHashId(r.ModelId)
	if err != nil || modelId == 0 {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	sId, err := models.DecodeHashId(r.SerialId)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	equipment := models.Equipment{
		UserId:   userId,
		ModelId:  modelId,
		SerialId: sId,
	}

	err = d.repo.AddEquipment(equipment)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"failed to add equipment",
			configs.ServerError,
			entities.E{"error": err.Error()},
		)
	}

	return nil
}

func (d DefaultEquipmentService) GetUserEquipments(userHid string) ([]dto.EquipmentResponse, *entities.ApiError) {
	userId, err := models.DecodeHashId(userHid)
	if err != nil || userId == 0 {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	dbRes, err := d.repo.GetEquipmentByField(models.QueryByField{
		Query: "user_id = ?",
		Value: userId,
	})
	if err != nil {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError(
			"failed to get equipments",
			configs.ServerError,
			entities.E{"error": err.Error()},
		)
	}

	var response []dto.EquipmentResponse
	response = []dto.EquipmentResponse{}
	for _, re := range dbRes {
		e := re.ToResponse()
		response = append(response, e)
	}

	return response, nil
}

func NewDefaultEquipmentService(repo equipments.EquipmentRepo) EquipmentService {
	return &DefaultEquipmentService{repo: repo}
}
