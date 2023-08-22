package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/zura-t/go_delivery_system-accounts/cmd/gapi"
	"github.com/zura-t/go_delivery_system-accounts/pb"
	"github.com/zura-t/go_delivery_system-accounts/internal"
	db "github.com/zura-t/go_delivery_system-accounts/internal/db/sqlc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := internal.LoadConfig(".")
	if err != nil {
		log.Fatal("can't load config", err)
	}
	
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("can't connect to db:", err)
	}
	
	store := db.New(conn)
	server, err := gapi.NewServer(store, config)
	if err != nil {
		log.Fatal("can't create server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUsersServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", config.ServerAddress)
	if err != nil {
		log.Fatalf("can't create listener: %s", err)
	}

	log.Printf("start GRPC server at %s", lis.Addr().String())
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("can't start GRPC server: %s", err)
	}
}