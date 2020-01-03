package publisher

import (
	"fmt"

	"github.com/streadway/amqp"
)

// Publisher of AMQP messages.
type Publisher struct {
	conn     *amqp.Connection
	exchange string
	key      string
}

// New Publisher instance.
func New(conn *amqp.Connection, exchange, key string) (*Publisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Couldn't create channel. Error: %s", err.Error())
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
		return nil, fmt.Errorf("Couldn't declare exchange. Error: %s", err.Error())
	}

	return &Publisher{
		conn:     conn,
		exchange: exchange,
		key:      key,
	}, nil
}

// Publish a message to the publisher's exchange with the given routing key.
func (p *Publisher) Publish(msg *amqp.Publishing) error {
	ch, err := p.conn.Channel()
	if err != nil {
		return err
	}

	if err = ch.Publish(
		p.exchange, // exchange
		p.key,      // key
		false,      // mandatory
		false,      // immediate
		*msg,       // message
	); err != nil {
		return err
	}

	return nil
}
