package internalhttp

import (
	"context"
	"fmt"
	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

type Server struct {
	Host string
	Port int
	Log  Logger
	Srv  *http.Server
}

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Debug(msg string)
	Error(msg string)
	InfoHttp(r *http.Request, httpStatus int, duration time.Duration)
}

type Application interface { // TODO
}

func NewServer(log Logger, app Application, cfg config.ServerConf) *Server {
	s := &Server{Log: log, Host: cfg.Host, Port: cfg.Port}

	m := http.NewServeMux()
	m.Handle("/", loggingMiddleware(s, s.Log))

	srv := &http.Server{Addr: fmt.Sprintf("%v:%v", s.Host, s.Port), Handler: m}
	s.Srv = srv
	return s
}

func (s *Server) Start(ctx context.Context) error {
	s.Log.Info(fmt.Sprintf("server starting: %v:%v", s.Host, s.Port))
	if err := s.Srv.ListenAndServe(); err != nil {
		s.Log.Info(err.Error())
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.Log.Info(fmt.Sprintf("server stopping:  %v:%v", s.Host, s.Port))
	if err := s.Srv.Shutdown(ctx); err != nil {
		s.Log.Error(err.Error())
		return err
	}
	return nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello-world")
}
