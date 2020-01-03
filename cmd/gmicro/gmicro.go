package main

import (
	"context"
	"log"
	"net/http"
	"os"

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
	conn, err := amqp.Dial("")
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

	// Retrieve server port from env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Fatal(http.ListenAndServe(":"+port, r))
}
