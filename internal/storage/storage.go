package storage

import (
	"context"
	"errors"

	"axis/ecommerce-backend/internal/models"
	postgres2 "axis/ecommerce-backend/internal/platform/postgres"
	redis2 "axis/ecommerce-backend/internal/platform/redis"
)

var Cache *redis2.Client

type Storage interface {
	ConnectDb() error
	GetUsers(limit, offset int) ([]models.User, error)
	FindAllOrders() ([]models.CustomerTempRequest, error)
	CreateUser(u models.User) error
	FindUserByField(fv models.FindByField) (*models.User, error)
	UpdateByField(model interface{}, fv map[string]interface{}) error
	CreateCustomerTempReq(cr models.CustomerTempRequest) error
	AddUserAddress(address models.UserAddress) (*models.UserAddress, error)
	GetUserAddresses(userId uint) ([]models.UserAddress, error)
	GetUserAddressByField(fv models.FindByField) (*models.UserAddress, error)
	UpdateUserAddress(data *models.UserAddress) error
	SetUserDefaultAddress(data *models.UserAddress) error
	DelUserAddress(data *models.UserAddress) error
	UpdateUserField(userId uint, data models.FindByField) error
	SearchUsers(limit, offset int, q string) ([]models.User, error)

	// CreateDistributor create a dist
	CreateDistributor(u models.Distributor) error
	GetDistributors(limit, offset int) ([]models.Distributor, error)
	DeleteDistributor(id uint) error
	UpdateDistributor(data *models.Distributor) error
	// FindAllDistributorByField(fv models.FindByField) (*models.Distributor, error)
	FindDistributorByField(fv models.FindByField) (*models.Distributor, error)

	// GetCategories Parts
	GetCategories(limit, offset int) ([]models.DiagramCat, error)
	GetCategoryByField(fv models.FindByField) (*models.DiagramCat, error)
	CreateDiagramCat(data models.DiagramCat) error
	UpdateDiagramCat(data *models.DiagramCat) error
	DeleteCategory(id uint) error
	CreateDiagramSubCat(data models.DiagramSubCat) error

	// CreateController create a new controller
	CreateController(data models.Controller) error
	UpdateController(data *models.Controller) error
	DeleteController(id uint) error
	GetControllers(limit, offset int) ([]models.Controller, error)

	// CreatePart create parts
	CreatePart(data models.Part) error
	UpdatePart(data models.Part) error
	MergeParts(mergeIntoId uint, delIds []uint) error
	GetParts(limit, offset int, returnZeroPrice bool) ([]models.Part, error)
	SearchParts(limit, offset int, q string, returnZeroPrice bool) ([]models.Part, error)
	GetDuplicates(q string) ([]models.Part, error)
	GetPartByField(ctx context.Context, fv models.QueryByField) (*models.Part, error)
	GetPartsByField(fv models.FindByField) ([]models.Part, error)

	// CreateManufacturer create a new manufacturer
	CreateManufacturer(data models.Manufacturer) error
	UpdateManufacturer(data *models.Manufacturer) error
	DeleteManufacturer(id uint) error
	GetManufacturers(limit, offset int) ([]models.Manufacturer, error)

	/*
		Payment transactions
	*/
	AddOrderPayment(data *models.Payment) error

	/*
		Cart transactions
	*/
	CreateCart(data models.Cart) (*models.Cart, error)
	GetUserActiveCart(userId uint) (*models.Cart, error)
	GetCartById(cartId uint) (*models.Cart, error)
	AddCartItem(data models.CartItem) (*models.CartItem, error)
	GetCartItemByField(fv models.FindByField) (*models.CartItem, error)
	GetCartItemsByQuery(qv models.QueryByField) ([]models.CartItem, error)
	GetCarts(limit, offset int, userId uint, status string) ([]models.Cart, error)
	UpdateCartItemQuantity(item *models.CartItem) error
	DeleteCartItem(cartItemId []uint) error
	UpdateCart(cart *models.Cart) error

	/*
		Orders transactions
	*/
	CreateOrder(data models.Order) (*models.Order, error)
	GetUserOrders(userId uint) ([]models.Order, error)
	GetOrders(limit, offset int, userId uint) ([]models.Order, error)
	GetOrderByField(fv models.FindByField) (*models.Order, error)
	UpdateOrderField(orderId uint, fv models.FindByField) error

	// AddEquipment Equipments
	AddEquipment(data models.Equipment) error
	GetEquipmentByQuery(query models.QueryByField) ([]models.Equipment, error)

	// GetModels get Models
	GetModels(limit, offset int) ([]models.Model, error)
	CreateModel(data models.Model, controllerIds []uint) error
	UpdateModel(data *models.Model) error
	GetModel(id uint) (*models.Model, error)
	DeleteModel(id uint) error
	GetSerialByID(id uint) (*models.Serial, error)

	// DeleteImage delete image
	DeleteImage(id uint) error

	// CreateDiagram diagrams
	CreateDiagram(data models.Diagram) error
	DeleteDiagram(id uint) error
	UpdateDiagram(data *models.Diagram) error
	GetDiagrams(limit, offset int, modelId []uint, catId []uint) ([]models.Diagram, error)
	Diagrams(limit, offset int, f bool) ([]models.Diagram, error)
	GetDiagramByField(fv models.FindByField) (*models.Diagram, error)

	/*
		Images and Editor
	*/
	CreateFigureImage(data models.FigureImage) (*models.FigureImage, error)
	UpdateFigImage(data *models.FigureImage) error
	SaveImagesFromCsv(data []models.FigureImage) error
	GetFigImage(id uint) (*models.FigureImage, error)
	DeleteFigImage(id uint) error
	GetFigureImages(limit, offset int, search *string) ([]models.FigureImage, error)

	// GetHeadTypes get head types
	GetHeadTypes(limit, offset int) ([]models.HeadType, error)
	GetHeads(limit, offset int) ([]models.Head, error)
	GetHeadsByType(id uint, limit, offset int) ([]models.Head, error)
	FindHeadByField(fv models.FindByField) (*models.Head, error)
	CreateHeadType(data models.HeadType) error
	CreateHead(data models.Head) error

	/*
		Sales transactions
	*/
	StoreControllerOrder(data models.ControllerOrder) error
	StoreKeslaOrder(data models.KeslaOrder) error
	StoreAxisHeadOrder(data models.AxisHead) error

	/*
		Taxes transactions
	*/
	CreateTax(data models.Tax) error
	DeleteTaxExemption(query models.TaxExemption) error
	GetTaxesByAddress(userId uint, query models.Tax) ([]models.Tax, []models.Tax, error)
	GetTaxesForTaxExemptions(userId uint) ([]models.Tax, []models.Tax, error)
	SaveTaxExemption(query models.TaxExemption) error

	/*
		Dashboard Stats for admin
	*/
	GetStats(limit int) (*models.Stat, error)
}

func ConnectDb() (Storage, error) {
	// Connect to Cache.
	Cache = redis2.NewRedisClient()
	if Cache == nil {
		return nil, errors.New("could not connect to redis")
	}

	// Any db connection integration.
	postgres := postgres2.NewPostgresStorage()
	if err := postgres.ConnectDb(); err != nil {
		return nil, err
	}
	return postgres, nil
}
