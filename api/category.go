package api

import (
	"database/sql"
	"errors"
	"log"
	"math"
	"net/http"

	db "github.com/alaminpu1007/inventory-system/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

/* CREATE NEW CATEGORY */
type createCategoryParams struct {
	Name string `json:"name"`
}

func (server *Server) createCategory(ctx *gin.Context) {
	var req createCategoryParams

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
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

/* UPDATE CATEGORY BY NAME */
type updateCategoryParams struct {
	Id   int32  `uri:"id" binding:"required,min=1"`
	Name string `json:"name"`
}

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

	// get category
	category, err := server.store.GetCategory(ctx, req.Id)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusForbidden, errorResponse(err))
	}

	arg := db.PatchCategoryParams{
		ID:   int32(req.Id),
		Name: req.Name,
	}

	category, err = server.store.PatchCategory(ctx, arg)

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

/* SEARCH CATEGORY BY ID */

type searchCategoryQuery struct {
	Name   string `form:"name" binding:"required"`
	Size   int32  `form:"size" binding:"required,min=1,max=100"` // page_size
	PageNo int32  `form:"page_no" binding:"required,min=0"`      // page number, 0-based
}

func (server *Server) searchCategoryByName(ctx *gin.Context) {
	var req searchCategoryQuery

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.SearchCategoryParams{
		Column1: req.Name,
		Limit:   req.Size,
		Offset:  (req.PageNo - 1) * req.Size,
	}

	categories, err := server.store.SearchCategory(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	totalCount, err := server.store.CountSearchCategories(ctx, req.Name)
	totalPages := int32(math.Ceil(float64(totalCount) / float64(req.Size)))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	data := map[string]interface{}{
		"categories": categories,
		"limit":      req.Size,
		"page":       req.PageNo,
		"totalCount": totalCount,
		"totalPages": int32(totalPages),
	}

	// Send custom response
	NewResponse(ctx, http.StatusOK, "Data found", data)
}

/* GET CATEGORY BY ID */
type getCategoryParams struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getCategoryById(ctx *gin.Context) {
	var req getCategoryParams

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	category, err := server.store.GetCategory(ctx, req.ID)

	if err != nil {

		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	NewResponse(ctx, http.StatusOK, "Data is found", category)
}

/* LIST CATEGORY */
type listCategoryParams struct {
	Size   int32 `form:"size" binding:"required,min=1,max=100"` // page_size
	PageNo int32 `form:"page_no" binding:"required,min=0"`      // page number, 0-based
}

func (server *Server) listCategory(ctx *gin.Context) {
	var req listCategoryParams

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListCategoryParams{
		Limit:  req.Size,
		Offset: (req.PageNo - 1) * req.Size,
	}

	categories, err := server.store.ListCategory(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	totalCount, err := server.store.CountListCategory(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	totalPages := int32(math.Ceil(float64(totalCount) / float64(req.Size)))

	data := map[string]interface{}{
		"categories": categories,
		"limit":      req.Size,
		"page":       req.PageNo,
		"totalCount": totalCount,
		"totalPages": totalPages,
	}

	NewResponse(ctx, http.StatusOK, "Data found", data)
}

/* REMOVED CATEGORY BY ID */
func (server *Server) deleteCategoryById(ctx *gin.Context) {
	var req getCategoryParams

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// find item to check it's validity
	_, err := server.store.GetCategory(ctx, req.ID)

	if err != nil {

		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	category, err := server.store.RemoveCategory(ctx, req.ID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	NewResponse(ctx, http.StatusOK, "Data deleted successfully", category)
}
