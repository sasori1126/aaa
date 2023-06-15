package carts

import "axis/ecommerce-backend/internal/models"

type CartRepo interface {
	CreateCart(data models.Cart) (*models.Cart, error)
	GetUserActiveCart(userId uint) (*models.Cart, error)
	AddCartItem(data models.CartItem) (*models.CartItem, error)
	GetCartItem(cartItemId uint) (*models.CartItem, error)
	GetCartItemsByQuery(qv models.QueryByField) ([]models.CartItem, error)
	GetCarts(limit, offset int, userId uint, status string) ([]models.Cart, error)
	GetCartIById(cartId uint) (*models.Cart, error)
	DelCartItem(cartItemIds []uint) error
	UpdateCartItemQuantity(item *models.CartItem) error
	UpdateCart(cart *models.Cart) error
}
