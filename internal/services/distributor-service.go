package services

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/domain/distributor"
	"axis/ecommerce-backend/internal/dto"
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/pkg/entities"
	"axis/ecommerce-backend/pkg/utils"
	"errors"
	"github.com/bugsnag/bugsnag-go/v2"
	"gorm.io/gorm"
)

type DistributorService interface {
	CreateDistributor(req *dto.DistributorRequest) *entities.ApiError
	EditDistributor()
	DeleteDistributor()
	GetDistributorById(hid string) (*dto.DistributorResponse, *entities.ApiError)
	UpdateDistributorById(hid string, req *dto.DistributorRequest) *entities.ApiError
	DeleteDistributorById(hid string) *entities.ApiError
	GetDistributors(limit, offset int) ([]dto.DistributorResponse, *entities.ApiError)
}

type DefaultDistributorService struct {
	repo distributor.RepoDistributor
}

func (d DefaultDistributorService) UpdateDistributorById(hid string, req *dto.DistributorRequest) *entities.ApiError {
	id, err := models.DecodeHashId(hid)
	if err != nil || id == 0 {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	dist, err := d.repo.GetDistributor(id)
	if err != nil {
		bugsnag.Notify(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.FormatApiError(
				"Distributor not found",
				configs.NotFound,
				entities.E{},
			)
		}

		return utils.FormatApiError(
			"Failed to process error "+err.Error(),
			configs.ServerError,
			entities.E{},
		)
	}

	data := models.Distributor{
		Name: req.Name,
		Address: models.Address{
			City:       req.Address.City,
			StreetName: req.Address.StreetName,
			Country:    req.Address.Country,
			Province:   req.Address.Province,
			State:      req.Address.State,
			ZipCode:    req.Address.ZipCode,
			Lat:        req.Address.Lat,
			Lon:        req.Address.Lon,
		},
		Contact: models.Contact{
			Name:   req.Contact.Name,
			Email:  req.Contact.Email,
			Phone:  req.Contact.Phone,
			Fax:    req.Contact.Fax,
			Postal: req.Contact.Postal,
		},
		Site: req.Site,
	}
	data.ID = dist.ID

	err = d.repo.UpdateDistributor(&data)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"distributor_id": "failed to update"},
		)
	}
	return nil
}

func (d DefaultDistributorService) DeleteDistributorById(hid string) *entities.ApiError {
	id, err := models.DecodeHashId(hid)
	if err != nil || id == 0 {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}
	err = d.repo.DeleteDistributor(id)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"distributor_id": "failed to delete distributor"},
		)
	}
	return nil
}

func (d DefaultDistributorService) GetDistributorById(hid string) (*dto.DistributorResponse, *entities.ApiError) {
	id, err := models.DecodeHashId(hid)
	if err != nil {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	dist, err := d.repo.GetDistributor(id)
	if err != nil {
		bugsnag.Notify(err)
		return nil, nil
	}

	var dr dto.DistributorResponse = dist.ToResponse()

	return &dr, nil
}

func (d DefaultDistributorService) CreateDistributor(req *dto.DistributorRequest) *entities.ApiError {
	data := models.Distributor{
		Name: req.Name,
		Address: models.Address{
			City:       req.Address.City,
			StreetName: req.Address.StreetName,
			Country:    req.Address.Country,
			Province:   req.Address.Province,
			State:      req.Address.State,
			ZipCode:    req.Address.ZipCode,
			Lat:        req.Address.Lat,
			Lon:        req.Address.Lon,
		},
		Contact: models.Contact{
			Name:   req.Contact.Name,
			Email:  req.Contact.Email,
			Phone:  req.Contact.Phone,
			Fax:    req.Contact.Fax,
			Postal: req.Contact.Postal,
		},
		Site: req.Site,
	}
	err := d.repo.CreateDistributor(data)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"failed to create create",
			configs.ServerError,
			entities.E{"error": err.Error()},
		)
	}

	return nil
}

func (d DefaultDistributorService) EditDistributor() {
	panic("implement me")
}

func (d DefaultDistributorService) DeleteDistributor() {
	panic("implement me")
}

func (d DefaultDistributorService) GetDistributors(limit, offset int) ([]dto.DistributorResponse, *entities.ApiError) {
	distributors, err := d.repo.GetDistributors(limit, offset)
	if err != nil {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	var dits []dto.DistributorResponse
	for _, m := range distributors {
		dt := m.ToResponse()
		dits = append(dits, dt)
	}

	return dits, nil
}

func NewDistributorService(repo distributor.RepoDistributor) DistributorService {
	return DefaultDistributorService{
		repo: repo,
	}
}
