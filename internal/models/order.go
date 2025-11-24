package models

import (
	db "github.com/alaminpu1007/inventory-system/db/sqlc"
	"github.com/alaminpu1007/inventory-system/utils"
)

type CreateOrderItemRequest struct {
	ProductID int32   `json:"product_id" binding:"required,min=1"`
	Quantity  int32   `json:"quantity" binding:"required,min=1"`
	Price     float64 `json:"price" binding:"required,min=0"`
}

type CreateOrderRequest struct {
	// UserID int32 `json:"user_id" binding:"required,min=1"`
	// TotalPrice float64 `json:"total_price" binding:"required,min=0"`
	Items []CreateOrderItemRequest `json:"items" binding:"required,dive,required"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type ListOrderQuery struct {
	Size   int32 `form:"size" binding:"required,min=1,max=100"`
	PageNo int32 `form:"page_no" binding:"required,min=0"`
}

type OrderResponse struct {
	Order      utils.OrderResponseItem `json:"order"`
	OrderItems []db.OrderItem          `json:"items"`
}
