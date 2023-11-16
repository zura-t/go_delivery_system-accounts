package usecase

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/zura-t/go_delivery_system-accounts/config"
	db "github.com/zura-t/go_delivery_system-accounts/pkg/db/sqlc"
	"github.com/zura-t/go_delivery_system-accounts/token"
)

type Server struct {
	router     *gin.Engine
	store      db.Store
	config     config.Config
	tokenMaker token.Maker
}

func NewServer(store db.Store, config config.Config) (*Server, error) {
	tokenMaker, err := token.NewJwtMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("can't create token maker: %w", err)
	}

	server := &Server{
		store:      store,
		config:     config,
		tokenMaker: tokenMaker,
	}
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/users", server.CreateUser)
	router.POST("/users/login", server.LoginUser)
	router.GET("/users/:id", server.GetUser)
	router.PATCH("/users/:id", server.UpdateUser)
	router.PATCH("/users/phone_number/:id", server.AddPhone)
	router.DELETE("/users/:id", server.DeleteUser)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
