package carts

import (
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/internal/storage"
)

type CartRepoDb struct {
	client storage.Storage
}

func (c CartRepoDb) GetCarts(limit, offset int, userId uint, status string) ([]models.Cart, error) {
	return c.client.GetCarts(limit, offset, userId, status)
}

func (c CartRepoDb) GetCartIById(cartId uint) (*models.Cart, error) {
	return c.client.GetCartById(cartId)
}

func (c CartRepoDb) GetCartItemsByQuery(qv models.QueryByField) ([]models.CartItem, error) {
	return c.client.GetCartItemsByQuery(qv)
}

func (c CartRepoDb) UpdateCart(cart *models.Cart) error {
	return c.client.UpdateCart(cart)
}

func (c CartRepoDb) GetCartItem(cartItemId uint) (*models.CartItem, error) {
	return c.client.GetCartItemByField(models.FindByField{Field: "id", Value: cartItemId})
}

func (c CartRepoDb) UpdateCartItemQuantity(item *models.CartItem) error {
	return c.client.UpdateCartItemQuantity(item)
}

func (c CartRepoDb) DelCartItem(cartItemIds []uint) error {
	return c.client.DeleteCartItem(cartItemIds)
}

func (c CartRepoDb) AddCartItem(data models.CartItem) (*models.CartItem, error) {
	return c.client.AddCartItem(data)
}

func (c CartRepoDb) CreateCart(data models.Cart) (*models.Cart, error) {
	return c.client.CreateCart(data)
}

func (c CartRepoDb) GetUserActiveCart(userId uint) (*models.Cart, error) {
	return c.client.GetUserActiveCart(userId)
}

func NewCartRepoDb(db storage.Storage) CartRepo {
	return &CartRepoDb{client: db}
}
