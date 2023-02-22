package rabbit

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	backoff "github.com/cenkalti/backoff/v3"
	"github.com/streadway/amqp"

	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/app"
	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/queue"
)

type Consumer struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	done         chan error
	dsn          string
	exchange     string
	exchangeType string
	queue        string
}

func NewConsumer(dsn, ex, exType, queue string) *Consumer {
	return &Consumer{
		dsn:          dsn,
		exchange:     ex,
		exchangeType: exType,
		queue:        queue,
		done:         make(chan error),
	}
}

type Worker func(context.Context, <-chan amqp.Delivery)

// Handler receive message -> send notify
func Handler(s *app.Sender) Worker {
	return func(ctx context.Context, m <-chan amqp.Delivery) {
		for {
			select {
			case <-ctx.Done():
				return
			case mes, ok := <-m:
				s.Log.Info("Delivery message")
				if !ok {
					s.Log.Info("Delivery err")
				}
				if len(mes.Body) == 0 {
					continue
				}

				var n queue.Notification
				if err := json.Unmarshal(mes.Body, &n); err != nil {
					s.Log.Info(err.Error())
					if err := mes.Nack(false, false); err != nil {
						s.Log.Info(err.Error())
					}
					continue
				}

				if err := s.SendNotify(ctx, n); err != nil {
					s.Log.Info(err.Error())
					if err := mes.Nack(false, false); err != nil {
						s.Log.Info(err.Error())
					}
					continue
				}

				if err := mes.Ack(false); err != nil {
					s.Log.Info(err.Error())
					return
				}
			}
		}
	}
}

func (c *Consumer) Handle(ctx context.Context, fn Worker, threads int) error {
	var err error
	if err = c.connect(); err != nil {
		return fmt.Errorf("connect: %w", err)
	}

	m, err := c.declareQueue()
	if err != nil {
		return fmt.Errorf("declare queue: %w", err)
	}

	for {
		for i := 0; i < threads; i++ {
			go fn(ctx, m)
		}

		if <-c.done != nil {
			if m, err = c.reConnect(ctx); err != nil {
				return fmt.Errorf("reConnect: %w", err)
			}
		}
		fmt.Println("Reconnected... possibly")
	}
}

func (c *Consumer) connect() error {
	var err error

	if c.conn, err = amqp.Dial(c.dsn); err != nil {
		return fmt.Errorf("dial: %w", err)
	}

	if c.channel, err = c.conn.Channel(); err != nil {
		return fmt.Errorf("channel: %w", err)
	}

	go func() {
		reason, ok := <-c.conn.NotifyClose(make(chan *amqp.Error))
		if !ok {
			if c.conn.IsClosed() {
				return
			}
			c.done <- reason
		}
	}()

	if err = c.channel.ExchangeDeclare(
		c.exchange,
		c.exchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("exchange declare: %w", err)
	}

	return nil
}

func (c *Consumer) declareQueue() (<-chan amqp.Delivery, error) {
	queue, err := c.channel.QueueDeclare(
		c.queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("queueDeclare: %w", err)
	}

	if err = c.channel.Qos(50, 0, false); err != nil {
		return nil, fmt.Errorf("qos: %w", err)
	}

	if err = c.channel.QueueBind(
		queue.Name,
		"",
		c.exchange,
		false,
		nil,
	); err != nil {
		return nil, fmt.Errorf("queueBind: %w", err)
	}

	m, err := c.channel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}

	return m, nil
}

func (c *Consumer) reConnect(ctx context.Context) (<-chan amqp.Delivery, error) {
	be := backoff.NewExponentialBackOff()
	be.MaxElapsedTime = time.Minute
	be.InitialInterval = 1 * time.Second
	be.Multiplier = 2
	be.MaxInterval = 15 * time.Second

	b := backoff.WithContext(be, ctx)
	for {
		d := b.NextBackOff()
		if d == backoff.Stop {
			return nil, fmt.Errorf("stop reconnecting")
		}

		for range time.After(d) {
			if err := c.connect(); err != nil {
				log.Printf("could not connect in reconnect call: %+v", err)
				continue
			}
			m, err := c.declareQueue()
			if err != nil {
				fmt.Printf("couldn't connect: %+v", err)
				continue
			}
			return m, nil
		}
	}
}
