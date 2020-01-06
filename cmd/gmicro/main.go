package main

import (
	"context"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/streadway/amqp"
	"github.com/varrrro/pay-up/internal/consumer"
	"github.com/varrrro/pay-up/internal/gmicro"
	"github.com/varrrro/pay-up/internal/gmicro/group"
	"github.com/varrrro/pay-up/internal/gmicro/member"
)

func init() {
	// Set log formatter
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	// Write logs to stdout
	log.SetOutput(os.Stdout)
}

func main() {
	rabbit := "amqp://guest:guest@rabbit:5672"
	dbconn := "host=db-gmicro port=5432 user=gmicro dbname=gmicro password=gmicro sslmode=disable"
	exchange := "transactions"
	queue := "balance"
	ctag := "gmicro"

	// Open AMQP connection
	log.WithField("url", rabbit).Info("Connecting to AMQP server")
	conn, err := amqp.Dial("amqp://guest:guest@rabbit:5672")
	if err != nil {
		log.WithField("url", rabbit).WithError(err).Fatal("AMQP server connection failure")
	}
	defer conn.Close()

	// Open database connection
	log.WithFields(log.Fields{
		"db":  "sqlite3",
		"url": dbconn,
	}).Info("Connecting to database")
	db, err := gorm.Open("postgres", dbconn)
	if err != nil {
		log.WithFields(log.Fields{
			"url": dbconn,
			"err": err,
		}).Fatal("Database connection failure")
	}
	defer db.Close()

	// Create tables
	db.CreateTable(&group.Group{}, &member.Member{})

	// Create data manager using database connection
	gm := gmicro.NewManager(db)

	// Create AMQP consumer
	log.WithFields(log.Fields{
		"exchange": exchange,
		"queue":    queue,
		"tag":      ctag,
	}).Info("Creating AMQP consumer")
	c, err := consumer.New(conn, exchange, queue, ctag)
	if err != nil {
		log.WithFields(log.Fields{
			"exchange": exchange,
			"queue":    queue,
			"tag":      ctag,
			"err":      err,
		}).Fatal("Can't create consumer")
	}

	// Create context that can be cancelled
	ctx, cfunc := context.WithCancel(context.Background())
	defer cfunc()
	c.Start(ctx, gmicro.MessageHandler(gm)) // start consumer

	// Build router with handlers
	r := mux.NewRouter().StrictSlash(true)
	r.Use(gmicro.LoggingMiddleware, gmicro.ContentTypeMiddleware)
	r.HandleFunc("/", gmicro.StatusHandler).Methods("GET")
	r.HandleFunc("/groups", gmicro.GroupsHandler(gm)).Methods("POST")
	r.HandleFunc("/groups/{groupid}", gmicro.GroupHandler(gm)).Methods("GET", "PUT", "DELETE")
	r.HandleFunc("/groups/{groupid}/members", gmicro.MembersHandler(gm)).Methods("POST")
	r.HandleFunc("/groups/{groupid}/members/{memberid}", gmicro.MemberHandler(gm)).Methods("GET", "PUT", "DELETE")

	// Start HTTP server
	log.WithField("port", 8080).Info("Starting HTTP server")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.WithError(err).Fatal("Server fail")
	}
}
