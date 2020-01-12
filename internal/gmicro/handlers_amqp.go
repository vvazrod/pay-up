package gmicro

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/streadway/amqp"
	"github.com/varrrro/pay-up/internal/tmicro/expense"
	"github.com/varrrro/pay-up/internal/tmicro/payment"
)

// MessageHandler for AMQP messages.
func MessageHandler(m Manager) func(*amqp.Delivery) {
	return func(msg *amqp.Delivery) {
		log.WithField("operation", msg.Headers["operation"]).Info("AMQP message received")

		switch msg.Headers["operation"] {
		case "add-expense":
			if err := addExpenseHandler(msg, m); err != nil {
				msg.Reject(false)
			} else {
				msg.Ack(false)
			}
			break
		case "delete-expense":
			if err := deleteExpenseHandler(msg, m); err != nil {
				msg.Reject(false)
			} else {
				msg.Ack(false)
			}
			break
		case "add-payment":
			if err := addPaymentHandler(msg, m); err != nil {
				msg.Reject(false)
			} else {
				msg.Ack(false)
			}
			break
		case "delete-payment":
			if err := deletePaymentHandler(msg, m); err != nil {
				msg.Reject(false)
			} else {
				msg.Ack(false)
			}
			break
		default:
			msg.Reject(false)
		}
	}
}

func addExpenseHandler(msg *amqp.Delivery, m Manager) error {
	logger := log.WithField("operation", "add-expense")

	// Decode JSON
	var e expense.Expense
	if err := json.Unmarshal(msg.Body, &e); err != nil {
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

func deleteExpenseHandler(msg *amqp.Delivery, m Manager) error {
	logger := log.WithField("operation", "add-expense")

	// Decode JSON
	var e expense.Expense
	if err := json.Unmarshal(msg.Body, &e); err != nil {
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

func addPaymentHandler(msg *amqp.Delivery, m Manager) error {
	logger := log.WithField("operation", "add-expense")

	// Decode JSON
	var p payment.Payment
	if err := json.Unmarshal(msg.Body, &p); err != nil {
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

func deletePaymentHandler(msg *amqp.Delivery, m Manager) error {
	logger := log.WithField("operation", "add-expense")

	// Decode JSON
	var p payment.Payment
	if err := json.Unmarshal(msg.Body, &p); err != nil {
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
