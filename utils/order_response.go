package utils

import (
	"time"

	db "github.com/alaminpu1007/inventory-system/db/sqlc"
)

type OrderResponseItem struct {
	ID          int32     `json:"id"`
	UserID      int32     `json:"user_id"`
	TotalAmount string    `json:"total_amount"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

func ToOrderResponse(order db.Order) OrderResponseItem {
	return OrderResponseItem{
		ID:          order.ID,
		UserID:      order.UserID,
		TotalAmount: order.TotalAmount,
		Status:      order.Status,
		CreatedAt:   order.CreatedAt.Time,
	}
}
