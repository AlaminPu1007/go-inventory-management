package api

import (
	"database/sql"
	"math"
	"net/http"

	db "github.com/alaminpu1007/inventory-system/db/sqlc"
	"github.com/alaminpu1007/inventory-system/internal/models"
	"github.com/alaminpu1007/inventory-system/token"
	"github.com/alaminpu1007/inventory-system/utils"
	"github.com/gin-gonic/gin"
)

/* GET LOGGED USERS ALL ORDER ITEMS LISTS */
func (server *Server) getOrdersItemForLoggedUsers(ctx *gin.Context) {
	var req models.PaginationQuery

	if err := ctx.ShouldBindQuery(&req); err != nil {
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

	arg := db.ListOrderItemsByUserParams{
		UserID: user.ID,
		Limit:  req.Size,
		Offset: (req.PageNo - 1) * req.Size,
	}

	rows, err := server.store.ListOrderItemsByUser(ctx, arg)

	if err != nil {
		NewResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Map SQLC rows to response DTO
	var orderItems []utils.OrderItemResponse
	for _, row := range rows {
		orderItems = append(orderItems, utils.ConvertOrderItem(row))
	}

	// get total orders
	totalCount, err := server.store.CountActiveOrderItemsByUser(ctx, user.ID)

	if err != nil {
		NewResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	totalPages := int32(math.Ceil(float64(totalCount) / float64(req.Size)))

	if orderItems == nil {
		orderItems = []utils.OrderItemResponse{}
	}

	var message string = "Data is not found"

	if len(orderItems) != 0 {
		message = "Data is found"
	}

	data := map[string]interface{}{
		"orderItems": orderItems,
		"limit":      req.Size,
		"page":       req.PageNo,
		"totalCount": totalCount,
		"totalPages": totalPages,
	}

	NewResponse(ctx, http.StatusOK, message, data)
}

/* DELETE AN EXISTING ORDERS IF IT'S ALREADY IN INITIAL STATUS (ACTIVE) */
func (server *Server) removedItemById(ctx *gin.Context) {
	var req models.OrderIdRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		NewResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// check it's an valid item or not
	order_item1, err := server.store.GetOrderItemById(ctx, req.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			NewResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
	}

	if order_item1.Status != "active" {
		NewResponse(ctx, http.StatusNotFound, "No item is found", nil)
		return
	}

	// check this order is exists: GetOrderById
	order, err := server.store.GetOrderById(ctx, order_item1.OrderID)

	if err != nil {
		if err == sql.ErrNoRows {
			NewResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
	}

	arg := db.UpdateOrderItemStatusParams{
		ID:     order_item1.ID,
		Status: "de-active",
	}

	err = server.store.UpdateOrderItemStatus(ctx, arg)

	if err != nil {
		NewResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	// after making status de-active, need to re-calculate total prices
	_, err = server.store.RecalculateOrderTotal(ctx, order.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			NewResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}
	}

	NewResponse(ctx, http.StatusOK, "Deleted successfully", nil)
}
