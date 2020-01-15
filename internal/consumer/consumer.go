package consumer

import (
	"context"
	"fmt"

	"github.com/streadway/amqp"
)

// Consumer of AMQP messages.
type Consumer struct {
	conn  *amqp.Connection
	queue string
	tag   string
}

// New Consumer instance.
func New(conn *amqp.Connection, exchange, queue, tag string) (*Consumer, error) {
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

	if _, err = ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	); err != nil {
		return nil, fmt.Errorf("Couldn't declare queue. Error: %s", err.Error())
	}

	if err = ch.QueueBind(
		queue,    // queue name
		queue,    // routing key
		exchange, // exchange name
		false,    // noWait
		nil,      // args
	); err != nil {
		return nil, fmt.Errorf("Couldn't bind queue to exchange. Error: %s", err.Error())
	}

	return &Consumer{
		conn:  conn,
		queue: queue,
		tag:   tag,
	}, nil
}

// Start consuming messages with the given handler.
//
// Based on https://codereview.stackexchange.com/a/199894
func (c *Consumer) Start(ctx context.Context, handle func(string, []byte) error) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		c.queue, // queue name
		c.tag,   // consumer tag
		false,   // autoAck
		false,   // exclusive
		false,   // noLocal
		false,   // noWait
		nil,     // args
	)
	if err != nil {
		return err
	}

	// Handle the messages
	go func() {
		select {
		case <-ctx.Done():
			return
		case msg := <-msgs:
			op, ok := msg.Headers["operation"]
			if !ok {
				msg.Nack(false, false)
			}

			if err := handle(op.(string), msg.Body); err != nil {
				msg.Nack(false, false)
			} else {
				msg.Ack(false)
			}
		}
	}()

	// Stop the consumer when context is cancelled
	go func() {
		<-ctx.Done()
		ch.Cancel(c.tag, false)
	}()

	return nil
}
