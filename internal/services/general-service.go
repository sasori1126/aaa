package services

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/domain/sales"
	"axis/ecommerce-backend/internal/dto"
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/internal/notification/mail"
	"axis/ecommerce-backend/pkg/entities"
	"axis/ecommerce-backend/pkg/utils"
	"github.com/bugsnag/bugsnag-go/v2"
)

type GeneralService interface {
	SupportRequest(req *dto.SupportRequest) error
	ControllerSaleOrder(req *dto.ControllerOrderRequest) *entities.ApiError
	KeslaSaleOrder(req *dto.KeslaOrderRequest) *entities.ApiError
	AxisHeadSaleOrder(req *dto.AxisHeadRequest) *entities.ApiError
}

type DefaultGeneralService struct {
	saleRepo sales.SaleRepo
}

func (d DefaultGeneralService) ControllerSaleOrder(req *dto.ControllerOrderRequest) *entities.ApiError {
	var cuttingList []models.CuttingList
	for _, list := range req.CuttingLists {
		var presets []models.Preset
		for _, preset := range list.Presets {
			p := models.Preset{
				TargetLength:       preset.TargetLength,
				MinDiameter:        preset.MinDiameter,
				MaxDiameter:        preset.MaxDiameter,
				MinDiameterGoToLog: preset.MinDiameterGoToLog,
				MaxDiameterGoToLog: preset.MaxDiameterGoToLog,
			}
			presets = append(presets, p)
		}
		cl := models.CuttingList{
			GormModel:   configs.GormModel{},
			SpeciesName: list.SpeciesName,
			Presets:     presets,
		}

		cuttingList = append(cuttingList, cl)
	}
	cm := models.ControllerOrder{
		Controller: req.Controller,
		HeadMake:   req.HeadMake,
		Carrier: models.Carrier{
			Year:  req.Carrier.Year,
			Make:  req.Carrier.Make,
			Model: req.Carrier.Model,
		},
		CabSystem:                    req.CabSystem,
		WillInstallGrappleAttachment: req.WillInstallGrappleAttachment,
		HasHeelRack:                  req.HasHeelRack,
		Joystick:                     req.Joystick,
		HasModelToTradeIn:            req.HasModelToTradeIn,
		ReportingEmail:               req.ReportingEmail,
		UnitsOfMeasurement: models.UnitsOfMeasurement{
			Length:      req.UnitsOfMeasurement.Length,
			Diameter:    req.UnitsOfMeasurement.Diameter,
			Volume:      req.UnitsOfMeasurement.Volume,
			OilPressure: req.UnitsOfMeasurement.OilPressure,
			Temperature: req.UnitsOfMeasurement.Temperature,
		},
		CuttingLists: cuttingList,
		ContactDetails: models.ContactDetails{
			GormModel:    configs.GormModel{},
			FirstName:    req.ContactDetails.FirstName,
			LastName:     req.ContactDetails.LastName,
			BusinessName: req.ContactDetails.BusinessName,
			Email:        req.ContactDetails.Email,
			Phone:        req.ContactDetails.Phone,
			Address:      req.ContactDetails.Address,
			City:         req.ContactDetails.City,
			State:        req.ContactDetails.State,
			Zip:          req.ContactDetails.Zip,
			Country:      req.ContactDetails.Country,
		},
		OfficeDetails: models.OfficeDetails{
			BusinessName:      req.OfficeDetails.BusinessName,
			ContactPersonName: req.OfficeDetails.ContactPersonName,
			OfficeEmail:       req.OfficeDetails.OfficeEmail,
			OfficePhone:       req.OfficeDetails.OfficePhone,
		},
		BillingContactStaff: models.ContactStaff{
			FullNames: req.BillingContactStaff.FullNames,
			Email:     req.BillingContactStaff.Email,
			Phone:     req.BillingContactStaff.Phone,
		},
		TechnicalContactStaff: models.ContactStaff{
			FullNames: req.TechnicalContactStaff.FullNames,
			Email:     req.TechnicalContactStaff.Email,
			Phone:     req.TechnicalContactStaff.Phone,
		},
	}

	err := d.saleRepo.StoreControllerOrder(cm)
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

func (d DefaultGeneralService) KeslaSaleOrder(req *dto.KeslaOrderRequest) *entities.ApiError {
	km := models.KeslaOrder{
		HeadName:                  req.HeadName,
		GripType:                  req.GripType,
		RequirePrologSystem:       req.RequirePrologSystem,
		IncludeSparePartsKit:      req.IncludeSparePartsKit,
		IncludeCompleteHoseKit:    req.IncludeCompleteHoseKit,
		IncludeSpecialToolkit:     req.IncludeSpecialToolkit,
		IncludeAuxiliaryCooler:    req.IncludeAuxiliaryCooler,
		IncludeHighPressureFilter: req.IncludeHighPressureFilter,
		IncludePressurizingValve:  req.IncludePressurizingValve,
		HasEquipmentDealer:        req.HasEquipmentDealer,
		EquipmentDealer: models.EquipmentDealer{
			Name:             req.EquipmentDealer.Name,
			City:             req.EquipmentDealer.City,
			SalesPersonName:  req.EquipmentDealer.SalesPersonName,
			SalesPersonPhone: req.EquipmentDealer.SalesPersonPhone,
		},
		ContactDetails: models.ContactDetails{
			FirstName:    req.ContactDetails.FirstName,
			LastName:     req.ContactDetails.LastName,
			BusinessName: req.ContactDetails.BusinessName,
			Email:        req.ContactDetails.Email,
			Phone:        req.ContactDetails.Phone,
			Address:      req.ContactDetails.Address,
			City:         req.ContactDetails.City,
			State:        req.ContactDetails.State,
			Zip:          req.ContactDetails.Zip,
			Country:      req.ContactDetails.Country,
		},
		OfficeDetails: models.OfficeDetails{
			BusinessName:      req.OfficeDetails.BusinessName,
			ContactPersonName: req.OfficeDetails.ContactPersonName,
			OfficeEmail:       req.OfficeDetails.OfficeEmail,
			OfficePhone:       req.OfficeDetails.OfficePhone,
		},
		BillingContactStaff: models.ContactStaff{
			FullNames: req.BillingContactStaff.FullNames,
			Email:     req.BillingContactStaff.Email,
			Phone:     req.BillingContactStaff.Phone,
		},
		TechnicalContactStaff: models.ContactStaff{
			FullNames: req.TechnicalContactStaff.FullNames,
			Email:     req.TechnicalContactStaff.Email,
			Phone:     req.TechnicalContactStaff.Phone,
		},
	}

	err := d.saleRepo.StoreKeslaOrder(km)
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

func (d DefaultGeneralService) AxisHeadSaleOrder(req *dto.AxisHeadRequest) *entities.ApiError {
	ah := models.AxisHead{
		MeasuringWheel:         req.MeasuringWheel,
		HeadName:               req.HeadName,
		GripType:               req.GripType,
		IncludeSparePartsKit:   req.IncludeSparePartsKit,
		IncludeCompleteHoseKit: req.IncludeCompleteHoseKit,
		HasEquipmentDealer:     req.HasEquipmentDealer,
		EquipmentDealer: models.EquipmentDealer{
			Name:             req.EquipmentDealer.Name,
			City:             req.EquipmentDealer.City,
			SalesPersonName:  req.EquipmentDealer.SalesPersonName,
			SalesPersonPhone: req.EquipmentDealer.SalesPersonPhone,
		},
		ContactDetails: models.ContactDetails{
			FirstName:    req.ContactDetails.FirstName,
			LastName:     req.ContactDetails.LastName,
			BusinessName: req.ContactDetails.BusinessName,
			Email:        req.ContactDetails.Email,
			Phone:        req.ContactDetails.Phone,
			Address:      req.ContactDetails.Address,
			City:         req.ContactDetails.City,
			State:        req.ContactDetails.State,
			Zip:          req.ContactDetails.Zip,
			Country:      req.ContactDetails.Country,
		},
		OfficeDetails: models.OfficeDetails{
			BusinessName:      req.OfficeDetails.BusinessName,
			ContactPersonName: req.OfficeDetails.ContactPersonName,
			OfficeEmail:       req.OfficeDetails.OfficeEmail,
			OfficePhone:       req.OfficeDetails.OfficePhone,
		},
		BillingContactStaff: models.ContactStaff{
			FullNames: req.BillingContactStaff.FullNames,
			Email:     req.BillingContactStaff.Email,
			Phone:     req.BillingContactStaff.Phone,
		},
		TechnicalContactStaff: models.ContactStaff{
			FullNames: req.TechnicalContactStaff.FullNames,
			Email:     req.TechnicalContactStaff.Email,
			Phone:     req.TechnicalContactStaff.Phone,
		},
	}
	err := d.saleRepo.StoreAxisHeadOrder(ah)
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

func (d DefaultGeneralService) SupportRequest(req *dto.SupportRequest) error {
	data := make(map[string]string)
	data["NAME"] = req.Name
	data["EMAIL"] = req.Email + " | " + req.PhoneNumber
	data["REASON"] = req.Reason
	data["MESSAGE"] = req.Message

	emailTo := dto.MailData{
		Name:  "Axis",
		Email: configs.AxisEmail,
	}

	replyTo := dto.MailData{
		Name:  req.Name,
		Email: req.Email,
	}

	err := mail.SendEmail(emailTo, replyTo, data, 23)
	if err != nil {
		bugsnag.Notify(err)
		return err
	}

	return nil
}

func NewDefaultGeneralService(saleRepo sales.SaleRepo) GeneralService {
	return &DefaultGeneralService{saleRepo: saleRepo}
}
