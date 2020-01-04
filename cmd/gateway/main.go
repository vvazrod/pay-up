package main

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

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
	rabbit := "amqp://guest:guest@rabbit:5672"
	gmicro := "http://gmicro:8080"
	exchange := "transactions"
	key := "management"

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
			"err":      err,
		}).Fatal("Can't create publisher")
	}

	// Create proxy
	log.WithField("url", gmicro).Info("Creating reverse proxy")
	url, err := url.Parse(gmicro)
	if err != nil {
		log.WithFields(log.Fields{
			"url": gmicro,
			"err": err,
		}).Fatal("Can't create reverse proxy")
	}
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Create router
	r := mux.NewRouter().StrictSlash(true)
	r.Use(loggingMiddleware) // add logging middleware
	r.HandleFunc("/", gateway.StatusHandler).Methods("GET")
	r.HandleFunc("/groups", gateway.ProxyHandler(proxy)).Methods("POST")
	r.HandleFunc("/groups/{groupid}", gateway.ProxyHandler(proxy)).Methods("GET", "POST", "PUT", "DELETE")
	r.HandleFunc("/groups/{groupid}/members/{memberid}", gateway.ProxyHandler(proxy)).Methods("GET", "PUT", "DELETE")
	r.HandleFunc("/groups/{groupid}/expenses", gateway.ExpensesHandler(pub)).Methods("POST", "DELETE")
	r.HandleFunc("/groups/{groupid}/payments", gateway.PaymentsHandler(pub)).Methods("POST", "DELETE")

	// Start HTTP server
	log.WithField("port", 8080).Info("Starting server")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"uri":    r.URL,
			"method": r.Method,
		}).Info("Request received")

		next.ServeHTTP(rw, r)
	})
}
