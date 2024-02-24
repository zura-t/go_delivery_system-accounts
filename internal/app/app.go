package app

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zura-t/go_delivery_system-accounts/config"
	v1 "github.com/zura-t/go_delivery_system-accounts/internal/controller/http/v1"
	"github.com/zura-t/go_delivery_system-accounts/internal/usecase"
	db "github.com/zura-t/go_delivery_system-accounts/pkg/db/sqlc"
	"github.com/zura-t/go_delivery_system-accounts/pkg/logger"
	"github.com/zura-t/go_delivery_system-accounts/rmq"
	"github.com/zura-t/go_delivery_system-accounts/token"
)

func Run(config *config.Config) {
	dbconn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("can't connect to db:", err)
	}

	store := db.NewStore(dbconn)
	if err != nil {
		log.Fatal("can't create store:", err)
	}

	// rabbitConn, err := connectRabbitmq()
	// if err != nil {
	// 	log.Println(err)
	// 	os.Exit(1)
	// }
	// defer rabbitConn.Close()

	l := logger.New(config.LogLevel)

	tokenMaker, err := token.NewJwtMaker(config.TokenSymmetricKey)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	usersUseCase := usecase.New(store, config, tokenMaker)

	runGinServer(l, config, usersUseCase)

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

func runGinServer(l *logger.Logger, config *config.Config, usersUsecase *usecase.UserUseCase) {
	handler := gin.New()

	server, err := v1.New(config)
	if err != nil {
		log.Fatalf("can't create server: %s", err)
	}

	server.NewRouter(handler, l, usersUsecase)
	handler.Run(config.HttpPort)
}

func connectRabbitmq() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection
	for {
		c, err := amqp.Dial("amqp://admin:admin@rabbitmq:5672")
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

func setupRabbitmq(rabbitConn *amqp.Connection) (*amqp.Channel, *rmq.Consumer, error) {
	channel, err := rabbitConn.Channel()
	if err != nil {
		log.Fatal("can't create rabbitmq consumer", err)
		return nil, nil, err
	}

	consumer, err := rmq.NewConsumer(rabbitConn, channel)
	if err != nil {
		return nil, nil, err
	}

	return channel, &consumer, nil
}
