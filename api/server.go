package api

import (
	"fmt"

	db "github.com/LeandroEstevez/budgetAppAPI/db/sqlc"
	"github.com/LeandroEstevez/budgetAppAPI/token"
	"github.com/LeandroEstevez/budgetAppAPI/util"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests
type Server struct {
	config util.Config
	store db.Store
	tokenMaker token.Maker
	router *gin.Engine
}

// Creates a new HTTP server and setup routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server {
		config: config,
		store: store,
		tokenMaker: tokenMaker,
	}

	server.setUpRouter()
	return server, nil
}

func (server *Server) setUpRouter() {
	router := gin.Default()
	// corsConfig := cors.DefaultConfig()
	// corsConfig.AllowOrigins = []string{"*"}
	// corsConfig.AllowCredentials = true
	// router.Use(cors.New(corsConfig))

	router.POST("/user", server.createUser)
	router.POST("/user/login", server.logInUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("/entry", server.addEntry)
	authRoutes.DELETE("/deleteEntry", server.deleteEntry)
	authRoutes.GET("/entries", server.getEntries)

	authRoutes.GET("/user/:username", server.getUser)

	server.router = router
}

// Runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H {"error": err.Error()}
}