package utils

import (
	"database/sql"
	"time"

	db "github.com/alaminpu1007/inventory-system/db/sqlc"
)

type OrderItemResponse struct {
	OrderItemID    int32     `json:"order_item_id"`
	OrderID        int32     `json:"order_id"`
	ProductID      int32     `json:"product_id"`
	Quantity       int32     `json:"quantity"`
	Price          string    `json:"price"`
	ItemCreatedAt  time.Time `json:"item_created_at"`
	ItemUpdatedAt  time.Time `json:"item_updated_at"`
	Status         string    `json:"status"`
	TotalAmount    string    `json:"total_amount"`
	OrderCreatedAt time.Time `json:"order_created_at"`
	OrderUpdatedAt time.Time `json:"order_updated_at"`
}

func NullTimeToTime(nt sql.NullTime) time.Time {
	if nt.Valid {
		return nt.Time
	}
	return time.Time{}
}

func ConvertOrderItem(row db.ListOrderItemsByUserRow) OrderItemResponse {
	return OrderItemResponse{
		OrderItemID:    row.OrderItemID,
		OrderID:        row.OrderID,
		ProductID:      row.ProductID,
		Quantity:       row.Quantity,
		Price:          row.Price,
		ItemCreatedAt:  row.ItemCreatedAt,
		ItemUpdatedAt:  row.ItemUpdatedAt,
		Status:         row.Status,
		TotalAmount:    row.TotalAmount,
		OrderCreatedAt: NullTimeToTime(row.OrderCreatedAt),
		OrderUpdatedAt: row.OrderUpdatedAt,
	}
}
