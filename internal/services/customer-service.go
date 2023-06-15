package services

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/domain/customers"
	"axis/ecommerce-backend/internal/dto"
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/pkg/entities"
	"axis/ecommerce-backend/pkg/utils"
	"github.com/bugsnag/bugsnag-go/v2"
)

type CustomerService interface {
	CreateTempCustomerReq(r *dto.TempCustomerRequest) *entities.ApiError
	DownloadTempRequests() ([]models.CustomerTempRequest, *entities.ApiError)
}

type DefaultCustomerService struct {
	repo customers.CustomerRepo
}

func (s DefaultCustomerService) CreateTempCustomerReq(r *dto.TempCustomerRequest) *entities.ApiError {
	cr := models.CustomerTempRequest{
		CustomerName:    r.CustomerName,
		CustomerAddress: r.CustomerAddress,
		CustomerContact: r.CustomerContact,
		CustomerEmail:   r.CustomerEmail,
	}
	err := cr.ConvertItemsToJson(r.Items)
	if err != nil {
		bugsnag.Notify(err)
		configs.Logger.Error(err)
		return utils.FormatApiError(
			"Failed to convert request items",
			configs.BadRequest,
			entities.E{"error": "failed to convert items"},
		)
	}
	err = s.repo.CreateCustomerTempReq(cr)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"Failed to save customer request, try again",
			configs.BadRequest,
			entities.E{"error": "failed to convert items"},
		)
	}
	return nil
}

func (s DefaultCustomerService) DownloadTempRequests() ([]models.CustomerTempRequest, *entities.ApiError) {
	orders, err := s.repo.FindAllOrders()
	if err != nil {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError(
			"Failed to get customers orders",
			configs.NotFound,
			entities.E{"orders": "failed to load orders"},
		)
	}

	return orders, nil
}

func NewDefaultCustomerService(repo customers.CustomerRepo) CustomerService {
	return DefaultCustomerService{repo: repo}
}
