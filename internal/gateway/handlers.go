package gateway

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
	"github.com/varrrro/pay-up/internal/publisher"
	"github.com/varrrro/pay-up/internal/tmicro/expense"
	"github.com/varrrro/pay-up/internal/tmicro/payment"
)

// StatusHandler returns a static message to know the server is working.
func StatusHandler(rw http.ResponseWriter, r *http.Request) {
	status := map[string]string{"status": "OK"}

	rw.WriteHeader(http.StatusOK)
	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(&status)
}

// ProxyHandler for requests that need to be sent to another service.
func ProxyHandler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		p.ServeHTTP(rw, r)
	}
}

// ExpensesHandler that publishes to an AMQP queue.
func ExpensesHandler(p *publisher.Publisher) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			postExpenseHandler(p, rw, r)
			break
		case "DELETE":
			deleteExpenseHandler(p, rw, r)
			break
		default:
			rw.WriteHeader(http.StatusBadRequest)
		}
	}
}

// PaymentsHandler that publishes to an AMQP queue.
func PaymentsHandler(p *publisher.Publisher) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			postPaymentHandler(p, rw, r)
			break
		case "DELETE":
			deletePaymentHandler(p, rw, r)
			break
		default:
			rw.WriteHeader(http.StatusBadRequest)
		}
	}
}

func postExpenseHandler(p *publisher.Publisher, rw http.ResponseWriter, r *http.Request) {
	// Check if body is nil
	if r.Body == nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Read bytes from request's body
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Parse JSON
	var e expense.Expense
	if err := json.Unmarshal(body, &e); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if expense, group and payer UUIDs are valid
	if _, err := uuid.Parse(e.ID); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	} else if _, err := uuid.Parse(e.GroupID); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	} else if _, err := uuid.Parse(e.Payer); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if recipient UUIDs are valid
	rec := strings.Split(e.Recipients, ";")
	for _, r := range rec {
		if _, err := uuid.Parse(r); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	// Create AMQP message
	msg := amqp.Publishing{
		Headers: amqp.Table{
			"operation": "add-expense",
		},
		ContentType:  "application/json",
		DeliveryMode: 2,
		Priority:     1,
		Body:         body,
	}

	// Publish AMQP message
	if err := p.Publish(&msg); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		rw.WriteHeader(http.StatusAccepted)
	}
}

func deleteExpenseHandler(p *publisher.Publisher, rw http.ResponseWriter, r *http.Request) {
	gid := mux.Vars(r)["group_id"]

	// Check if group UUID is valid
	if _, err := uuid.Parse(gid); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Encode JSON body
	body, err := json.Marshal(&map[string]string{"group_id": gid})
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Create AMQP message
	msg := amqp.Publishing{
		Headers: amqp.Table{
			"operation": "delete-expense",
		},
		ContentType:  "application/json",
		DeliveryMode: 2,
		Priority:     1,
		Body:         body,
	}

	if err := p.Publish(&msg); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		rw.WriteHeader(http.StatusAccepted)
	}
}

func postPaymentHandler(p *publisher.Publisher, rw http.ResponseWriter, r *http.Request) {
	// Check if body is nil
	if r.Body == nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Read bytes from request's body
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Parse JSON
	var pay payment.Payment
	if err := json.Unmarshal(body, &pay); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if UUIDs are valid
	if _, err := uuid.Parse(pay.ID); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	} else if _, err := uuid.Parse(pay.GroupID); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	} else if _, err := uuid.Parse(pay.Payer); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	} else if _, err := uuid.Parse(pay.Recipient); err != nil {
		rw.WriteHeader((http.StatusBadRequest))
		return
	}

	// Create AMQP message
	msg := amqp.Publishing{
		Headers: amqp.Table{
			"operation": "add-payment",
		},
		ContentType:  "application/json",
		DeliveryMode: 2,
		Priority:     1,
		Body:         body,
	}

	// Publish AMQP message
	if err := p.Publish(&msg); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		rw.WriteHeader(http.StatusAccepted)
	}
}

func deletePaymentHandler(p *publisher.Publisher, rw http.ResponseWriter, r *http.Request) {
	gid := mux.Vars(r)["group_id"]

	// Check if group UUID is valid
	if _, err := uuid.Parse(gid); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Encode JSON body
	body, err := json.Marshal(&map[string]string{"group_id": gid})
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Create AMQP message
	msg := amqp.Publishing{
		Headers: amqp.Table{
			"operation": "delete-payment",
		},
		ContentType:  "application/json",
		DeliveryMode: 2,
		Priority:     1,
		Body:         body,
	}

	if err := p.Publish(&msg); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		rw.WriteHeader(http.StatusAccepted)
	}
}
