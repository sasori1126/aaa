package models

import "axis/ecommerce-backend/internal/dto"

type Stat struct {
	TotalUsers        int64               `json:"total_users"`
	TotalOrders       int64               `json:"total_orders"`
	ActiveCarts       int64               `json:"active_carts"`
	TotalParts        int64               `json:"total_parts"`
	UndeliveredOrders []dto.OrderResponse `json:"undelivered_orders"`
	PendingOrders     []dto.OrderResponse `json:"pending_orders"`
}
