package models

import "axis/ecommerce-backend/configs"

type ControllerOrder struct {
	configs.GormModel
	Controller                   string  `json:"controller"`
	HeadMake                     string  `json:"head_make"`
	Carrier                      Carrier `json:"carrier"`
	CarrierId                    uint
	CabSystem                    string `json:"cab_system"`
	WillInstallGrappleAttachment bool   `json:"will_install_grapple_attachment"`
	HasHeelRack                  bool   `json:"has_heel_rack"`
	Joystick                     string `json:"joystick"`
	HasModelToTradeIn            bool   `json:"has_model_to_trade_in"`
	ReportingEmail               string `json:"reporting_email"`
	UnitsOfMeasurement           UnitsOfMeasurement
	UnitsOfMeasurementId         uint
	CuttingLists                 []CuttingList  `gorm:"many2many:controller_order_cutting_lists;"`
	ContactDetails               ContactDetails `json:"contact_details"`
	ContactDetailsId             uint
	OfficeDetails                OfficeDetails `json:"office_details"`
	OfficeDetailsId              uint
	BillingContactStaff          ContactStaff `json:"billing_contact_staff"`
	BillingContactStaffId        uint
	TechnicalContactStaff        ContactStaff `json:"technical_contact_staff"`
	TechnicalContactStaffId      uint
}

type Carrier struct {
	configs.GormModel
	Year  string `json:"year"`
	Make  string `json:"make"`
	Model string `json:"model"`
}

type CuttingList struct {
	configs.GormModel
	SpeciesName string   `json:"species_name"`
	Presets     []Preset `gorm:"many2many:cutting_list_presets;"`
}

type Preset struct {
	configs.GormModel
	TargetLength       int64   `json:"target_length"`
	MinDiameter        int64   `json:"min_diameter"`
	MaxDiameter        float64 `json:"max_diameter"`
	MinDiameterGoToLog float64 `json:"min_diameter_go_to_log"`
	MaxDiameterGoToLog float64 `json:"max_diameter_go_to_log"`
}

type UnitsOfMeasurement struct {
	configs.GormModel
	Length      string `json:"length"`
	Diameter    string `json:"diameter"`
	Volume      string `json:"volume"`
	OilPressure string `json:"oil_pressure"`
	Temperature string `json:"temperature"`
}
