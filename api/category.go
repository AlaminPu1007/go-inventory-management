package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"

	db "github.com/alaminpu1007/inventory-system/db/sqlc"
	"github.com/alaminpu1007/inventory-system/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createCategoryParams struct {
	Name string `json:"name"`
}

// INSERT NEW CATEGORY INTO DB
func (server *Server) createCategory(ctx *gin.Context) {
	var req createCategoryParams

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// check for admin only
	if strings.ToLower(authPayload.Role) != "admin" {
		ctx.JSON(http.StatusForbidden, map[string]string{
			"error": "You are not authorized to perform this action",
		})
		return
	}

	category, err := server.store.CreateCategory(ctx, req.Name)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			log.Println((pqErr.Code.Name()))
			switch pqErr.Code.Name() {
			case "unique_violation", "foreign_key_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// send custom response to the user
	NewResponse(ctx, http.StatusOK, "Created successfully", category)
}

type updateCategoryParams struct {
	Id   int32  `uri:"id" binding:"required,min=1"`
	Name string `json:"name"`
}

// UPDATE CATEGORY BY ID
func (server *Server) updateCategoryById(ctx *gin.Context) {
	var req updateCategoryParams

	// Bind URI
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// check for admin only
	if strings.ToLower(authPayload.Role) != "admin" {
		err := errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	// get category
	category, err := server.store.GetCategory(ctx, req.Id)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusForbidden, errorResponse(err))
	}

	arg := db.UpdateCategoryParams{
		ID:   int32(req.Id),
		Name: req.Name,
	}

	category, err = server.store.UpdateCategory(ctx, arg)

	if err != nil {
		if err == sql.ErrNoRows {
			value := errors.New("Category name already exists")
			ctx.JSON(http.StatusNotFound, errorResponse(value))
			return
		}
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	// send custom response to the user
	NewResponse(ctx, http.StatusOK, "Updated successfully", category)
}
