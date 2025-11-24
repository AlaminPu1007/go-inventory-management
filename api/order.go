package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/alaminpu1007/inventory-system/db/sqlc"
	"github.com/alaminpu1007/inventory-system/internal/models"
	"github.com/alaminpu1007/inventory-system/token"
	"github.com/alaminpu1007/inventory-system/utils"
	"github.com/gin-gonic/gin"
)

func (server *Server) createOrder(ctx *gin.Context) {
	var req models.CreateOrderRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		NewResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// Get current user from context
	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	user, err := server.store.GetUser(ctx, payload.Username)

	if err != nil {
		if err != sql.ErrNoRows {
			NewResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}

		NewResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	var total float64

	for _, item := range req.Items {
		total += item.Price * float64(item.Quantity)
	}

	arg := db.CreateOrderParams{
		UserID:      user.ID,
		TotalAmount: fmt.Sprintf("%.2f", total),
	}

	// insert into order table
	order, err := server.store.CreateOrder(ctx, arg)

	if err != nil {
		NewResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	var createdItems []db.OrderItem

	// Insert order items table:
	for _, item := range req.Items {

		itemArg := db.CreateOrderItemParams{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     fmt.Sprintf("%.2f", item.Price),
		}

		createdItem, err := server.store.CreateOrderItem(ctx, itemArg)

		if err != nil {
			NewResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		// add to array
		createdItems = append(createdItems, createdItem)
	}

	createdOrder := db.Order{
		ID:          order.ID,
		UserID:      order.UserID,
		TotalAmount: order.TotalAmount,
		Status:      order.Status,
		CreatedAt:   order.CreatedAt,
	}

	// final response
	response := models.OrderResponse{
		Order:      utils.ToOrderResponse(createdOrder),
		OrderItems: createdItems,
	}

	NewResponse(ctx, http.StatusOK, "Order created successfully", response)
}
