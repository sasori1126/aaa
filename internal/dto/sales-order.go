package dto

type ControllerOrderRequest struct {
	Controller                   string                    `json:"controller"`
	HeadMake                     string                    `json:"head_make"`
	Carrier                      CarrierRequest            `json:"carrier"`
	CabSystem                    string                    `json:"cab_system"`
	WillInstallGrappleAttachment bool                      `json:"will_install_grapple_attachment"`
	HasHeelRack                  bool                      `json:"has_heel_rack"`
	Joystick                     string                    `json:"joystick"`
	HasModelToTradeIn            bool                      `json:"has_model_to_trade_in"`
	ReportingEmail               string                    `json:"reporting_email"`
	UnitsOfMeasurement           UnitsOfMeasurementRequest `json:"units_of_measurement"`
	CuttingLists                 []CuttingListRequest      `json:"cutting_lists"`
	ContactDetails               ContactDetailsRequest     `json:"contact_details"`
	OfficeDetails                OfficeDetailsRequest      `json:"office_details"`
	BillingContactStaff          ContactStaffRequest       `json:"billing_contact_staff"`
	TechnicalContactStaff        ContactStaffRequest       `json:"technical_contact_staff"`
}

type CarrierRequest struct {
	Year  string `json:"year"`
	Make  string `json:"make"`
	Model string `json:"model"`
}

type CuttingListRequest struct {
	SpeciesName string          `json:"species_name"`
	Presets     []PresetRequest `json:"presets"`
}

type PresetRequest struct {
	TargetLength       int64   `json:"target_length"`
	MinDiameter        int64   `json:"min_diameter"`
	MaxDiameter        float64 `json:"max_diameter"`
	MinDiameterGoToLog float64 `json:"min_diameter_go_to_log"`
	MaxDiameterGoToLog float64 `json:"max_diameter_go_to_log"`
}

type UnitsOfMeasurementRequest struct {
	Length      string `json:"length"`
	Diameter    string `json:"diameter"`
	Volume      string `json:"volume"`
	OilPressure string `json:"oil_pressure"`
	Temperature string `json:"temperature"`
}

type AxisHeadRequest struct {
	HeadName               string                 `json:"head_name"`
	GripType               string                 `json:"grip_type"`
	IncludeSparePartsKit   bool                   `json:"include_spare_parts_kit"`
	IncludeCompleteHoseKit bool                   `json:"include_complete_hose_kit"`
	MeasuringWheel         string                 `json:"measuring_wheel"`
	HasEquipmentDealer     bool                   `json:"has_equipment_dealer"`
	EquipmentDealer        EquipmentDealerRequest `json:"equipment_dealer"`
	ContactDetails         ContactDetailsRequest  `json:"contact_details"`
	OfficeDetails          OfficeDetailsRequest   `json:"office_details"`
	BillingContactStaff    ContactStaffRequest    `json:"billing_contact_staff"`
	TechnicalContactStaff  ContactStaffRequest    `json:"technical_contact_staff"`
}

type KeslaOrderRequest struct {
	HeadName                  string                 `json:"head_name"`
	GripType                  string                 `json:"grip_type"`
	RequirePrologSystem       bool                   `json:"require_prolog_system"`
	IncludeSparePartsKit      bool                   `json:"include_spare_parts_kit"`
	IncludeCompleteHoseKit    bool                   `json:"include_complete_hose_kit"`
	IncludeSpecialToolkit     bool                   `json:"include_special_toolkit"`
	IncludeAuxiliaryCooler    bool                   `json:"include_auxiliary_cooler"`
	IncludeHighPressureFilter bool                   `json:"include_high_pressure_filter"`
	IncludePressurizingValve  bool                   `json:"include_pressurizing_valve"`
	HasEquipmentDealer        bool                   `json:"has_equipment_dealer"`
	EquipmentDealer           EquipmentDealerRequest `json:"equipment_dealer"`
	ContactDetails            ContactDetailsRequest  `json:"contact_details"`
	OfficeDetails             OfficeDetailsRequest   `json:"office_details"`
	BillingContactStaff       ContactStaffRequest    `json:"billing_contact_staff"`
	TechnicalContactStaff     ContactStaffRequest    `json:"technical_contact_staff"`
}

type ContactStaffRequest struct {
	FullNames string `json:"full_names"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type ContactDetailsRequest struct {
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

type EquipmentDealerRequest struct {
	Name             string `json:"name"`
	City             string `json:"city"`
	SalesPersonName  string `json:"sales_person_name"`
	SalesPersonPhone string `json:"sales_person_phone"`
}

type OfficeDetailsRequest struct {
	BusinessName      string `json:"business_name"`
	ContactPersonName string `json:"contact_person_name"`
	OfficeEmail       string `json:"office_email"`
	OfficePhone       string `json:"office_phone"`
}
