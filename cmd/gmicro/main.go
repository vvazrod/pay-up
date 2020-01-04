package main

import (
	"context"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/streadway/amqp"
	"github.com/varrrro/pay-up/internal/consumer"
	"github.com/varrrro/pay-up/internal/gmicro"
	"github.com/varrrro/pay-up/internal/gmicro/group"
	"github.com/varrrro/pay-up/internal/gmicro/member"
)

func main() {
	// Open AMQP connection
	conn, err := amqp.Dial("amqp://guest:guest@rabbit:5672")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Open database connection
	db, _ := gorm.Open("sqlite3", ":memory:")
	defer db.Close()

	// Create tables
	db.CreateTable(&group.Group{}, &member.Member{})

	// Create data manager using database connection
	gm := gmicro.NewManager(db)

	// Build handler functions using data manager
	httpHandlers := gmicro.NewHTTPHandlers(gm)
	amqpHandler := gmicro.NewMessageHandler(gm)

	// Create AMQP consumer
	c, err := consumer.New(conn, "transactions", "balance", "gmicro")
	if err != nil {
		log.Fatal(err)
	}

	// Create context that can be cancelled
	ctx, cfunc := context.WithCancel(context.Background())
	defer cfunc()
	c.Start(ctx, amqpHandler) // start consumer

	// Build router with handlers
	r := gmicro.NewRouter(httpHandlers)

	// Start server
	log.Fatal(http.ListenAndServe(":8080", r))
}
