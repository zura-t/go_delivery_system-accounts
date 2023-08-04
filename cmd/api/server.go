package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zura-t/go_delivery_system-accounts/internal/config"
	db "github.com/zura-t/go_delivery_system-accounts/internal/db/sqlc"
)

type Server struct {
	store  db.Queries
	config config.Config
	router *gin.Engine
}

func NewServer(store *db.Queries, config config.Config) (*Server, error) {
	server := &Server{store: *store, config:  config}
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
      "message": "pong",
    })
  })
	router.GET("/accounts/:id", server.getUser)
	router.POST("/accounts", server.createUser)
	router.PATCH("/accounts/:id", server.updateUser)
	router.DELETE("/accounts/:id", server.deleteUser)

	server.router = router
	return server, nil

}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
