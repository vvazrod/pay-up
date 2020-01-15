package gmicro

import (
	"encoding/json"
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/varrrro/pay-up/internal/tmicro/expense"
	"github.com/varrrro/pay-up/internal/tmicro/payment"
)

// MessageHandler for AMQP messages.
func MessageHandler(m Manager) func(string, []byte) error {
	return func(op string, body []byte) error {
		log.WithField("operation", op).Info("AMQP message received")

		switch op {
		case "add-expense":
			return addExpenseHandler(body, m)
		case "delete-expense":
			return deleteExpenseHandler(body, m)
		case "add-payment":
			return addPaymentHandler(body, m)
		case "delete-payment":
			return deletePaymentHandler(body, m)
		default:
			err := errors.New("Wrong operation type")
			log.WithError(err).Warn("Can't handle message")
			return err
		}
	}
}

func addExpenseHandler(body []byte, m Manager) error {
	logger := log.WithField("operation", "add-expense")

	// Decode JSON
	var e expense.Expense
	if err := json.Unmarshal(body, &e); err != nil {
		logger.WithError(err).Error("Can't decode body")
		return err
	}

	// Add expense
	if err := m.AddExpense(&e); err != nil {
		logger.WithError(err).Error("Can't add expense")
		return err
	}

	return nil
}

func deleteExpenseHandler(body []byte, m Manager) error {
	logger := log.WithField("operation", "delete-expense")

	// Decode JSON
	var e expense.Expense
	if err := json.Unmarshal(body, &e); err != nil {
		logger.WithError(err).Error("Can't decode body")
		return err
	}

	// Remove expense
	if err := m.RemoveExpense(&e); err != nil {
		logger.WithError(err).Error("Can't remove expense")
		return err
	}

	return nil
}

func addPaymentHandler(body []byte, m Manager) error {
	logger := log.WithField("operation", "add-payment")

	// Decode JSON
	var p payment.Payment
	if err := json.Unmarshal(body, &p); err != nil {
		logger.WithError(err).Error("Can't decode body")
		return err
	}

	// Add payment
	if err := m.AddPayment(&p); err != nil {
		logger.WithError(err).Error("Can't add payment")
		return err
	}

	return nil
}

func deletePaymentHandler(body []byte, m Manager) error {
	logger := log.WithField("operation", "delete-payment")

	// Decode JSON
	var p payment.Payment
	if err := json.Unmarshal(body, &p); err != nil {
		logger.WithError(err).Error("Can't decode body")
		return err
	}

	// Remove payment
	if err := m.RemovePayment(&p); err != nil {
		logger.WithError(err).Error("Can't remove payment")
		return err
	}

	return nil
}
