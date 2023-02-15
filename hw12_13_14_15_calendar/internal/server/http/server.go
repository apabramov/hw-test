package internalhttp

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/server/pb"
)

type Server struct {
	Host string
	Port string
	Log  Logger
	Srv  *http.Server
}

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Debug(msg string)
	Error(msg string)
}

func NewServer(ctx context.Context, log Logger, cfg *config.Config) (*Server, error) {
	s := &Server{Log: log, Host: cfg.HttpServ.Host, Port: cfg.HttpServ.Port}
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pb.RegisterEventServiceHandlerFromEndpoint(ctx, mux, net.JoinHostPort(cfg.GrpsServ.Host, cfg.GrpsServ.Port), opts)
	if err != nil {
		return nil, err
	}
	srv := &http.Server{
		Addr:    net.JoinHostPort(cfg.HttpServ.Host, cfg.HttpServ.Port),
		Handler: mux,
	}
	s.Srv = srv
	return s, nil
}

func (s *Server) Start(ctx context.Context) error {
	s.Log.Info(fmt.Sprintf("HTTP starting: %v:%v", s.Host, s.Port))

	if err := s.Srv.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.Log.Info(fmt.Sprintf("HTTP stopping:  %v:%v", s.Host, s.Port))
	s.Srv.Shutdown(ctx)
	return nil
}
