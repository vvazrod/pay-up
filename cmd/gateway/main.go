package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
	"github.com/varrrro/pay-up/internal/gateway"
	"github.com/varrrro/pay-up/internal/publisher"
)

func main() {
	// Open AMQP connection
	conn, err := amqp.Dial("amqp://guest:guest@rabbit:5672")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Create AMQP publisher
	pub, err := publisher.New(conn, "transactions", "management")
	if err != nil {
		log.Fatal(err)
	}

	// Create proxy
	url, err := url.Parse("http://gmicro:8080")
	if err != nil {
		log.Fatal(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Create router
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", gateway.StatusHandler).Methods("GET")
	r.HandleFunc("/groups", gateway.ProxyHandler(proxy)).Methods("POST")
	r.HandleFunc("/groups/{groupid}", gateway.ProxyHandler(proxy)).Methods("GET", "POST", "PUT", "DELETE")
	r.HandleFunc("/groups/{groupid}/members/{memberid}", gateway.ProxyHandler(proxy)).Methods("GET", "PUT", "DELETE")
	r.HandleFunc("/groups/{groupid}/expenses", gateway.ExpensesHandler(pub)).Methods("POST", "DELETE")
	r.HandleFunc("/groups/{groupid}/payments", gateway.PaymentsHandler(pub)).Methods("POST", "DELETE")

	// Start HTTP server
	log.Fatal(http.ListenAndServe(":8080", r))
}
