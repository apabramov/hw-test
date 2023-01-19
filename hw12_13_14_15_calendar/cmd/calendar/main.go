package main

import (
	"context"
	"flag"
	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/util"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/app"
	cfg "github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/server/http"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := cfg.NewConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}
	logg, err := logger.New(config.Logger.Level)
	if err != nil {
		log.Fatal(err)
	}

	storage := util.NewStorage(logg, config.Storage)
	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar, config.Servers)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
