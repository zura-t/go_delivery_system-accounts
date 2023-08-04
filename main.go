package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/zura-t/go_delivery_system-accounts/cmd/api"
	"github.com/zura-t/go_delivery_system-accounts/internal/config"
	db "github.com/zura-t/go_delivery_system-accounts/internal/db/sqlc"
)

type Server struct {
	config config.Config
	router *gin.Engine
}

func main() {
	config, err := config.LoadConfig(".")

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("can't connect to db:", err)
	}
	
	store := db.New(conn)
	server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatal("can't create server:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("can't start server", err)
	}
}