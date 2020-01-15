package publisher

import (
	"github.com/streadway/amqp"
)

// Publisher of AMQP messages.
type Publisher interface {
	Publish(op string, body []byte) error
}

// MockPublisher used in tests.
type MockPublisher func(op string, body []byte) error

// Publish function that calls the mock function.
func (p MockPublisher) Publish(op string, body []byte) error {
	return p(op, body)
}

// AMQPPublisher implementation.
type AMQPPublisher struct {
	conn     *amqp.Connection
	exchange string
	key      string
}

// New AMQPPublisher instance.
func New(conn *amqp.Connection, exchange, key string) (*AMQPPublisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	if err = ch.ExchangeDeclare(
		exchange, // name
		"direct", // type
		true,     // durable
		false,    // autoDelete
		false,    // internal
		false,    // noWait
		nil,      // args
	); err != nil {
		return nil, err
	}

	return &AMQPPublisher{
		conn:     conn,
		exchange: exchange,
		key:      key,
	}, nil
}

// Publish a message to the publisher's exchange with the given routing key.
func (p *AMQPPublisher) Publish(op string, body []byte) error {
	ch, err := p.conn.Channel()
	if err != nil {
		return err
	}

	msg := amqp.Publishing{
		Headers: amqp.Table{
			"operation": op,
		},
		ContentType:  "application/json",
		DeliveryMode: 2,
		Priority:     1,
		Body:         body,
	}

	if err = ch.Publish(
		p.exchange, // exchange
		p.key,      // key
		false,      // mandatory
		false,      // immediate
		msg,        // message
	); err != nil {
		return err
	}

	return nil
}
