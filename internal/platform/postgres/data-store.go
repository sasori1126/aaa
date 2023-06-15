package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"

	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal"
	"axis/ecommerce-backend/internal/dto"
	"axis/ecommerce-backend/internal/models"

	"go.opentelemetry.io/otel"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

func (s *Storage) SaveImagesFromCsv(data []models.FigureImage) error {
	return s.db.Save(data).Error
}

func (s *Storage) UpdateFigImage(data *models.FigureImage) error {
	return s.db.Save(data).Error
}

func (s *Storage) GetFigImage(id uint) (*models.FigureImage, error) {
	image := models.FigureImage{}
	err := s.db.Where("id = ?", id).First(&image).Error
	if err != nil {
		return nil, err
	}

	return &image, nil
}

func (s *Storage) DeleteFigImage(id uint) error {
	return s.db.Delete(&models.FigureImage{}, id).Error
}

func (s *Storage) GetSerialByID(id uint) (*models.Serial, error) {
	serial := models.Serial{}
	err := s.db.Where("id = ?", id).First(&serial).Error
	if err != nil {
		return nil, err
	}

	return &serial, nil
}

func (s *Storage) UpdateDiagram(data *models.Diagram) error {
	d := &models.Diagram{
		Name:         data.Name,
		Description:  data.Description,
		Status:       data.Status,
		BgImage:      data.BgImage,
		Thumbnail:    data.Thumbnail,
		Draft:        data.Draft,
		ControllerId: data.ControllerId,
	}

	d.ID = data.ID
	err := s.db.Save(d).Error
	if err != nil {
		return err
	}

	err = s.db.Model(d).Association("Models").Replace(data.Models)
	if err != nil {
		return err
	}

	err = s.db.Model(d).Association("Cats").Replace(data.Cats)
	if err != nil {
		return err
	}

	dps, _, err := s.sortOutParts(d.ID, data.Parts)
	if err != nil {
		return err
	}

	err = s.db.Where("diagram_id = ?", d.ID).Delete(&models.DiagramPart{}).Error
	if err != nil {
		return err
	}

	if len(dps) != 0 {
		err = s.db.Table("diagram_parts").Create(dps).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) GetFigureImages(limit, offset int, search *string) ([]models.FigureImage, error) {
	var fis []models.FigureImage
	query := s.db.Order("created_at desc").Limit(limit).Offset(offset)
	if search != nil {
		query = query.Where("title ILike ?", "%"+*search+"%")
	}

	err := query.Find(&fis).Error
	if err != nil {
		return nil, err
	}

	return fis, nil
}

func (s *Storage) CreateFigureImage(data models.FigureImage) (*models.FigureImage, error) {
	dd := data
	err := s.db.Create(&dd).Error
	if err != nil {
		return nil, err
	}

	return &dd, nil
}

func (s *Storage) DeleteImage(id uint) error {
	img := &models.Image{}
	err := s.db.Where("id = ?", id).First(img).Error
	if err != nil {
		return err
	}

	err = s.db.Model(img).Association("Heads").Delete(&models.Image{}, id)
	if err != nil {
		return err
	}
	err = s.db.Model(img).Association("Parts").Delete(&models.Image{}, id)
	if err != nil {
		return err
	}

	return s.db.Delete(&models.Image{}, id).Error
}

func (s *Storage) UpdatePart(data models.Part) error {
	return s.db.Save(data).Error
}

func (s *Storage) GetStats(limit int) (*models.Stat, error) {
	stat := models.Stat{}
	var countOrders int64
	err := s.db.Model(&models.Order{}).Count(&countOrders).Error
	if err != nil {
		return nil, err
	}
	stat.TotalOrders = countOrders

	var countUsers int64
	err = s.db.Model(&models.User{}).Count(&countUsers).Error
	if err != nil {
		return nil, err
	}
	stat.TotalUsers = countUsers

	var countActiveCarts int64
	err = s.db.Model(&models.Cart{}).Where("status = ?", "active").Count(&countActiveCarts).Error
	if err != nil {
		return nil, err
	}
	stat.ActiveCarts = countActiveCarts

	var countParts int64
	err = s.db.Model(&models.Part{}).Count(&countParts).Error
	if err != nil {
		return nil, err
	}
	stat.TotalParts = countParts
	var orders []models.Order
	err = s.db.Limit(limit).Where("status in ?", []string{
		configs.OrderPaid,
		configs.OrderPending,
		configs.OrderPendingOnAccount,
		configs.OrderPendingOnPhone,
	}).Preload("Items").
		Preload("User").
		Order("id desc").Find(&orders).Error
	if err != nil {
		return nil, err
	}

	var pendingOrders []dto.OrderResponse
	pendingOrders = []dto.OrderResponse{}
	var ordersPendingDelivery []dto.OrderResponse
	ordersPendingDelivery = []dto.OrderResponse{}
	for _, order := range orders {
		o := order.ToResponse()
		if o.Status == configs.OrderPaid {
			ordersPendingDelivery = append(ordersPendingDelivery, o)
		} else {
			pendingOrders = append(pendingOrders, o)
		}
	}
	stat.PendingOrders = pendingOrders
	stat.UndeliveredOrders = ordersPendingDelivery

	return &stat, nil
}

func (s *Storage) MergeParts(mergeIntoId uint, delIds []uint) error {
	// update order items
	err := s.db.Model(&models.OrderItem{}).
		Where("part_id IN (?)", delIds).
		Update("part_id", mergeIntoId).Error
	if err != nil {
		return err
	}

	// update cart items
	err = s.db.Model(&models.CartItem{}).
		Where("part_id IN (?)", delIds).
		Update("part_id", mergeIntoId).Error
	if err != nil {
		return err
	}

	// Update diagram parts
	err = s.db.Model(&models.DiagramPart{}).
		Where("part_id IN (?)", delIds).
		Update("part_id", mergeIntoId).Error
	if err != nil {
		return err
	}

	// Update part images
	err = s.db.Table("part_images").
		Where("part_id IN (?)", delIds).
		Update("part_id", mergeIntoId).Error
	if err != nil {
		return err
	}

	// Update part models
	err = s.db.Table("part_models").
		Where("part_id IN (?)", delIds).
		Update("part_id", mergeIntoId).Error
	if err != nil {
		return err
	}

	// Update parts part types
	err = s.db.Table("part_part_types").
		Where("part_id IN (?)", delIds).
		Update("part_id", mergeIntoId).Error
	if err != nil {
		return err
	}

	err = s.db.Where("id IN (?)", delIds).
		Delete(&models.Part{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetCarts(limit, offset int, userId uint, status string) ([]models.Cart, error) {
	var carts []models.Cart
	query := s.db.Order("created_at desc").Limit(limit).Offset(offset).
		Preload("Items.Part").
		Preload("User")
	if userId != 0 {
		query = query.Where("user_id = ?", userId)
	}

	if status != "" {
		query = query.Where("status = ? ", status)
	}

	err := query.Find(&carts).Error
	if err != nil {
		return nil, err
	}
	return carts, nil
}

func (s *Storage) GetPartsByField(fv models.FindByField) ([]models.Part, error) {
	value := fv.Value.(string)
	var parts []models.Part
	err := s.db.Where("part_number SIMILAR TO ?", value).Find(&parts).Error
	if err != nil {
		return nil, err
	}
	return parts, nil
}

func (s *Storage) SearchUsers(limit, offset int, q string) ([]models.User, error) {
	var users []models.User
	res := s.db.Order("name").Where("name ILike ?", "%"+q+"%").
		Or("business_name ILIKE ?", "%"+q+"%").
		Or("email ILIKE ?", "%"+q+"%").
		Limit(limit).Offset(offset).
		Find(&users)

	if res.Error != nil {
		return nil, res.Error
	}

	return users, nil
}

func (s *Storage) UpdateModel(data *models.Model) error {
	return s.db.Save(data).Error
}

func (s *Storage) GetModel(id uint) (*models.Model, error) {
	mdl := &models.Model{}
	db := s.db.Where("id = ?", id).
		Preload("Manufacturer").
		Preload("Controllers").
		Preload("Serials").First(mdl)
	if db.Error != nil {
		return nil, db.Error
	}

	return mdl, nil
}

func (s *Storage) DeleteModel(id uint) error {
	md := &models.Model{}
	err := s.db.Where("id = ?", id).
		Preload("Manufacturer").
		Preload("Controllers").
		Preload("Serials").First(md).Error
	if err != nil {
		return err
	}

	err = s.db.Model(md).Association("Controllers").Delete(md.Controllers)
	if err != nil {
		return err
	}
	err = s.db.Model(md).Association("Manufacturer").Delete(md.Manufacturer)
	if err != nil {
		return err
	}

	return s.db.Select("Serials").Delete(&models.Model{}, id).Error
}

func (s *Storage) UpdateController(data *models.Controller) error {
	return s.db.Save(data).Error
}

func (s *Storage) DeleteController(id uint) error {
	return s.db.Delete(&models.Controller{}, id).Error
}

func (s *Storage) UpdateDiagramCat(data *models.DiagramCat) error {
	return s.db.Save(data).Error
}

func (s *Storage) DeleteCategory(id uint) error {
	return s.db.Delete(&models.DiagramCat{}, id).Error
}

func (s *Storage) Diagrams(limit, offset int, f bool) ([]models.Diagram, error) {
	var diagrams []models.Diagram
	res := s.db.Order("name").Limit(limit).Offset(offset)

	if f {
		res = res.Where("draft != ? ", "")
	}

	res = res.Find(&diagrams)

	if res.Error != nil {
		return nil, res.Error
	}

	return diagrams, nil
}

func (s *Storage) DeleteManufacturer(id uint) error {
	return s.db.Delete(&models.Manufacturer{}, id).Error
}

func (s *Storage) UpdateManufacturer(data *models.Manufacturer) error {
	return s.db.Save(data).Error
}

func (s *Storage) UpdateDistributor(data *models.Distributor) error {
	return s.db.Save(data).Error
}

func (s *Storage) DeleteDistributor(id uint) error {
	return s.db.Delete(&models.Distributor{}, id).Error
}

func (s *Storage) UpdateUserField(userId uint, data models.FindByField) error {
	return s.db.Model(&models.User{}).
		Where("id = ?", userId).
		Update(data.Field, data.Value).Error
}

func (s *Storage) UpdateOrderField(orderId uint, fv models.FindByField) error {
	return s.db.Model(&models.Order{}).Where("id = ?", orderId).Update(fv.Field, fv.Value).Error
}

func (s *Storage) GetOrderByField(fv models.FindByField) (*models.Order, error) {
	order := &models.Order{}
	err := s.db.Where(fv.Field+" = ?", fv.Value).
		Preload("BillingAddress.Address").
		Preload("BillingAddress.Contact").
		Preload("ShippingAddress.Address").
		Preload("ShippingAddress.Contact").
		Preload("Items.Part").
		Preload("User").
		Preload("Taxes").
		Find(order).Error
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *Storage) GetOrders(limit, offset int, userId uint) ([]models.Order, error) {
	var orders []models.Order
	query := s.db.Order("created_at desc").Limit(limit).Offset(offset).
		Preload("BillingAddress.Address").
		Preload("BillingAddress.Contact").
		Preload("ShippingAddress.Address").
		Preload("ShippingAddress.Contact").
		Preload("Items.Part").
		Preload("User").
		Preload("Taxes")
	if userId != 0 {
		query = query.Where("user_id = ?", userId)
	}

	err := query.Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *Storage) StoreControllerOrder(data models.ControllerOrder) error {
	dd := data
	return s.db.Create(&dd).Error
}

func (s *Storage) StoreKeslaOrder(data models.KeslaOrder) error {
	dd := data
	return s.db.Create(&dd).Error
}

func (s *Storage) StoreAxisHeadOrder(data models.AxisHead) error {
	dd := data
	return s.db.Create(&dd).Error
}

func (s *Storage) GetUserOrders(userId uint) ([]models.Order, error) {
	var orders []models.Order
	dbRes := s.db.Order("created_at desc").Where("user_id = ?", userId).
		Preload("BillingAddress.Address").
		Preload("BillingAddress.Contact").
		Preload("ShippingAddress.Address").
		Preload("ShippingAddress.Contact").
		Preload("Items.Part").
		Preload("Taxes").
		Find(&orders)
	if dbRes.Error != nil {
		return nil, dbRes.Error
	}

	return orders, nil
}

func (s *Storage) CreateOrder(data models.Order) (*models.Order, error) {
	orderAvailable := models.Order{}
	dbRes := s.db.Where("cart_id = ?", data.CartId).First(&orderAvailable)
	if !errors.Is(dbRes.Error, gorm.ErrRecordNotFound) {
		return nil, configs.ErrOrderAlreadyExist
	}
	orderAvailable = data
	err := s.db.Create(&orderAvailable)
	return &orderAvailable, err.Error
}

func (s *Storage) AddOrderPayment(data *models.Payment) error {
	payment := data
	err := s.db.Create(&payment)
	return err.Error
}

func (s *Storage) GetCartById(cartId uint) (*models.Cart, error) {
	cart := &models.Cart{}
	res := s.db.Where("id = ?", cartId).
		Preload("Items.Part.Images").
		Preload("User").
		First(cart)
	return cart, res.Error
}

func (s *Storage) DelUserAddress(data *models.UserAddress) error {
	dbRes := s.db.Delete(data)
	return dbRes.Error
}

func (s *Storage) SetUserDefaultAddress(data *models.UserAddress) error {
	dbRes := s.db.Model(&models.UserAddress{}).
		Where("is_default_address = ?", true).
		Update("is_default_address", false)
	if dbRes.Error != nil {
		return dbRes.Error
	}

	dbRes = s.db.Save(data)
	return dbRes.Error
}

func (s *Storage) UpdateUserAddress(data *models.UserAddress) error {
	if data.IsDefaultAddress {
		return s.db.Model(&models.UserAddress{}).
			Where("is_default_address = ?", true).
			Update("is_default_address", false).Error
	}

	return s.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(data).Error
}

func (s *Storage) GetCartItemsByQuery(qv models.QueryByField) ([]models.CartItem, error) {
	var items []models.CartItem
	res := s.db.Where(qv.Query, qv.Value).Find(&items)
	if res.Error != nil {
		return nil, res.Error
	}

	return items, nil
}

func (s *Storage) UpdateCart(cart *models.Cart) error {
	return s.db.Save(cart).Error
}

func (s *Storage) GetTaxes(limit, offset int) ([]models.Tax, error) {
	var taxes []models.Tax
	res := s.db.Limit(limit).Offset(offset).Find(&taxes)
	if res.Error != nil {
		return nil, res.Error
	}

	return taxes, nil
}

func (s *Storage) GetTaxesByQuery(qv models.QueryByField) ([]models.Tax, error) {
	var taxes []models.Tax
	res := s.db.Where(qv.Query, qv.Value).Find(&taxes)
	if res.Error != nil {
		return nil, res.Error
	}
	return taxes, nil
}

func (s *Storage) GetTaxesByAddress(userId uint, query models.Tax) ([]models.Tax, []models.Tax, error) {
	var taxes []models.Tax
	taxResult := s.db.Where(&query).Order("name ASC").Find(&taxes)
	if taxResult.Error != nil {
		return nil, nil, taxResult.Error
	}

	var taxExemptions []models.TaxExemption
	taxExemptionResult := s.db.Where(&models.TaxExemption{UserId: userId}).Find(&taxExemptions)
	if taxExemptionResult.Error != nil {
		return nil, nil, taxExemptionResult.Error
	}
	if taxExemptions == nil || len(taxExemptions) == 0 {
		return taxes, nil, nil
	}

	var exemptedTaxes []models.Tax
	var filteredTaxes []models.Tax
	for i := range taxes {
		tax := taxes[i]
		isTaxExempted := false
		for j := range taxExemptions {
			if tax.ID == taxExemptions[j].TaxId {
				exemptedTaxes = append(exemptedTaxes, tax)
				isTaxExempted = true
				continue
			}
		}
		if !isTaxExempted {
			filteredTaxes = append(filteredTaxes, tax)
		}
	}
	return filteredTaxes, exemptedTaxes, nil
}

func (s *Storage) GetTaxesForTaxExemptions(userId uint) ([]models.Tax, []models.Tax, error) {
	var taxes []models.Tax
	query := models.Tax{IsAllowedTaxExemption: true}
	taxResult := s.db.Where(&query).Order("name ASC").Find(&taxes)
	if taxResult.Error != nil {
		return nil, nil, taxResult.Error
	}

	var taxExemptions []models.TaxExemption
	taxExemptionResult := s.db.Where(&models.TaxExemption{UserId: userId}).Find(&taxExemptions)
	if taxExemptionResult.Error != nil {
		return nil, nil, taxExemptionResult.Error
	}
	if taxExemptions == nil || len(taxExemptions) == 0 {
		return taxes, nil, nil
	}

	var exemptedTaxes []models.Tax
	var filteredTaxes []models.Tax
	for i := range taxes {
		tax := taxes[i]
		isTaxExempted := false
		for j := range taxExemptions {
			if tax.ID == taxExemptions[j].TaxId {
				exemptedTaxes = append(exemptedTaxes, tax)
				isTaxExempted = true
				continue
			}
		}
		if !isTaxExempted {
			filteredTaxes = append(filteredTaxes, tax)
		}
	}
	return filteredTaxes, exemptedTaxes, nil
}

func (s *Storage) SaveTaxExemption(query models.TaxExemption) error {
	res := s.db.Save(&query)
	return res.Error
}

func (s *Storage) DeleteTaxExemption(query models.TaxExemption) error {
	res := s.db.Where("user_id = ? AND tax_id = ?", query.UserId, query.TaxId).Unscoped().Delete(&models.TaxExemption{})
	return res.Error
}

func (s *Storage) CreateTax(data models.Tax) error {
	tax := data
	err := s.db.Create(&tax)
	return err.Error
}

func (s *Storage) GetUserAddressByField(fv models.FindByField) (*models.UserAddress, error) {
	address := &models.UserAddress{}
	dbRes := s.db.Where(fv.Field+" = ?", fv.Value).
		Preload("Address").
		Preload("Contact").
		First(address)
	if dbRes.Error != nil {
		return nil, dbRes.Error
	}

	return address, nil
}

func (s *Storage) GetUserAddresses(userId uint) ([]models.UserAddress, error) {
	var addresses []models.UserAddress
	dbRes := s.db.Where("user_id = ?", userId).
		Preload("Address").
		Preload("Contact").
		Find(&addresses)
	if dbRes.Error != nil {
		return nil, dbRes.Error
	}

	return addresses, nil
}

func (s *Storage) AddUserAddress(address models.UserAddress) (*models.UserAddress, error) {
	dbRes := s.db.Model(&models.UserAddress{}).
		Where("is_default_address = ?", true).
		Update("is_default_address", false)
	if dbRes.Error != nil {
		return nil, dbRes.Error
	}
	userAddress := address
	err := s.db.Create(&userAddress)
	return &userAddress, err.Error
}

func (s *Storage) GetCartItemByField(fv models.FindByField) (*models.CartItem, error) {
	ci := &models.CartItem{}
	dbRes := s.db.Where(fv.Field+" = ?", fv.Value).Preload("Cart").First(ci)
	if dbRes.Error != nil {
		return nil, dbRes.Error
	}

	return ci, nil
}

func (s *Storage) UpdateCartItemQuantity(item *models.CartItem) error {
	res := s.db.Save(item)
	return res.Error
}

func (s *Storage) DeleteCartItem(cartItemId []uint) error {
	res := s.db.Where("id IN (?)", cartItemId).Delete(&models.CartItem{})
	return res.Error
}

func (s *Storage) AddCartItem(data models.CartItem) (*models.CartItem, error) {
	c := models.CartItem{}
	res := s.db.Where("cart_id = ?", data.Cart.ID).
		Where("part_id = ?", data.PartId).First(&c)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		c = data
		err := s.db.Create(&c)
		return &c, err.Error
	}

	c.Amount += data.Amount
	c.Quantity += data.Quantity
	c.Weight += data.Weight
	res = s.db.Save(&c)
	return &c, res.Error
}

func (s *Storage) GetUserActiveCart(userId uint) (*models.Cart, error) {
	cart := &models.Cart{}
	res := s.db.Where("user_id = ?", userId).
		Where("status = ?", internal.CartStatusActive).
		Preload("Items.Part.Images").
		Preload("User").
		First(cart)
	return cart, res.Error
}

func (s *Storage) CreateCart(data models.Cart) (*models.Cart, error) {
	c := data
	err := s.db.Create(&c)
	return &c, err.Error
}

func (s *Storage) AddEquipment(data models.Equipment) error {
	err := s.db.Create(&data)
	return err.Error
}

func (s *Storage) GetEquipmentByQuery(query models.QueryByField) ([]models.Equipment, error) {
	var eqs []models.Equipment
	res := s.db.Where(query.Query, query.Value).
		Preload("Model").Preload("Serial").
		Find(&eqs)

	if res.Error != nil {
		return nil, res.Error
	}

	return eqs, res.Error
}

func (s *Storage) SearchParts(limit, offset int, q string, returnZeroPrice bool) ([]models.Part, error) {
	var parts []models.Part
	searchQuery := s.db.Order("name").Where("name ILike ?", "%"+q+"%").
		Or("part_number ILIKE ?", "%"+q+"%").
		Or("oem_number ILIKE ?", "%"+q+"%").
		Or("description ILIKE ?", "%"+q+"%").
		Or("description ILIKE ?", "%"+q+"%").
		Preload("Images").
		Preload("Types").
		Limit(limit).Offset(offset)

	if !returnZeroPrice {
		searchQuery = searchQuery.Where("price > ?", 0)
	}
	err := searchQuery.Find(&parts).Error
	if err != nil {
		return nil, err
	}

	return parts, nil
}

func (s *Storage) GetDuplicates(q string) ([]models.Part, error) {
	var parts []models.Part
	res := s.db.Order("name").
		Where("part_number SIMILAR TO ?", q).
		Find(&parts)

	if res.Error != nil {
		return nil, res.Error
	}

	return parts, nil
}

func (s *Storage) sortOutParts(diagramId uint, ps []models.Part) ([]models.DiagramPart, []uint, error) {
	var partsId []uint
	type ss struct {
		No             int
		RecommendedQty int
		CatNote        string
	}

	inResult := make(map[uint]ss)
	for _, part := range ps {
		inResult[part.ID] = ss{
			No:             part.PartDiagramNumber,
			RecommendedQty: part.RecommendedQty,
			CatNote:        part.CatNote,
		}
		partsId = append(partsId, part.ID)
	}

	var parts []models.Part
	res := s.db.Where("id IN (?)", partsId).Find(&parts)
	if res.Error != nil {
		return nil, nil, res.Error
	}

	dps := make([]models.DiagramPart, 0, len(parts))
	for _, p := range parts {
		dp := models.DiagramPart{
			DiagramId:         diagramId,
			PartId:            p.ID,
			PartDiagramNumber: inResult[p.ID].No,
			PartNumber:        p.PartNumber,
			PartDescription:   p.Description,
			RecommendedQty:    inResult[p.ID].RecommendedQty,
			CatNote:           inResult[p.ID].CatNote,
		}
		dps = append(dps, dp)
	}
	if res.Error != nil {
		return nil, nil, res.Error
	}

	return dps, partsId, nil
}

func (s *Storage) DeleteDiagram(id uint) error {
	err := s.db.Table("diagram_parts").Where("diagram_id = ?", id).Delete(&models.DiagramPart{}).Error
	if err != nil {
		return err
	}

	err = s.db.Exec(fmt.Sprintf("delete from diagram_models where diagram_id=%d", id)).Error
	if err != nil {
		return err
	}

	err = s.db.Exec(fmt.Sprintf("delete from diagram_models where diagram_id=%d", id)).Error
	if err != nil {
		return err
	}

	err = s.db.Exec(fmt.Sprintf("delete from diagram_diagram_cats where diagram_id=%d", id)).Error
	if err != nil {
		return err
	}

	err = s.db.Exec(fmt.Sprintf("delete from diagram_images where diagram_id=%d", id)).Error
	if err != nil {
		return err
	}

	err = s.db.Exec(fmt.Sprintf("delete from diagram_diagram_sub_cats where diagram_id=%d", id)).Error
	if err != nil {
		return err
	}

	var d models.Diagram
	err = s.db.Where("id = ?", id).First(&d).Error
	if err != nil {
		return err
	}

	return s.db.Delete(&d).Error
}

func (s *Storage) CreateDiagram(data models.Diagram) error {
	d := data
	d.Parts = []models.Part{}

	err := s.db.Create(&d).Error
	if err != nil {
		return err
	}

	dps, _, err := s.sortOutParts(d.ID, data.Parts)
	if err != nil {
		return err
	}

	if len(dps) > 0 {
		return s.db.Table("diagram_parts").Create(dps).Error
	}

	return nil
}

func (s *Storage) GetParts(limit, offset int, returnZeroPrice bool) ([]models.Part, error) {
	var getParts []models.Part
	getPartsQuery := s.db.Order("name").Limit(limit).Offset(offset).
		Preload("Images").
		Preload("Types")

	if !returnZeroPrice {
		getPartsQuery = getPartsQuery.Where("price > ?", 0)
	}

	err := getPartsQuery.Find(&getParts).Error
	if err != nil {
		return nil, err
	}

	return getParts, nil
}

func (s *Storage) GetPartByField(ctx context.Context, fv models.QueryByField) (*models.Part, error) {
	_, span := otel.Tracer("").Start(ctx, "getPartByFieldDbQuery")
	defer span.End()
	part := &models.Part{}
	db := s.db.Where(fv.Query, fv.Value).
		Preload("Types").
		Preload("Images").
		First(part)
	if db.Error != nil {
		return nil, db.Error
	}
	return part, nil
}

func (s *Storage) CreatePart(data models.Part) error {
	err := s.db.Create(&data)
	return err.Error
}

func (s *Storage) GetDiagramByField(fv models.FindByField) (*models.Diagram, error) {
	diagram := &models.Diagram{}
	stringValue := fmt.Sprintf("%v", fv.Value)
	db := s.db.Where(fv.Field+" = ?", fv.Value).
		Preload("Parts", func(db *gorm.DB) *gorm.DB {
			return db.Select("parts.*, diagram_parts.part_diagram_number").
				Joins("LEFT JOIN diagram_parts ON diagram_parts.diagram_id = " + stringValue + " AND diagram_parts.part_id = parts.id")
		}).
		Preload("Parts.Images").
		Preload("Parts.Types").
		Preload("Controller").
		Preload("Models").
		Preload("Cats").
		Preload("SubCats").
		First(diagram)
	if db.Error != nil {
		return nil, db.Error
	}

	return diagram, nil
}

func (s *Storage) GetDiagrams(limit, offset int, modelId []uint, catId []uint) ([]models.Diagram, error) {
	var diagrams []models.Diagram
	var mt []struct {
		DiagramId uint
		ModelId   uint
	}
	res := s.db.Table("diagram_models").Where("model_id IN ?", modelId).Find(&mt)
	if res.Error != nil {
		return nil, res.Error
	}

	log.Println(mt, modelId, "mdls")

	var ct []struct {
		DiagramId    uint
		DiagramCatId uint
	}
	res = s.db.Table("diagram_diagram_cats").Where("diagram_cat_id IN ?", catId).Find(&ct)
	if res.Error != nil {
		return nil, res.Error
	}

	var diagramIds []uint
	inResult := make(map[uint]bool)
	for _, d := range ct {
		if _, ok := inResult[d.DiagramId]; !ok {
			getId := d.DiagramId
			inResult[getId] = true
		}
	}

	for _, d := range mt {
		getId := d.DiagramId
		if _, ok := inResult[d.DiagramId]; ok && len(ct) > 0 {
			diagramIds = append(diagramIds, getId)
		} else if len(ct) == 0 {
			diagramIds = append(diagramIds, getId)
		}
	}

	res = s.db.Order("name").Where("id IN (?)", diagramIds).Where("status = ? ", "Active").Limit(limit).Offset(offset).Find(&diagrams)

	if res.Error != nil {
		return nil, res.Error
	}

	return diagrams, nil
}

func (s *Storage) GetCategoryByField(fv models.FindByField) (*models.DiagramCat, error) {
	cat := &models.DiagramCat{}
	db := s.db.Where(fv.Field+" = ?", fv.Value).Preload("Diagrams").First(cat)
	if db.Error != nil {
		return nil, db.Error
	}
	return cat, nil
}

func (s *Storage) GetControllers(limit, offset int) ([]models.Controller, error) {
	var getControllers []models.Controller
	res := s.db.Order("name").Limit(limit).Offset(offset).Find(&getControllers)
	if res.Error != nil {
		return nil, res.Error
	}

	return getControllers, nil
}

func (s *Storage) CreateController(data models.Controller) error {
	err := s.db.Create(&data)
	return err.Error
}

func (s *Storage) CreateModel(data models.Model, controllerIds []uint) error {
	err := s.db.Create(&data)

	if err.Error != nil {
		return err.Error
	}

	var modelControllers []models.ModelController
	for _, controllerId := range controllerIds {
		mc := models.ModelController{ModelId: data.ID, ControllerId: controllerId}
		modelControllers = append(modelControllers, mc)
	}

	dbErr := s.db.Table("model_controllers").Create(&modelControllers)
	return dbErr.Error
}

func (s *Storage) GetManufacturers(limit, offset int) ([]models.Manufacturer, error) {
	var ma []models.Manufacturer
	res := s.db.Order("name").Limit(limit).Offset(offset).Find(&ma)
	if res.Error != nil {
		return nil, res.Error
	}

	return ma, nil
}

func (s *Storage) CreateManufacturer(data models.Manufacturer) error {
	err := s.db.Create(&data)
	return err.Error
}

func (s *Storage) CreateDiagramSubCat(data models.DiagramSubCat) error {
	err := s.db.Create(&data)
	return err.Error
}

func (s *Storage) CreateDiagramCat(data models.DiagramCat) error {
	err := s.db.Create(&data)
	return err.Error
}

func (s *Storage) GetDistributor(fv models.FindByField) (*models.Distributor, error) {
	dist := &models.Distributor{}
	db := s.db.Where(fv.Field+" = ?", fv.Value).Preload("Address").Preload("Contact").First(dist)
	if db.Error != nil {
		return nil, db.Error
	}
	return dist, nil
}

func (s *Storage) GetDistributors(limit, offset int) ([]models.Distributor, error) {
	var dits []models.Distributor
	res := s.db.Order("name").Limit(limit).Offset(offset).Preload("Address").Preload("Contact").Find(&dits)
	if res.Error != nil {
		return nil, res.Error
	}

	return dits, nil
}

func (s *Storage) CreateHead(data models.Head) error {
	err := s.db.Create(&data)
	return err.Error
}

func (s *Storage) CreateHeadType(data models.HeadType) error {
	err := s.db.Create(&data)
	return err.Error
}

func (s *Storage) GetHeadsByType(id uint, limit, offset int) ([]models.Head, error) {
	var heads []models.Head
	res := s.db.Order("name").Limit(limit).Offset(offset).Where("head_type_id = ?", id).Preload("Resources").Preload("Images").Preload("HeadType").Find(&heads)
	if res.Error != nil {
		return nil, res.Error
	}

	return heads, nil
}

func (s *Storage) FindHeadByField(fv models.FindByField) (*models.Head, error) {
	head := &models.Head{}
	db := s.db.Where(fv.Field+" = ?", fv.Value).Preload("Resources").Preload("Images").Preload("HeadType").First(head)
	if db.Error != nil {
		return nil, db.Error
	}
	return head, nil
}

func (s *Storage) GetHeads(limit, offset int) ([]models.Head, error) {
	var heads []models.Head
	res := s.db.Order("name").Limit(limit).Offset(offset).Preload("Resources").Preload("Images").Preload("HeadType").Find(&heads)
	if res.Error != nil {
		return nil, res.Error
	}

	return heads, nil
}

func (s *Storage) GetHeadTypes(limit, offset int) ([]models.HeadType, error) {
	var hTypes []models.HeadType
	err := s.db.Order("name").Limit(limit).Offset(offset).Find(&hTypes)
	if err.Error != nil {
		return nil, err.Error
	}

	return hTypes, nil
}

func (s *Storage) GetModels(limit, offset int) ([]models.Model, error) {
	var ms []models.Model
	res := s.db.Order("name").Preload("Manufacturer").Preload("Serials").Limit(limit).Offset(offset).Find(&ms)
	if res.Error != nil {
		return nil, res.Error
	}

	return ms, nil
}

func (s *Storage) GetCategories(limit, offset int) ([]models.DiagramCat, error) {
	var cats []models.DiagramCat
	err := s.db.Order("name").Limit(limit).Offset(offset).Preload("DiagramSubCats").Find(&cats)
	if err.Error != nil {
		return nil, err.Error
	}

	return cats, nil
}

func NewPostgresStorage() *Storage {
	return &Storage{}
}

func (s *Storage) ConnectDb() error {
	dsn := configs.PostgresDns()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	if err = migrateDb(db); err != nil {
		return err
	}

	s.db = db
	return nil
}

func (s *Storage) CreateUser(u models.User) error {
	err := s.db.Create(&u)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (s *Storage) CreateCustomerTempReq(cr models.CustomerTempRequest) error {
	err := s.db.Create(&cr)
	if err != nil {
		return err.Error
	}
	return nil
}

func (s *Storage) GetUsers(limit, offset int) ([]models.User, error) {
	var users []models.User
	result := s.db.Order("name").Limit(limit).Offset(offset).Find(&users)
	err := result.Error
	if err != nil {
		return nil, err
	}

	return users, err
}

func (s *Storage) FindAllOrders() ([]models.CustomerTempRequest, error) {
	var orders []models.CustomerTempRequest
	res := s.db.Find(&orders)
	err := res.Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *Storage) CreateDistributor(data models.Distributor) error {
	err := s.db.Create(&data)
	return err.Error
}

func (s *Storage) FindUserByField(fv models.FindByField) (*models.User, error) {
	user := &models.User{}
	db := s.db.Where(fv.Field+" = ?", fv.Value).First(user)
	if db.Error != nil {
		return nil, db.Error
	}

	return user, nil
}

func (s *Storage) UpdateByField(model interface{}, fv map[string]interface{}) error {
	return s.db.Model(model).Updates(fv).Error
}

func (s *Storage) FindDistributorByField(fv models.FindByField) (*models.Distributor, error) {
	dist := &models.Distributor{}
	db := s.db.Where(fv.Field+" = ?", fv.Value).Preload("Address").Preload("Contact").First(dist)
	if db.Error != nil {
		return nil, db.Error
	}
	return dist, nil
}

func (s *Storage) FindAllDistributorByField(fv models.FindByField) ([]models.Distributor, error) {
	var distributors []models.Distributor

	db := s.db.Order("name").Where(fv.Field+" = ?", fv.Value).Find(&distributors)
	if db.Error != nil {
		return nil, db.Error
	}
	return distributors, nil
}

func migrateDb(db *gorm.DB) error {
	if err := db.SetupJoinTable(&models.Diagram{}, "Parts", &models.DiagramPart{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(
		&models.Address{},
		&models.AxisHead{},
		&models.Carrier{},
		&models.Cart{},
		&models.CartItem{},
		&models.Contact{},
		&models.ContactStaff{},
		&models.ContactDetails{},
		&models.Controller{},
		&models.ControllerOrder{},
		&models.CustomerTempRequest{},
		&models.CuttingList{},
		&models.Diagram{},
		&models.DiagramCat{},
		&models.DiagramSubCat{},
		&models.Distributor{},
		&models.Equipment{},
		&models.EquipmentDealer{},
		&models.FigureImage{},
		&models.Head{},
		&models.HeadType{},
		&models.KeslaOrder{},
		&models.Manufacturer{},
		&models.Model{},
		&models.Oem{},
		&models.OfficeDetails{},
		&models.Order{},
		&models.OrderItem{},
		&models.OrderTax{},
		&models.Part{},
		&models.Payment{},
		&models.Preset{},
		&models.Resource{},
		&models.Serial{},
		&models.Tax{},
		&models.TaxExemption{},
		&models.UnitsOfMeasurement{},
		&models.User{},
		&models.UserAddress{},
	); err != nil {
		return err
	}
	return nil
}
