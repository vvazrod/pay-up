package gateway

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"

	log "github.com/sirupsen/logrus"

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
		}
	}
}

func postExpenseHandler(p *publisher.Publisher, rw http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"uri":    r.URL,
		"method": r.Method,
	})

	// Check if body is nil
	if r.Body == nil {
		logger.Warn("Request body is empty")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Read bytes from request's body
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.WithField("err", err).Warn("Can't read request body")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Parse JSON
	var e expense.Expense
	if err := json.Unmarshal(body, &e); err != nil {
		logger.WithField("err", err).Warn("Can't parse request body as expense")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if group IDs in path and body match
	if e.GroupID != mux.Vars(r)["groupid"] {
		logger.WithField("id", e.GroupID).Warn("Group IDs in body and path don't match")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if expense, group and payer UUIDs are valid
	if _, err := uuid.Parse(e.ID); err != nil {
		logger.WithField("id", e.ID).Warn("Expense ID isn't valid UUID")
		rw.WriteHeader(http.StatusBadRequest)
		return
	} else if _, err := uuid.Parse(e.GroupID); err != nil {
		logger.WithField("id", e.GroupID).Warn("Group ID isn't valid UUID")
		rw.WriteHeader(http.StatusBadRequest)
		return
	} else if _, err := uuid.Parse(e.Payer); err != nil {
		logger.WithField("id", e.Payer).Warn("Payer ID isn't valid UUID")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if recipient UUIDs are valid
	rec := strings.Split(e.Recipients, ";")
	for _, r := range rec {
		if _, err := uuid.Parse(r); err != nil {
			logger.WithField("id", r).Warn("Recipient ID isn't valid UUID")
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
		logger.WithField("err", err).Warn("Can't publish AMQP message")
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		rw.WriteHeader(http.StatusAccepted)
	}
}

func deleteExpenseHandler(p *publisher.Publisher, rw http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"uri":    r.URL,
		"method": r.Method,
	})

	gid := mux.Vars(r)["group_id"]

	// Check if group UUID is valid
	if _, err := uuid.Parse(gid); err != nil {
		logger.WithField("id", gid).Warn("Group ID isn't valid UUID")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Encode JSON body
	body, err := json.Marshal(&map[string]string{"group_id": gid})
	if err != nil {
		logger.WithField("err", err).Warn("Can't encode body")
		rw.WriteHeader(http.StatusInternalServerError)
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
		logger.WithField("err", err).Warn("Can't publish AMQP message")
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		rw.WriteHeader(http.StatusAccepted)
	}
}

func postPaymentHandler(p *publisher.Publisher, rw http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"uri":    r.URL,
		"method": r.Method,
	})

	// Check if body is nil
	if r.Body == nil {
		logger.Warn("Request body is empty")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Read bytes from request's body
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.WithField("err", err).Warn("Can't read request body")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Parse JSON
	var pay payment.Payment
	if err := json.Unmarshal(body, &pay); err != nil {
		logger.WithField("err", err).Warn("Can't parse request body as payment")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if group IDs in path and body match
	if pay.GroupID != mux.Vars(r)["groupid"] {
		logger.WithField("id", pay.GroupID).Warn("Group IDs in body and path don't match")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if UUIDs are valid
	if _, err := uuid.Parse(pay.ID); err != nil {
		logger.WithField("id", pay.ID).Warn("Payment ID isn't valid UUID")
		rw.WriteHeader(http.StatusBadRequest)
		return
	} else if _, err := uuid.Parse(pay.GroupID); err != nil {
		logger.WithField("id", pay.GroupID).Warn("Group ID isn't valid UUID")
		rw.WriteHeader(http.StatusBadRequest)
		return
	} else if _, err := uuid.Parse(pay.Payer); err != nil {
		logger.WithField("id", pay.Payer).Warn("Payer ID isn't valid UUID")
		rw.WriteHeader(http.StatusBadRequest)
		return
	} else if _, err := uuid.Parse(pay.Recipient); err != nil {
		logger.WithField("id", pay.Recipient).Warn("Recipient ID isn't valid UUID")
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
		logger.WithField("err", err).Warn("Can't publish AMQP message")
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		rw.WriteHeader(http.StatusAccepted)
	}
}

func deletePaymentHandler(p *publisher.Publisher, rw http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"uri":    r.URL,
		"method": r.Method,
	})

	gid := mux.Vars(r)["group_id"]

	// Check if group UUID is valid
	if _, err := uuid.Parse(gid); err != nil {
		logger.WithField("id", gid).Warn("Group ID isn't valid UUID")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Encode JSON body
	body, err := json.Marshal(&map[string]string{"group_id": gid})
	if err != nil {
		logger.WithField("err", err).Warn("Can't encode body")
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
		logger.WithField("err", err).Warn("Can't publish AMQP message")
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		rw.WriteHeader(http.StatusAccepted)
	}
}
