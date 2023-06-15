package models

import "axis/ecommerce-backend/configs"

type AxisHead struct {
	configs.GormModel
	HeadName                string          `json:"head_name"`
	GripType                string          `json:"grip_type"`
	IncludeSparePartsKit    bool            `json:"include_spare_parts_kit"`
	IncludeCompleteHoseKit  bool            `json:"include_complete_hose_kit"`
	MeasuringWheel          string          `json:"measuring_wheel"`
	HasEquipmentDealer      bool            `json:"has_equipment_dealer"`
	EquipmentDealer         EquipmentDealer `json:"equipment_dealer"`
	EquipmentDealerId       uint
	ContactDetails          ContactDetails `json:"contact_details"`
	ContactDetailsId        uint
	OfficeDetails           OfficeDetails `json:"office_details"`
	OfficeDetailsId         uint
	BillingContactStaff     ContactStaff `json:"billing_contact_staff"`
	BillingContactStaffId   uint
	TechnicalContactStaff   ContactStaff `json:"technical_contact_staff"`
	TechnicalContactStaffId uint
}
