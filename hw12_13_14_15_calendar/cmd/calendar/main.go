package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/app"
	cfg "github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/server/http"
	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/util"
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

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	logg.Info("calendar is running...")
	startHTTP(ctx, &config, logg)
	startGRPC(ctx, &config, logg, calendar)
}

func startGRPC(ctx context.Context, cfg *cfg.Config, logg *logger.Logger, a *app.App) {
	srv := internalgrpc.NewServer(logg, a, cfg.GrpsServ)

	go func() {
		<-ctx.Done()
		srv.Stop()
	}()
	if err := srv.Start(); err != nil {
		logg.Error("failed to start grpc server: " + err.Error())
	}
}

func startHTTP(ctx context.Context, config *cfg.Config, logg *logger.Logger) {
	srv := internalhttp.NewServer(ctx, logg, config)

	go func() {
		if err := srv.Start(ctx); err != nil {
			logg.Error("failed to start http server: " + err.Error())
		}
	}()

	go func() {
		<-ctx.Done()
		if err := srv.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()
}
