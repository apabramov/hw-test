package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"github.com/streadway/amqp"

	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/app"
	cfg "github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/queue"
	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/queue/rabbit"
	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/util"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config_sender.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config, err := cfg.NewSenderCfg(configFile)
	if err != nil {
		log.Fatal(err)
	}
	logg, err := logger.New(config.Logger.Level, "calender_sender")
	if err != nil {
		log.Fatal(err)
	}

	storage := util.NewStorage(logg, config.Storage)
	sender := app.NewSender(logg, storage)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go startSender(ctx, sender, config.Queue)

	logg.Info("sender is running...")

	<-ctx.Done()
}

func startSender(ctx context.Context, s *app.Sender, c cfg.QueueConf) {
	cn := rabbit.NewConsumer(c.Dsn, c.Exchange, c.ExchangeType, c.Queue)

	if err := cn.Handle(ctx, Worker(s), 1); err != nil {
		s.Log.Info(err.Error())
	}
}

// Worker receive message -> send notify
func Worker(s *app.Sender) rabbit.Worker {
	return func(ctx context.Context, m <-chan amqp.Delivery) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			select {
			case <-ctx.Done():
				return
			case mes := <-m:
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
