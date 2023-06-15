package services

import (
	"fmt"
	"strings"

	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/domain/taxes"
	"axis/ecommerce-backend/internal/dto"
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/pkg/entities"
	"axis/ecommerce-backend/pkg/utils"

	"github.com/bugsnag/bugsnag-go/v2"
)

type TaxService interface {
	AddTax(request *dto.TaxRequest) *entities.ApiError
	DeleteTaxExemption(request models.TaxExemption) *entities.ApiError
	GetTaxesByAddress(userId uint, address *models.Address) ([]models.Tax, []models.Tax, error)
	GetTaxesForTaxExemptions(userId uint) ([]models.Tax, []models.Tax, error)
	SaveTaxExemption(request models.TaxExemption) *entities.ApiError
}

type DefaultTaxService struct {
	repo taxes.TaxRepo
}

func (d DefaultTaxService) AddTax(request *dto.TaxRequest) *entities.ApiError {
	data := models.Tax{
		Country:     request.Country,
		Description: request.Description,
		Name:        request.Name,
		Rate:        request.Rate,
		State:       request.State,
	}

	err := d.repo.AddTax(data)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError("failed to create tax", configs.ServerError, entities.E{"error": err.Error()})
	}
	return nil
}

func (d DefaultTaxService) GetTaxesByAddress(
	userId uint,
	addresss *models.Address,
) ([]models.Tax, []models.Tax, error) {
	country := strings.ToUpper(addresss.Country)
	state := strings.ToUpper(addresss.State)

	isTaxableCountry := false
	taxableCountries := []string{"CA", "US"}
	for _, v := range taxableCountries {
		if country == v {
			isTaxableCountry = true
			break
		}
	}
	if !isTaxableCountry {
		return nil, nil, nil // No tax for other international countries.
	}

	taxes, taxExemptions, err := d.repo.GetTaxesByAddress(userId, models.Tax{Country: country, State: state})
	if err != nil {
		bugsnag.Notify(err)
		return nil, nil, fmt.Errorf("failed to get taxes by address, error: %s", err.Error())
	}
	return taxes, taxExemptions, nil
}

func (d DefaultTaxService) GetTaxesForTaxExemptions(userId uint) ([]models.Tax, []models.Tax, error) {
	taxes, taxExemptions, err := d.repo.GetTaxesForTaxExemptions(userId)
	if err != nil {
		bugsnag.Notify(err)
		return nil, nil, fmt.Errorf("failed to get taxes by address, error: %s", err.Error())
	}
	return taxes, taxExemptions, nil
}

func (d DefaultTaxService) SaveTaxExemption(request models.TaxExemption) *entities.ApiError {
	if err := d.repo.SaveTaxExemption(request); err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError("failed to save tax exemption", configs.ServerError, entities.E{"error": err.Error()})
	}
	return nil
}

func (d DefaultTaxService) DeleteTaxExemption(request models.TaxExemption) *entities.ApiError {
	if err := d.repo.DeleteTaxExemption(request); err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError("failed to delete tax exemption", configs.ServerError, entities.E{"error": err.Error()})
	}
	return nil
}

func NewDefaultTaxService(repo taxes.TaxRepo) TaxService {
	return &DefaultTaxService{repo: repo}
}
