package api

import (
	"net/http"
	"strings"

	"github.com/alaminpu1007/inventory-system/token"
	"github.com/gin-gonic/gin"
)

type createCategoryParams struct {
	Name string `json:"name"`
}

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

}
