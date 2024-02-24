package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zura-t/go_delivery_system-accounts/internal/usecase"
	"github.com/zura-t/go_delivery_system-accounts/pkg/logger"
)

func (server *Server) NewRouter(handler *gin.Engine, logger logger.Interface, userUsecase *usecase.UserUseCase) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	handler.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	{
		server.newUserRoutes(handler, userUsecase, logger)
	}
}