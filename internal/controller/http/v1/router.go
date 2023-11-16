package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zura-t/go_delivery_system-accounts/internal/usecase"
	"github.com/zura-t/go_delivery_system-accounts/pkg/logger"
)

func (server *Server) NewRouter(handler *gin.Engine, logger logger.Interface, userUsecase usecase.User) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	h := handler.Group("/v1")
	handler.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	{
		server.newUserRoutes(h, userUsecase, logger)
	}
}