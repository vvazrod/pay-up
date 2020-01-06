package main

import (
	"context"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/streadway/amqp"
	"github.com/varrrro/pay-up/internal/consumer"
	"github.com/varrrro/pay-up/internal/publisher"
	"github.com/varrrro/pay-up/internal/tmicro"
	"github.com/varrrro/pay-up/internal/tmicro/expense"
	"github.com/varrrro/pay-up/internal/tmicro/payment"
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
	dbconn := "host=db-tmicro port=5432 user=tmicro dbname=tmicro password=tmicro sslmode=disable"
	exchange := "transactions"
	key := "balance"
	queue := "management"
	ctag := "tmicro"

	// Open AMQP connection
	log.WithField("url", rabbit).Info("Connecting to AMQP server")
	conn, err := amqp.Dial(rabbit)
	if err != nil {
		log.WithFields(log.Fields{
			"url": rabbit,
			"err": err,
		}).Fatal("AMQP server connection failure")
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
	db.CreateTable(&expense.Expense{}, &payment.Payment{})

	// Create data manager
	tm := tmicro.NewManager(db)

	// Create AMQP publisher
	log.WithFields(log.Fields{
		"exchange": exchange,
		"key":      key,
	}).Info("Creating AMQP publisher")
	pub, err := publisher.New(conn, exchange, key)
	if err != nil {
		log.WithFields(log.Fields{
			"exchange": exchange,
			"key":      key,
			"err":      err,
		}).Fatal("Can't create publisher")
	}

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

	// Create channel to listen for OS signals
	sch := make(chan os.Signal, 1)
	signal.Notify(sch, os.Interrupt, os.Kill)

	// Create context that can be cancelled
	ctx, cfunc := context.WithCancel(context.Background())
	defer cfunc()

	log.Info("Starting AMQP consumer")
	c.Start(ctx, tmicro.MessageHandler(tm, pub)) // start consumer

	<-sch // blocking until we receive a signal
}
