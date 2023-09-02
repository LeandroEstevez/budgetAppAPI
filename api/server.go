package api

import (
	"fmt"

	db "github.com/LeandroEstevez/budgetAppAPI/db/sqlc"
	"github.com/LeandroEstevez/budgetAppAPI/token"
	"github.com/LeandroEstevez/budgetAppAPI/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Creates a new HTTP server and setup routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setUpRouter()
	return server, nil
}

func (server *Server) setUpRouter() {
	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AllowCredentials = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PATCH", "DELETE"}
	// router.Use(cors.New(corsConfig))
	router.Use(CORSMiddleware())

	router.POST("/forgotpassword", server.forgotPassword)
	router.POST("/user", server.createUser)
	router.POST("/user/login", server.logInUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("/entry", server.addEntry)
	authRoutes.PATCH("/updateEntry", server.updateEntry)
	authRoutes.DELETE("/deleteEntry/:id", server.deleteEntry)
	authRoutes.GET("/entries", server.getEntries)
	authRoutes.GET("/categories", server.getCategories)
	authRoutes.DELETE("/deleteUser/:username", server.deleteUser)
	authRoutes.PATCH("/resetPassword", server.resetPassword)

	authRoutes.GET("/user/:username", server.getUser)

	server.router = router
}

// Runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
