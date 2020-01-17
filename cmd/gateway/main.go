package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
	"github.com/varrrro/pay-up/internal/gateway"
	"github.com/varrrro/pay-up/internal/publisher"
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
	rabbit := os.Getenv("RABBIT_CONN")
	gmicro := os.Getenv("PROXY_URL")
	exchange := os.Getenv("EXCHANGE")
	key := os.Getenv("KEY")

	// Open AMQP connection
	log.WithField("url", rabbit).Info("Connecting to AMQP server")
	conn, err := amqp.Dial(rabbit)
	if err != nil {
		log.WithField("url", rabbit).WithError(err).Fatal("AMQP server connection failure")
	}
	defer conn.Close()

	// Create AMQP publisher
	log.WithFields(log.Fields{
		"exchange": exchange,
		"key":      key,
	}).Info("Creating AQMP publisher")
	pub, err := publisher.New(conn, exchange, key)
	if err != nil {
		log.WithFields(log.Fields{
			"exchange": exchange,
			"key":      key,
		}).WithError(err).Fatal("Can't create publisher")
	}

	// Create proxy
	log.WithField("url", gmicro).Info("Creating reverse proxy")
	url, err := url.Parse(gmicro)
	if err != nil {
		log.WithField("url", gmicro).WithError(err).Fatal("Can't create reverse proxy")
	}
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Create router
	r := mux.NewRouter().StrictSlash(true)
	r.Use(gateway.LoggingMiddleware)
	r.HandleFunc("/", gateway.StatusHandler).Methods("GET")
	r.HandleFunc("/groups", gateway.ProxyHandler(proxy)).Methods("POST")
	r.HandleFunc("/groups/{groupid}", gateway.ProxyHandler(proxy)).Methods("GET", "PUT", "DELETE")
	r.HandleFunc("/groups/{groupid}/members", gateway.ProxyHandler(proxy)).Methods("POST")
	r.HandleFunc("/groups/{groupid}/members/{memberid}", gateway.ProxyHandler(proxy)).Methods("GET", "PUT", "DELETE")
	r.HandleFunc("/groups/{groupid}/expenses", gateway.ExpensesHandler(pub)).Methods("POST", "DELETE")
	r.HandleFunc("/groups/{groupid}/payments", gateway.PaymentsHandler(pub)).Methods("POST", "DELETE")

	// Start HTTP server
	log.WithField("port", 8080).Info("Starting HTTP server")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.WithError(err).Fatal("Server fail")
	}
}
