package rabbit

import (
	"context"
	"fmt"

	"github.com/streadway/amqp"
)

type Producer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	dsn     string
	queue   string
}

func NewProducer(dsn, queue string) *Producer {
	return &Producer{dsn: dsn, queue: queue}
}

func (p *Producer) Publish(ctx context.Context, body []byte) error {
	if p.conn == nil || p.conn.IsClosed() {
		if err := p.Connect(); err != nil {
			return fmt.Errorf("connect: %v", err)
		}
	}

	queue, err := p.channel.QueueDeclarePassive(p.queue, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("queue declare passive: %v", err)
	}

	if err = p.channel.Publish("", queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         body,
	}); err != nil {
		return fmt.Errorf("publish: %v", err)
	}
	return nil
}

func (p *Producer) Connect() error {
	var err error

	if p.conn, err = amqp.Dial(p.dsn); err != nil {
		return fmt.Errorf("dial: %v", err)
	}

	if p.channel, err = p.conn.Channel(); err != nil {
		return fmt.Errorf("channel: %v", err)
	}
	return nil
}

func (p *Producer) Close() error {
	if p.conn == nil {
		return nil
	}
	if p.channel == nil {
		return nil
	}
	if err := p.channel.Close(); err != nil {
		return err
	}
	return p.conn.Close()
}
