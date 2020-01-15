package tmicro

import (
	"encoding/json"
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"github.com/varrrro/pay-up/internal/publisher"
	"github.com/varrrro/pay-up/internal/tmicro/expense"
	"github.com/varrrro/pay-up/internal/tmicro/payment"
)

// MessageHandler using a data manager and message publisher.
func MessageHandler(m Manager, p publisher.Publisher) func(string, []byte) error {
	return func(op string, body []byte) error {
		log.WithField("operation", op).Info("AMQP message received")

		switch op {
		case "add-expense":
			return addExpenseHandler(body, m, p)
		case "delete-expense":
			return deleteExpenseHandler(body, m, p)
		case "add-payment":
			return addPaymentHandler(body, m, p)
		case "delete-payment":
			return deletePaymentHandler(body, m, p)
		default:
			err := errors.New("Wrong operation type")
			log.WithError(err).Warn("Can't handle message")
			return err
		}
	}
}

func addExpenseHandler(body []byte, m Manager, pub publisher.Publisher) error {
	logger := log.WithField("operation", "add-expense")

	// Decode JSON
	var e expense.Expense
	if err := json.Unmarshal(body, &e); err != nil {
		logger.WithError(err).Error("Can't parse message body as expense")
		return err
	}

	// Create expense
	if err := m.CreateExpense(&e); err != nil {
		logger.WithError(err).Error("Can't create expense")
		return err
	}

	// Publish AMQP message
	if err := pub.Publish("add-expense", body); err != nil {
		logger.WithError(err).Warn("Can't publish AMQP message")
		return err
	}

	return nil
}

func deleteExpenseHandler(body []byte, m Manager, pub publisher.Publisher) error {
	logger := log.WithField("operation", "delete-expense")

	// Decode JSON
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		logger.WithError(err).Error("Can't parse message body")
		return err
	}

	// Check if group ID is present
	groupidstr, ok := data["group_id"].(string)
	if !ok {
		logger.Warn("Group ID not present in message body")
		return errors.New("No group ID in message body")
	}

	// Check if group ID is valid UUID
	groupid, err := uuid.Parse(groupidstr)
	if err != nil {
		logger.WithField("id", groupidstr).Warn("Group ID isn't valid UUID")
		return err
	}

	// Remove last expense
	exp, err := m.RemoveLastExpense(groupid)
	if err != nil {
		logger.WithFields(log.Fields{
			"group_id": groupidstr,
			"err":      err,
		}).Warn("Can't delete expense")
		return err
	}

	// Encode JSON
	newBody, err := json.Marshal(&exp)
	if err != nil {
		logger.WithError(err).Error("Can't encode expense as JSON")
		return err
	}

	// Publish AMQP message
	if err := pub.Publish("delete-expense", newBody); err != nil {
		logger.WithError(err).Warn("Can't publish AMQP message")
		return err
	}

	return nil
}

func addPaymentHandler(body []byte, m Manager, pub publisher.Publisher) error {
	logger := log.WithField("operation", "add-payment")

	// Decode JSON
	var p payment.Payment
	if err := json.Unmarshal(body, &p); err != nil {
		logger.WithError(err).Error("Can't parse message body as payment")
		return err
	}

	// Create payment
	if err := m.CreatePayment(&p); err != nil {
		logger.WithError(err).Error("Can't create payment")
		return err
	}

	// Publish AMQP message
	if err := pub.Publish("add-payment", body); err != nil {
		logger.WithError(err).Warn("Can't publish AMQP message")
		return err
	}

	return nil
}

func deletePaymentHandler(body []byte, m Manager, pub publisher.Publisher) error {
	logger := log.WithField("operation", "delete-payment")

	// Decode JSON
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		logger.WithError(err).Error("Can't parse message body")
		return err
	}

	// Check if group ID is present
	groupidstr, ok := data["group_id"].(string)
	if !ok {
		logger.Error("Group ID not present in message body")
		return errors.New("No group ID in message body")
	}

	// Check if group ID is valid UUID
	groupid, err := uuid.Parse(groupidstr)
	if err != nil {
		logger.WithField("id", groupidstr).Error("Group ID isn't valid UUID")
		return err
	}

	// Remove last payment
	payment, err := m.RemoveLastPayment(groupid)
	if err != nil {
		logger.WithField("group_id", groupidstr).WithError(err).Error("Can't delete payment")
		return err
	}

	// Encode JSON
	newBody, err := json.Marshal(&payment)
	if err != nil {
		logger.WithError(err).Error("Can't encode payment as JSON")
		return err
	}

	// Publish AMQP message
	if err := pub.Publish("delete-payment", newBody); err != nil {
		logger.WithError(err).Warn("Can't publish AMQP message")
		return err
	}

	return nil
}
