package internalhttp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/server/pb"
	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	Host string
	Port string
	Log  Logger
	Srv  *runtime.ServeMux
}

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Debug(msg string)
	Error(msg string)
}

type Application interface {
	AddEvent(ctx context.Context, event storage.Event) error
	UpdateEvent(ctx context.Context, event storage.Event) error
	DelEvent(ctx context.Context, event storage.Event) error
	GetEvent(ctx context.Context, event storage.Event) (storage.Event, error)
	ListByDayEvents(ctx context.Context, t time.Time) ([]storage.Event, error)
	ListByWeekEvents(ctx context.Context, t time.Time) ([]storage.Event, error)
	ListByMonthEvents(ctx context.Context, t time.Time) ([]storage.Event, error)
}

func NewServer(ctx context.Context, log Logger, cfg *config.Config) *Server {
	s := &Server{Log: log, Host: cfg.HttpServ.Host, Port: cfg.HttpServ.Port}
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pb.RegisterEventServiceHandlerFromEndpoint(ctx, mux, net.JoinHostPort(cfg.GrpsServ.Host, cfg.GrpsServ.Port), opts)

	if err != nil {
		log.Info(err.Error())
	}
	s.Srv = mux
	return s
}

func (s *Server) Start(ctx context.Context) error {
	s.Log.Info(fmt.Sprintf("HTTP starting: %v:%v", s.Host, s.Port))
	if err := http.ListenAndServe(net.JoinHostPort(s.Host, s.Port), s.Srv); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.Log.Info(fmt.Sprintf("HTTP stopping:  %v:%v", s.Host, s.Port))
	return nil
}
