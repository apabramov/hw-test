package internalgrpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/server/pb"
	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedEventServiceServer
	App  Application
	Addr string
	Log  Logger
	Srv  *grpc.Server
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
	DelEvent(ctx context.Context, id string) error
	GetEvent(ctx context.Context, id string) (storage.Event, error)
	ListByDayEvents(ctx context.Context, t time.Time) ([]storage.Event, error)
	ListByWeekEvents(ctx context.Context, t time.Time) ([]storage.Event, error)
	ListByMonthEvents(ctx context.Context, t time.Time) ([]storage.Event, error)
}

func NewServer(log Logger, app Application, cfg config.GrpcServerConf) *Server {
	s := &Server{
		App:  app,
		Addr: net.JoinHostPort(cfg.Host, cfg.Port),
		Log:  log,
	}

	g := grpc.NewServer(
		grpc.UnaryInterceptor(
			loggingMiddleware(log),
		),
	)

	s.Srv = g
	pb.RegisterEventServiceServer(g, s)

	return s
}

func (s *Server) Start() error {
	list, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	s.Log.Info(fmt.Sprintf("GRPC starting: %s", s.Addr))
	return s.Srv.Serve(list)
}

func (s *Server) Stop() error {
	s.Log.Info(fmt.Sprintf("GRPC stopping: %s", s.Addr))
	s.Srv.GracefulStop()
	return nil
}

func (s *Server) Add(ctx context.Context, r *pb.EventRequest) (*pb.ErrorResponse, error) {
	se, err := s.getStorageEvent(r.GetEvent())
	if err != nil {
		return &pb.ErrorResponse{Error: err.Error()}, err
	}

	if err = s.App.AddEvent(ctx, se); err != nil {
		return &pb.ErrorResponse{Error: err.Error()}, err
	}

	return &pb.ErrorResponse{}, nil
}

func (s *Server) Update(ctx context.Context, r *pb.EventRequest) (*pb.ErrorResponse, error) {
	se, err := s.getStorageEvent(r.GetEvent())
	if err != nil {
		return &pb.ErrorResponse{Error: err.Error()}, err
	}

	if err = s.App.UpdateEvent(ctx, se); err != nil {
		return &pb.ErrorResponse{Error: err.Error()}, err
	}

	return &pb.ErrorResponse{}, nil
}

func (s *Server) Del(ctx context.Context, r *pb.IDRequest) (*pb.ErrorResponse, error) {
	if err := s.App.DelEvent(ctx, r.GetID()); err != nil {
		return &pb.ErrorResponse{Error: err.Error()}, err
	}

	return &pb.ErrorResponse{}, nil
}

func (s *Server) Get(ctx context.Context, r *pb.IDRequest) (*pb.Event, error) {
	event, err := s.App.GetEvent(ctx, r.GetID())
	if err != nil {
		return &pb.Event{}, err
	}

	return getEvent(event), nil
}

func (s *Server) ListByDay(ctx context.Context, r *pb.ListRequest) (*pb.ListResponse, error) {
	t := r.GetDate().AsTime()

	events, err := s.App.ListByDayEvents(ctx, t)
	if err != nil {
		return nil, err
	}
	return getListResponse(events), nil
}

func (s *Server) ListByWeek(ctx context.Context, r *pb.ListRequest) (*pb.ListResponse, error) {
	t := r.GetDate().AsTime()

	events, err := s.App.ListByWeekEvents(ctx, t)
	if err != nil {
		return nil, err
	}
	return getListResponse(events), nil
}

func (s *Server) ListByMonth(ctx context.Context, r *pb.ListRequest) (*pb.ListResponse, error) {
	t := r.GetDate().AsTime()

	events, err := s.App.ListByMonthEvents(ctx, t)
	if err != nil {
		return nil, err
	}
	return getListResponse(events), nil
}

func (*Server) getStorageEvent(event *pb.Event) (storage.Event, error) {
	id, err := uuid.Parse(event.GetID())
	if err != nil {
		return storage.Event{}, err
	}

	userID, err := uuid.Parse(event.GetUserId())
	if err != nil {
		return storage.Event{}, err
	}

	return storage.Event{
		ID:          id,
		Title:       event.GetTitle(),
		Date:        event.GetDate().AsTime(),
		Duration:    event.GetDuration().AsDuration(),
		Description: event.GetDescription(),
		UserId:      userID,
		Notify:      event.GetNotify().AsDuration(),
	}, nil
}

func getEvent(event storage.Event) *pb.Event {
	return &pb.Event{
		ID:          event.ID.String(),
		Title:       event.Title,
		Date:        timestamppb.New(event.Date),
		Duration:    durationpb.New(event.Duration),
		Description: event.Description,
		UserId:      event.UserId.String(),
		Notify:      durationpb.New(event.Notify),
	}
}

func getListResponse(events []storage.Event) *pb.ListResponse {
	pbEvents := make([]*pb.Event, 0)

	for _, event := range events {
		pbEvents = append(pbEvents, getEvent(event))
	}

	return &pb.ListResponse{
		Events: pbEvents,
	}
}
