package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/app"
	cfg "github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/queue"
	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/queue/rabbit"
	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/util"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config_scheduler.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config, err := cfg.NewSchedulerCfg(configFile)
	if err != nil {
		log.Fatal(err)
	}
	logg, err := logger.New(config.Logger.Level, "calender_scheduler")
	if err != nil {
		log.Fatal(err)
	}

	storage := util.NewStorage(logg, config.Storage)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	scheduler := app.NewScheduler(logg, storage)

	go startScheduler(ctx, cfg.SchedulerConfig{}, scheduler)
	go startDelete(ctx, cfg.SchedulerConfig{}, scheduler)

	logg.Info("scheduler is running...")

	<-ctx.Done()
}

func startScheduler(ctx context.Context, c cfg.SchedulerConfig, s *app.Scheduler) {
	p := rabbit.NewProducer(c.Queue.Dsn, c.Queue.Queue)

	d, err := time.ParseDuration(c.Ticker.Duration)
	if err != nil {
		s.Log.Error(err.Error())
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.Tick(d):
			events, err := s.ListNotify(ctx)
			if err != nil {
				s.Log.Info(err.Error())
				continue
			}

			for _, e := range events {
				m := queue.Notification{
					ID:     e.ID,
					Title:  e.Title,
					Date:   e.Date,
					UserId: e.UserId,
				}

				body, err := json.Marshal(m)
				if err != nil {
					s.Log.Info(err.Error())
					continue
				}

				if err = p.Publish(ctx, body); err != nil {
					s.Log.Info(err.Error())
				}

				s.Log.Info("published to queue: " + m.GetId())
			}
		}
	}
}

func startDelete(ctx context.Context, c cfg.SchedulerConfig, s *app.Scheduler) {
	d, err := time.ParseDuration(c.Ticker.Duration)
	if err != nil {
		s.Log.Info(err.Error())
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		select {
		case <-ctx.Done():
			return
		case <-time.Tick(d):
			if err := s.DeleteOutDate(ctx); err != nil {
				s.Log.Info(err.Error())
				continue
			}
			s.Log.Info("deleted event")
		}
	}
}
