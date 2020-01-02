package tmicro

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"github.com/varrrro/pay-up/internal/publisher"
	"github.com/varrrro/pay-up/internal/tmicro/expense"
	"github.com/varrrro/pay-up/internal/tmicro/payment"
)

// NewMessageHandler using the given handler functions.
func NewMessageHandler(m Manager, p *publisher.Publisher) func(*amqp.Delivery) {
	return func(msg *amqp.Delivery) {
		switch msg.Headers["operation"] {
		case "add-expense":
			if err := addExpenseHandler(msg, m, p); err != nil {
				msg.Reject(false)
			} else {
				msg.Ack(false)
			}
			break
		case "delete-expense":
			if err := deleteExpenseHandler(msg, m, p); err != nil {
				msg.Reject(false)
			} else {
				msg.Ack(false)
			}
			break
		case "add-payment":
			if err := addPaymentHandler(msg, m, p); err != nil {
				msg.Reject(false)
			} else {
				msg.Ack(false)
			}
			break
		case "delete-payment":
			if err := deletePaymentHandler(msg, m, p); err != nil {
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

func addExpenseHandler(msg *amqp.Delivery, m Manager, pub *publisher.Publisher) error {
	var e expense.Expense

	if err := json.Unmarshal(msg.Body, &e); err != nil {
		return err
	}

	err := m.CreateExpense(&e)
	if err != nil {
		return err
	}

	body, err := json.Marshal(&e)
	if err != nil {
		return err
	}

	updateMsg := amqp.Publishing{
		Headers: amqp.Table{
			"operation": "add-expense",
		},
		ContentType:  "application/json",
		DeliveryMode: 2,
		Priority:     1,
		Body:         body,
	}

	pub.Publish(&updateMsg)

	return nil
}

func deleteExpenseHandler(msg *amqp.Delivery, m Manager, pub *publisher.Publisher) error {
	var data map[string]interface{}

	if err := json.Unmarshal(msg.Body, &data); err != nil {
		return err
	}

	groupidstr, ok := data["group_id"].(string)
	if !ok {

	}

	groupid, err := uuid.Parse(groupidstr)
	if err != nil {
		return err
	}

	exp, err := m.RemoveLastExpense(groupid)
	if err != nil {
		return err
	}

	body, err := json.Marshal(&exp)
	if err != nil {
		return err
	}

	updateMsg := amqp.Publishing{
		Headers: amqp.Table{
			"operation": "delete-expense",
		},
		ContentType:  "application/json",
		DeliveryMode: 2,
		Priority:     1,
		Body:         body,
	}

	pub.Publish(&updateMsg)

	return nil
}

func addPaymentHandler(msg *amqp.Delivery, m Manager, pub *publisher.Publisher) error {
	var p payment.Payment

	if err := json.Unmarshal(msg.Body, &p); err != nil {
		return err
	}

	err := m.CreatePayment(&p)
	if err != nil {
		return err
	}

	body, err := json.Marshal(&p)
	if err != nil {
		return err
	}

	updateMsg := amqp.Publishing{
		Headers: amqp.Table{
			"operation": "add-payment",
		},
		ContentType:  "application/json",
		DeliveryMode: 2,
		Priority:     1,
		Body:         body,
	}

	pub.Publish(&updateMsg)

	return nil
}

func deletePaymentHandler(msg *amqp.Delivery, m Manager, pub *publisher.Publisher) error {
	var data map[string]interface{}

	if err := json.Unmarshal(msg.Body, &data); err != nil {
		return err
	}

	groupidstr, ok := data["group_id"].(string)
	if !ok {

	}

	groupid, err := uuid.Parse(groupidstr)
	if err != nil {
		return err
	}

	payment, err := m.RemoveLastPayment(groupid)
	if err != nil {
		return err
	}

	body, err := json.Marshal(&payment)
	if err != nil {
		return err
	}

	updateMsg := amqp.Publishing{
		Headers: amqp.Table{
			"operation": "delete-payment",
		},
		ContentType:  "application/json",
		DeliveryMode: 2,
		Priority:     1,
		Body:         body,
	}

	pub.Publish(&updateMsg)

	return nil
}
