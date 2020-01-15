package gateway

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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
func ExpensesHandler(p publisher.Publisher) func(http.ResponseWriter, *http.Request) {
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

func postExpenseHandler(p publisher.Publisher, rw http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"uri":    r.URL,
		"method": r.Method,
	})

	// Decode JSON
	var e expense.Expense
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		logger.WithError(err).Error("Can't parse request body as expense")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if group IDs in path and body match
	if e.GroupID.String() != mux.Vars(r)["groupid"] {
		logger.WithField("id", e.GroupID).Error("Group IDs in body and path don't match")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if recipient UUIDs are valid
	rec := strings.Split(e.Recipients, ";")
	for _, r := range rec {
		if _, err := uuid.Parse(r); err != nil {
			logger.WithField("id", r).Error("Recipient ID isn't valid UUID")
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	// Encode JSON
	body, err := json.Marshal(&e)
	if err != nil {
		logger.WithError(err).Error("Can't encode body")
	}

	// Publish AMQP message
	if err := p.Publish("add-expense", body); err != nil {
		logger.WithError(err).Warn("Can't publish AMQP message")
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		rw.WriteHeader(http.StatusAccepted)
	}
}

func deleteExpenseHandler(p publisher.Publisher, rw http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"uri":    r.URL,
		"method": r.Method,
	})

	gid := mux.Vars(r)["group_id"]

	// Check if group UUID is valid
	if _, err := uuid.Parse(gid); err != nil {
		logger.WithField("id", gid).Error("Group ID isn't valid UUID")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Encode JSON
	body, err := json.Marshal(&map[string]string{"group_id": gid})
	if err != nil {
		logger.WithError(err).Error("Can't encode body")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Publish AMQP message
	if err := p.Publish("delete-expense", body); err != nil {
		logger.WithError(err).Warn("Can't publish AMQP message")
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		rw.WriteHeader(http.StatusAccepted)
	}
}

// PaymentsHandler that publishes to an AMQP queue.
func PaymentsHandler(p publisher.Publisher) func(http.ResponseWriter, *http.Request) {
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

func postPaymentHandler(p publisher.Publisher, rw http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"uri":    r.URL,
		"method": r.Method,
	})

	// Decode JSON
	var pay payment.Payment
	if err := json.NewDecoder(r.Body).Decode(&pay); err != nil {
		logger.WithError(err).Error("Can't parse request body as payment")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if group IDs in path and body match
	if pay.GroupID.String() != mux.Vars(r)["groupid"] {
		logger.WithField("id", pay.GroupID).Error("Group IDs in body and path don't match")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Encode JSON
	body, err := json.Marshal(&pay)
	if err != nil {
		logger.WithError(err).Error("Can't encode body")
	}

	// Publish AMQP message
	if err := p.Publish("add-payment", body); err != nil {
		logger.WithError(err).Warn("Can't publish AMQP message")
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		rw.WriteHeader(http.StatusAccepted)
	}
}

func deletePaymentHandler(p publisher.Publisher, rw http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"uri":    r.URL,
		"method": r.Method,
	})

	gid := mux.Vars(r)["group_id"]

	// Check if group UUID is valid
	if _, err := uuid.Parse(gid); err != nil {
		logger.WithField("id", gid).Error("Group ID isn't valid UUID")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Encode JSON
	body, err := json.Marshal(&map[string]string{"group_id": gid})
	if err != nil {
		logger.WithError(err).Error("Can't encode body")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Publish AMQP message
	if err := p.Publish("delete-payment", body); err != nil {
		logger.WithError(err).Warn("Can't publish AMQP message")
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		rw.WriteHeader(http.StatusAccepted)
	}
}
