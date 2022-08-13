package api

import (
	db "github.com/LeandroEstevez/budgetAppAPI/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests
type Server struct {
	store db.Store
	router *gin.Engine
}

// Creates a new HTTP server and setup routing
func NewServer(store db.Store) *Server {
	server := &Server {
		store: store,
	}
	router := gin.Default()

	router.POST("/user", server.createUser)
	router.POST("/entry", server.addEntry)
	router.DELETE("/deleteEntry", server.deleteEntry)
	router.GET("/user/:username", server.getUser)
	router.GET("/entries", server.getEntries)

	server.router = router

	return server
}

// Runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H {"error": err.Error()}
}