package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zura-t/go_delivery_system-accounts/cmd/api"
	"github.com/zura-t/go_delivery_system-accounts/internal"
	db "github.com/zura-t/go_delivery_system-accounts/internal/db/sqlc"
	"github.com/zura-t/go_delivery_system-accounts/rmq"
)

func main() {
	config, err := internal.LoadConfig(".")
	if err != nil {
		log.Fatal("can't load config", err)
	}

	dbconn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("can't connect to db:", err)
	}

	store := db.NewStore(dbconn)
	if err != nil {
		log.Fatal("can't create store:", err)
	}

	rabbitConn, err := connectRabbitmq()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	_, err = runGinServer(store, config)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// channel, consumer, err := setupRabbitmq(rabbitConn, server)
	// if err != nil {
	// 	panic(err)
	// }
	// defer channel.Close()

	// err = consumer.Listen([]string{})
	// if err != nil {
	// 	log.Println(err)
	// }
}

func runGinServer(store db.Store, config internal.Config) (*api.Server, error) {
	server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatalf("can't create server: %s", err)
		return nil, err
	}

	err = server.Start(config.HttpServerAddress)
	if err != nil {
		log.Fatalf("can't start server: %s", err)
		return nil, err
	}
	return server, nil
}

func connectRabbitmq() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection
	for {
		c, err := amqp.Dial("amqp://guest:guest@localhost:5672")
		if err != nil {
			fmt.Println("rabbitmq not yet ready")
			counts++
		} else {
			log.Println("Connected to rabbitmq.")

			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("back off")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}

func setupRabbitmq(rabbitConn *amqp.Connection, server *api.Server) (*amqp.Channel, *rmq.Consumer, error) {
	channel, err := rabbitConn.Channel()
	if err != nil {
		log.Fatal("can't create rabbitmq consumer", err)
		return nil, nil, err
	}

	consumer, err := rmq.NewConsumer(rabbitConn, channel, server)
	if err != nil {
		return nil, nil, err
	}

	return channel, &consumer, nil
}
