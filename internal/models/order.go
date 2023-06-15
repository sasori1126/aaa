package models

import (
	"time"

	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/dto"
)

type Order struct {
	configs.GormModel

	BillingAddress           UserAddress
	BillingAddressId         uint
	CancelledAt              *time.Time
	CancelReason             string
	Cart                     Cart
	CartId                   uint
	ConvergePayTransactionId string
	Currency                 string
	DeliveryInstruction      string
	DeliveryMethod           string
	DeliveryRateId           string
	Items                    []OrderItem
	PaidAmount               float64
	PaymentMethod            string
	PurchaseOrderNumber      string
	PickupPoint              string
	ReturnedAt               *time.Time
	ShippedAt                *time.Time
	ShippingAddressId        uint
	ShippingAddress          UserAddress
	ShippingAmount           float64
	Status                   string
	SubTotalAmount           float64
	Taxes                    []OrderTax
	TotalAmount              float64
	TotalTaxAmount           float64
	TrackingNumber           string
	UserId                   uint
	User                     User
}

type OrderItem struct {
	configs.GormModel

	Description          string
	Name                 string
	Note                 string
	Order                Order
	OrderId              uint
	Part                 Part
	PartId               uint
	PricePerUnit         float64
	Quantity             float64
	ShippingMethodDetail string
	Total                float64
	Unit                 string
	Weight               float64
}

func (i OrderItem) ToResponse() dto.OrderItemResponse {
	itemId, _ := EncodeHashId(i.ID)
	oid, _ := EncodeHashId(i.OrderId)
	pid, _ := EncodeHashId(i.PartId)
	part := i.Part
	return dto.OrderItemResponse{
		Description: i.Description,
		Id:          itemId,
		Name:        i.Name,
		Note:        i.Note,
		OrderId:     oid,
		Part: dto.EmbedPart{
			Code:                  part.Code,
			CountryOfOrigin:       part.CountryOfOrigin,
			CylinderType:          part.CylinderType,
			Description:           part.Description,
			Detail:                part.Detail,
			DealerPrice:           part.DealerPrice,
			DealerPricePercentage: part.DealerPricePercentage,
			DealerPricePerMeter:   part.DealerPricePerMeter,
			Drivers:               part.Drivers,
			Featured:              part.Featured,
			GuideLink:             part.GuideLink,
			Height:                part.Height,
			Id:                    pid,
			Length:                part.Length,
			Material:              part.Material,
			MetaDescription:       part.MetaDescription,
			MetaKeywords:          part.MetaKeywords,
			Name:                  part.Name,
			OemCompatible:         part.OemCompatible,
			OemNumber:             part.OemNumber,
			PartDiagramNumber:     part.PartDiagramNumber,
			PartNumber:            part.PartNumber,
			Price:                 part.Price,
			PricePerMeter:         part.PricePerMeter,
			QuantityOnHand:        part.QuantityOnHand,
			QuantityOnOrder:       part.QuantityOnOrder,
			QuantityOnSaleOrder:   part.QuantityOnSaleOrder,
			QuantityRecommended:   part.QuantityRecommended,
			SalePrice:             part.SalePrice,
			SalePricePerMeter:     part.SalePricePerMeter,
			Seo:                   part.Seo,
			Status:                part.Status,
			VideoUrl:              part.VideoUrl,
			Weight:                part.Weight,
			WhereUsed:             part.WhereUsed,
			Width:                 part.Width,
		},
		PricePerUnit: i.PricePerUnit,
		Quantity:     i.Quantity,
		Total:        i.Total,
		Unit:         i.Unit,
		Weight:       i.Weight,
	}
}

type OrderTax struct {
	configs.GormModel

	Amount            float64
	Description       string
	Name              string
	OrderId           uint
	ProductTaxAmount  float64
	Rate              float64 // Delete column RateType
	ShippingTaxAmount float64
}

func (t OrderTax) ToResponse(totalTaxAmount float64) dto.OrderTaxResponse {
	taxId, _ := EncodeHashId(t.ID)
	oid, _ := EncodeHashId(t.OrderId)

	productTaxAmount := t.ProductTaxAmount
	shippingTaxAmount := t.ShippingTaxAmount
	// For backward compatiblity, there was no product_tax_amount/shipping_tax_amount column in DB. And the system was not
	// charge shipping tax. Hence, totalTaxAmount meant productTaxAmount.
	if productTaxAmount == 0 {
		productTaxAmount = t.Amount
		shippingTaxAmount = 0.0
	}

	return dto.OrderTaxResponse{
		Amount:            t.Amount,
		Description:       t.Description,
		Id:                taxId,
		Name:              t.Name,
		OrderId:           oid,
		ProductTaxAmount:  productTaxAmount,
		Rate:              t.Rate,
		ShippingTaxAmount: shippingTaxAmount,
	}
}

func (o Order) ToResponse() dto.OrderResponse {
	oid, _ := EncodeHashId(o.ID)
	uid, _ := EncodeHashId(o.UserId)
	cid, _ := EncodeHashId(o.CartId)

	sAddress := o.ShippingAddress.ToResponse()
	bAddress := o.BillingAddress.ToResponse()

	var items []dto.OrderItemResponse
	items = []dto.OrderItemResponse{}
	for _, item := range o.Items {
		t := item.ToResponse()
		items = append(items, t)
	}

	var taxes []dto.OrderTaxResponse
	taxes = []dto.OrderTaxResponse{}
	for _, tax := range o.Taxes {
		t := tax.ToResponse(o.TotalTaxAmount)
		taxes = append(taxes, t)
	}

	return dto.OrderResponse{
		BillingAddress:      bAddress,
		CancelledAt:         o.CancelledAt,
		CancelReason:        o.CancelReason,
		CartId:              cid,
		Currency:            o.Currency,
		DeliveryInstruction: o.DeliveryInstruction,
		DeliveryMethod:      o.DeliveryMethod,
		DeliveryRateId:      o.DeliveryRateId,
		Id:                  oid,
		Items:               items,
		OrderDate:           o.CreatedAt,
		PaidAmount:          o.PaidAmount,
		PaymentMethod:       o.PaymentMethod,
		PickupPoint:         o.PickupPoint,
		PurchaseOrderNumber: o.PurchaseOrderNumber,
		ReturnedAt:          o.ReturnedAt,
		ShippedAt:           o.ShippedAt,
		ShippingAddress:     sAddress,
		ShippingAmount:      o.ShippingAmount,
		Status:              o.Status,
		SubTotalAmount:      o.SubTotalAmount,
		Taxes:               taxes,
		TotalAmount:         o.TotalAmount,
		TotalTaxAmount:      o.TotalTaxAmount,
		TrackingNumber:      o.TrackingNumber,
		User: dto.EmbedUser{
			Email:       o.User.Email,
			Id:          uid,
			IsActive:    o.User.IsActive,
			Name:        o.User.Name,
			PhoneNumber: o.User.PhoneNumber,
			Verified:    o.User.IsActive,
		},
		UserId: uid,
	}
}
