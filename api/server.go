package api

import (
	"fmt"

	db "github.com/alaminpu1007/inventory-system/db/sqlc"
	"github.com/alaminpu1007/inventory-system/token"
	"github.com/alaminpu1007/inventory-system/utils"
	"github.com/gin-gonic/gin"
)

// Serve http request for our banking service
type Server struct {
	config     utils.Config
	store      *db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config utils.Config, store *db.Store) (*Server, error) {

	// create a token maker
	// if you want to use JWT token maker, just replace the method with: token.NewJWTMaker()

	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)

	if err != nil {
		return nil, fmt.Errorf("can not create token maker: %w", err)
	}

	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config, // we will get token maker related info later
	}

	server.setupRouter()

	return server, nil
}

// ALL INITIALIZED ROUTER WILL BE GOES HERE
func (server *Server) setupRouter() {

	router := gin.Default()
	// create a router group for protected route or users
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// create a user
	router.POST("/users/register", server.createUser)
	// login user
	router.POST("/users/login", server.loginUser)

	// CATEGORY ROUTES GOES HERE
	// create category
	authRoutes.POST(`/category/create`, adminMiddleware(), server.createCategory)
	// update category by id
	authRoutes.PATCH("/category/:id", adminMiddleware(), server.updateCategoryById)
	// search category by name
	authRoutes.GET("/category/search", adminMiddleware(), server.searchCategoryByName)
	// get list of category by {limit, offset}
	authRoutes.GET("/categories", adminMiddleware(), server.listCategory)
	// get category by id
	authRoutes.GET("/category/:id", adminMiddleware(), server.getCategoryById)

	server.router = router
}

// START: runs the HTTP server on a specif address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// CREATE ERROR HANDLER TO SERVER ERROR JSON GLOBALLY
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
