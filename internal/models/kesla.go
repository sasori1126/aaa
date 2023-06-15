package models

import "axis/ecommerce-backend/configs"

type KeslaOrder struct {
	configs.GormModel
	HeadName                  string          `json:"head_name"`
	GripType                  string          `json:"grip_type"`
	RequirePrologSystem       bool            `json:"require_prolog_system"`
	IncludeSparePartsKit      bool            `json:"include_spare_parts_kit"`
	IncludeCompleteHoseKit    bool            `json:"include_complete_hose_kit"`
	IncludeSpecialToolkit     bool            `json:"include_special_toolkit"`
	IncludeAuxiliaryCooler    bool            `json:"include_auxiliary_cooler"`
	IncludeHighPressureFilter bool            `json:"include_high_pressure_filter"`
	IncludePressurizingValve  bool            `json:"include_pressurizing_valve"`
	HasEquipmentDealer        bool            `json:"has_equipment_dealer"`
	EquipmentDealer           EquipmentDealer `json:"equipment_dealer"`
	EquipmentDealerId         uint
	ContactDetails            ContactDetails `json:"contact_details"`
	ContactDetailsId          uint
	OfficeDetails             OfficeDetails `json:"office_details"`
	OfficeDetailsId           uint
	BillingContactStaff       ContactStaff `json:"billing_contact_staff"`
	BillingContactStaffId     uint
	TechnicalContactStaff     ContactStaff `json:"technical_contact_staff"`
	TechnicalContactStaffId   uint
}

type ContactStaff struct {
	configs.GormModel
	FullNames string `json:"full_names"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type ContactDetails struct {
	configs.GormModel
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	BusinessName string `json:"business_name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Address      string `json:"address"`
	City         string `json:"city"`
	State        string `json:"state"`
	Zip          string `json:"zip"`
	Country      string `json:"country"`
}

type EquipmentDealer struct {
	configs.GormModel
	Name             string `json:"name"`
	City             string `json:"city"`
	SalesPersonName  string `json:"sales_person_name"`
	SalesPersonPhone string `json:"sales_person_phone"`
}

type OfficeDetails struct {
	configs.GormModel
	BusinessName      string `json:"business_name"`
	ContactPersonName string `json:"contact_person_name"`
	OfficeEmail       string `json:"office_email"`
	OfficePhone       string `json:"office_phone"`
}
