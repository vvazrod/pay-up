package gmicro

import (
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"github.com/varrrro/pay-up/internal/tmicro/expense"
	"github.com/varrrro/pay-up/internal/tmicro/payment"
)

// NewMessageHandler for AMQP messages.
func NewMessageHandler(m Manager) func(*amqp.Delivery) {
	return func(msg *amqp.Delivery) {
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
	var e expense.Expense
	if err := json.Unmarshal(msg.Body, &e); err != nil {
		return err
	}

	groupid, err := uuid.Parse(e.GroupID)
	if err != nil {
		return err
	}

	payerid, err := uuid.Parse(e.Payer)
	if err != nil {
		return err
	}

	recipients := strings.Split(e.Recipients, ";")
	var recids []uuid.UUID
	for _, r := range recipients {
		id, err := uuid.Parse(r)
		if err != nil {
			return err
		}

		recids = append(recids, id)
	}

	if err := m.AddExpense(e.Amount, groupid, payerid, &recids); err != nil {
		return err
	}

	return nil
}

func deleteExpenseHandler(msg *amqp.Delivery, m Manager) error {
	var e expense.Expense
	if err := json.Unmarshal(msg.Body, &e); err != nil {
		return err
	}

	groupid, err := uuid.Parse(e.GroupID)
	if err != nil {
		return err
	}

	payerid, err := uuid.Parse(e.Payer)
	if err != nil {
		return err
	}

	recipients := strings.Split(e.Recipients, ";")
	var recids []uuid.UUID
	for _, r := range recipients {
		id, err := uuid.Parse(r)
		if err != nil {
			return err
		}

		recids = append(recids, id)
	}

	if err := m.RemoveExpense(e.Amount, groupid, payerid, &recids); err != nil {
		return err
	}

	return nil
}

func addPaymentHandler(msg *amqp.Delivery, m Manager) error {
	var p payment.Payment
	if err := json.Unmarshal(msg.Body, &p); err != nil {
		return err
	}

	groupid, err := uuid.Parse(p.GroupID)
	if err != nil {
		return err
	}

	payerid, err := uuid.Parse(p.Payer)
	if err != nil {
		return err
	}

	recipientid, err := uuid.Parse(p.Recipient)
	if err != nil {
		return err
	}

	if err := m.AddPayment(p.Amount, groupid, payerid, recipientid); err != nil {
		return err
	}

	return nil
}

func deletePaymentHandler(msg *amqp.Delivery, m Manager) error {
	var p payment.Payment
	if err := json.Unmarshal(msg.Body, &p); err != nil {
		return err
	}

	groupid, err := uuid.Parse(p.GroupID)
	if err != nil {
		return err
	}

	payerid, err := uuid.Parse(p.Payer)
	if err != nil {
		return err
	}

	recipientid, err := uuid.Parse(p.Recipient)
	if err != nil {
		return err
	}

	if err := m.RemovePayment(p.Amount, groupid, payerid, recipientid); err != nil {
		return err
	}

	return nil
}
