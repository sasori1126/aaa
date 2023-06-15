package models

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/dto"
)

type Cart struct {
	configs.GormModel

	Items       []CartItem
	Status      string
	SubTotal    float64
	TotalAmount float64
	User        User
	UserId      uint
}

func (c Cart) ToResponse() dto.CartResponse {
	cid, _ := EncodeHashId(c.ID)
	uid, _ := EncodeHashId(c.UserId)

	var total float64
	var items []dto.CartItemResponse
	items = []dto.CartItemResponse{}
	for _, item := range c.Items {
		total += item.Amount
		itemResponse := item.ToResponse()
		items = append(items, itemResponse)
	}

	return dto.CartResponse{
		Id: cid,
		User: dto.EmbedUser{
			Id:          uid,
			Name:        c.User.Name,
			Email:       c.User.Email,
			PhoneNumber: c.User.PhoneNumber,
			IsActive:    c.User.IsActive,
			Verified:    c.User.IsActive,
		},
		Items:       items,
		SubTotal:    total,
		TotalAmount: total,
	}
}

type CartItem struct {
	configs.GormModel
	CartId       uint
	Cart         Cart
	PartId       uint
	Part         Part
	PricePerUnit float64
	Quantity     float64
	Unit         string
	Name         string
	Description  string
	Weight       float64
	Note         string
	Title        string
	Amount       float64
}

func (ci CartItem) ToResponse() dto.CartItemResponse {
	itemId, _ := EncodeHashId(ci.ID)
	cartId, _ := EncodeHashId(ci.CartId)
	partId, _ := EncodeHashId(ci.PartId)
	partResp := ci.Part.ToResponse()
	return dto.CartItemResponse{
		Id:           itemId,
		CartId:       cartId,
		PartId:       partId,
		PricePerUnit: ci.PricePerUnit,
		Quantity:     ci.Quantity,
		Unit:         ci.Unit,
		Description:  ci.Description,
		Note:         ci.Note,
		Title:        ci.Title,
		Amount:       ci.Amount,
		Part:         partResp,
	}
}
