package gateway_test

import (
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/varrrro/pay-up/internal/gateway"
	"github.com/varrrro/pay-up/internal/publisher"
)

var r *mux.Router

func TestMain(m *testing.M) {
	// Create mock publisher
	pub := publisher.MockPublisher(func(op string, body []byte) error {
		return nil
	})

	// Create router
	r = mux.NewRouter().StrictSlash(true)
	r.Use(gateway.LoggingMiddleware)
	r.HandleFunc("/", gateway.StatusHandler).Methods("GET")
	r.HandleFunc("/groups/{groupid}/expenses", gateway.ExpensesHandler(pub)).Methods("POST", "DELETE")
	r.HandleFunc("/groups/{groupid}/payments", gateway.PaymentsHandler(pub)).Methods("POST", "DELETE")

	// Run tests
	os.Exit(m.Run())
}
