package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/alaminpu1007/inventory-system/db/sqlc"
	"github.com/alaminpu1007/inventory-system/utils"
	"github.com/gin-gonic/gin"
)

/* CREATE PRODUCTS */
type createProductsParams struct {
	Name        string  `json:"name" binding:"required,min=1,max=50"`
	Description *string `json:"description,omitempty"`
	Price       string  `json:"price" binding:"required,min=1,max=9223372036854775807"`
	Quantity    int32   `json:"quantity" binding:"required,max=5000"`
	CategoryID  int32   `json:"category_id" binding:"required"`
}

func (server *Server) createProducts(ctx *gin.Context) {
	var req createProductsParams

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// check description validity
	var dbDescription sql.NullString

	if req.Description != nil {
		dbDescription = sql.NullString{
			String: *req.Description,
			Valid:  true,
		}
	} else {
		dbDescription = sql.NullString{
			Valid: false,
		}
	}

	_, err := server.store.GetCategory(ctx, req.CategoryID)

	if err != nil {
		if err == sql.ErrNoRows {
			value := errors.New("Invalid category")
			ctx.JSON(http.StatusNotFound, errorResponse(value))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateProductParams{
		Name:        req.Name,
		Description: dbDescription,
		Price:       req.Price,
		Quantity:    req.Quantity,
		CategoryID:  req.CategoryID,
	}

	product, err := server.store.CreateProduct(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	NewResponse(ctx, http.StatusOK, "Product created successfully", utils.ToProductResponse(product))
}

/* CREATE PRODUCTS */
type updateProductParams struct {
	ID          int32   `uri:"id" binding:"required,min=1"`
	Name        string  `json:"name" binding:"required,min=1,max=50"`
	Description *string `json:"description,omitempty"`
	Price       string  `json:"price" binding:"required,min=1,max=9223372036854775807"`
	Quantity    int32   `json:"quantity" binding:"required,max=5000"`
	CategoryID  int32   `json:"category_id" binding:"required"`
}

func (server *Server) updateProductById(ctx *gin.Context) {
	var req updateProductParams

	if err := ctx.ShouldBindUri(&req.ID); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.store.GetCategory(ctx, req.CategoryID)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// check description validity
	var dbDescription sql.NullString

	if req.Description != nil {
		dbDescription = sql.NullString{
			String: *req.Description,
			Valid:  true,
		}
	} else {
		dbDescription = sql.NullString{
			Valid: false,
		}
	}

	arg := db.UpdateProductParams{
		ID:          req.ID,
		Name:        req.Name,
		Description: dbDescription,
		CategoryID:  req.CategoryID,
		Price:       req.Price,
		Quantity:    req.Quantity,
	}

	product, err := server.store.UpdateProduct(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	NewResponse(ctx, http.StatusOK, "Updated successfully", product)
}
