package internalgrpc

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/app"
	cfg "github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/server/pb"
	memorystorage "github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/storage/memory"
)

func TestGRPC(t *testing.T) {
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	clientOptions := []grpc.DialOption{grpc.WithInsecure()}
	_, err = grpc.Dial(l.Addr().String(), clientOptions...)
	require.NoError(t, err)

	logg, err := logger.New("info", "calendar")
	require.NoError(t, err)

	storage := memorystorage.New()
	calendar := app.New(logg, storage)

	cfg := cfg.Config{
		GrpsServ: cfg.GrpcServerConf{
			Host: "",
			Port: "9000",
		},
	}

	server := NewServer(logg, calendar, cfg.GrpsServ)
	require.NoError(t, err)

	go func() {
		server.Srv.Serve(l)
	}()
}

func TestGRPCServerAdd(t *testing.T) {
	logg, err := logger.New("info", "calendar")
	require.NoError(t, err)

	storage := memorystorage.New()
	calendar := app.New(logg, storage)

	cfg := cfg.Config{
		GrpsServ: cfg.GrpcServerConf{
			Host: "",
			Port: "8082",
		},
	}

	go runGrpc(calendar, &cfg, logg)

	t.Run("add", func(t *testing.T) {
		conn, err := grpc.Dial(net.JoinHostPort(cfg.GrpsServ.Host, cfg.GrpsServ.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		c := pb.NewEventServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := c.Add(ctx, &pb.EventRequest{Event: &pb.Event{ID: "2bb0d64e-8f6e-4863-b1d8-8b20018c743d", UserId: "2bb0d64e-8f6e-4863-b1d8-8b20018c743f"}})
		require.NoError(t, err)
		require.Equal(t, "", r.GetError())
	})
}

func TestGRPCServerUpdate(t *testing.T) {
	logg, err := logger.New("info", "calendar")
	require.NoError(t, err)

	storage := memorystorage.New()
	calendar := app.New(logg, storage)

	cfg := cfg.Config{
		GrpsServ: cfg.GrpcServerConf{
			Host: "",
			Port: "8083",
		},
	}

	go runGrpc(calendar, &cfg, logg)

	t.Run("update", func(t *testing.T) {
		conn, err := grpc.Dial(net.JoinHostPort(cfg.GrpsServ.Host, cfg.GrpsServ.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		c := pb.NewEventServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		r, err := c.Add(ctx, &pb.EventRequest{Event: &pb.Event{
			ID:     "2bb0d64e-8f6e-4863-b1d8-8b20018c743d",
			Title:  "test",
			UserId: "2bb0d64e-8f6e-4863-b1d8-8b20018c743f",
		}})
		require.NoError(t, err)
		require.Equal(t, "", r.GetError())

		r, err = c.Update(ctx, &pb.EventRequest{Event: &pb.Event{
			ID:     "2bb0d64e-8f6e-4863-b1d8-8b20018c743d",
			Title:  "test update",
			UserId: "2bb0d64e-8f6e-4863-b1d8-8b20018c743f",
		}})
		require.NoError(t, err)
		require.Equal(t, "", r.GetError())

		e, err := c.Get(ctx, &pb.IDRequest{ID: "2bb0d64e-8f6e-4863-b1d8-8b20018c743d"})
		require.NoError(t, err)
		require.Equal(t, "test update", e.Event.GetTitle())
	})
}

func TestGRPCServerDel(t *testing.T) {
	logg, err := logger.New("info", "calendar")
	require.NoError(t, err)

	storage := memorystorage.New()
	calendar := app.New(logg, storage)

	cfg := cfg.Config{
		GrpsServ: cfg.GrpcServerConf{
			Host: "",
			Port: "50055",
		},
	}

	go runGrpc(calendar, &cfg, logg)

	t.Run("delete", func(t *testing.T) {
		conn, err := grpc.Dial(net.JoinHostPort(cfg.GrpsServ.Host, cfg.GrpsServ.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		c := pb.NewEventServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		r, err := c.Add(ctx, &pb.EventRequest{Event: &pb.Event{
			ID:     "2bb0d64e-8f6e-4863-b1d8-8b20018c743d",
			Title:  "test",
			UserId: "2bb0d64e-8f6e-4863-b1d8-8b20018c743f",
		}})
		require.NoError(t, err)
		require.Equal(t, "", r.GetError())

		r, err = c.Del(ctx, &pb.IDRequest{ID: "2bb0d64e-8f6e-4863-b1d8-8b20018c743d"})
		require.NoError(t, err)
		require.Equal(t, "", r.GetError())

		_, err = c.Get(ctx, &pb.IDRequest{ID: "2bb0d64e-8f6e-4863-b1d8-8b20018c743d"})
		require.Error(t, err)
	})
}

func TestGRPCServerListByDay(t *testing.T) {
	logg, err := logger.New("info", "calendar")
	require.NoError(t, err)

	storage := memorystorage.New()
	calendar := app.New(logg, storage)

	cfg := cfg.Config{
		GrpsServ: cfg.GrpcServerConf{
			Host: "",
			Port: "50056",
		},
	}

	go runGrpc(calendar, &cfg, logg)

	t.Run("list by day", func(t *testing.T) {
		conn, err := grpc.Dial(net.JoinHostPort(cfg.GrpsServ.Host, cfg.GrpsServ.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		c := pb.NewEventServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		dt := time.Now()
		ev := &pb.Event{
			ID:     "2bb0d64e-8f6e-4863-b1d8-8b20018c743d",
			Title:  "test",
			Date:   timestamppb.New(dt.Add(time.Second * 10)),
			UserId: "2bb0d64e-8f6e-4863-b1d8-8b20018c743f",
		}

		r, err := c.Add(ctx, &pb.EventRequest{Event: ev})
		require.NoError(t, err)
		require.Equal(t, "", r.GetError())

		p := pb.ListRequest{Bg: timestamppb.New(dt), Fn: timestamppb.New(dt.AddDate(0, 0, 1))}

		l, err := c.ListByDay(ctx, &p)
		require.NoError(t, err)
		m := l.GetEvents()
		require.True(t, len(m) == 1)

		require.Equal(t, ev.Title, m[0].Title)
	})
}

func TestGRPCServerListByWeek(t *testing.T) {
	logg, err := logger.New("info", "calendar")
	require.NoError(t, err)

	storage := memorystorage.New()
	calendar := app.New(logg, storage)

	cfg := cfg.Config{
		GrpsServ: cfg.GrpcServerConf{
			Host: "",
			Port: "50057",
		},
	}

	go runGrpc(calendar, &cfg, logg)

	t.Run("list by week", func(t *testing.T) {
		conn, err := grpc.Dial(net.JoinHostPort(cfg.GrpsServ.Host, cfg.GrpsServ.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		c := pb.NewEventServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		dt := time.Now()
		ev := &pb.Event{
			ID:     "2bb0d64e-8f6e-4863-b1d8-8b20018c743d",
			Title:  "test",
			Date:   timestamppb.New(dt.Add(time.Second * 10)),
			UserId: "2bb0d64e-8f6e-4863-b1d8-8b20018c743f",
		}

		r, err := c.Add(ctx, &pb.EventRequest{Event: ev})
		require.NoError(t, err)
		require.Equal(t, "", r.GetError())

		p := pb.ListRequest{Bg: timestamppb.New(dt), Fn: timestamppb.New(dt.AddDate(0, 0, 1))}

		l, err := c.ListByWeek(ctx, &p)
		require.NoError(t, err)
		m := l.GetEvents()
		require.Equal(t, ev.Title, m[0].Title)
	})
}

func TestGRPCServerListByMonth(t *testing.T) {
	logg, err := logger.New("info", "calendar")
	require.NoError(t, err)

	storage := memorystorage.New()
	calendar := app.New(logg, storage)

	cfg := cfg.Config{
		GrpsServ: cfg.GrpcServerConf{
			Host: "",
			Port: "50058",
		},
	}

	go runGrpc(calendar, &cfg, logg)

	t.Run("list by month", func(t *testing.T) {
		conn, err := grpc.Dial(net.JoinHostPort(cfg.GrpsServ.Host, cfg.GrpsServ.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		c := pb.NewEventServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		dt := time.Now()
		ev := &pb.Event{
			ID:     "2bb0d64e-8f6e-4863-b1d8-8b20018c743d",
			Title:  "test",
			Date:   timestamppb.New(dt.Add(time.Second * 10)),
			UserId: "2bb0d64e-8f6e-4863-b1d8-8b20018c743f",
		}

		r, err := c.Add(ctx, &pb.EventRequest{Event: ev})
		require.NoError(t, err)
		require.Equal(t, "", r.GetError())

		p := pb.ListRequest{Bg: timestamppb.New(dt), Fn: timestamppb.New(dt.AddDate(0, 0, 1))}

		l, err := c.ListByMonth(ctx, &p)
		require.NoError(t, err)
		m := l.GetEvents()

		require.True(t, len(m) == 1)

		require.Equal(t, ev.Title, m[0].Title)
	})
}

func runGrpc(app *app.App, cfg *cfg.Config, log *logger.Logger) {
	srv := NewServer(log, app, cfg.GrpsServ)

	lis, err := net.Listen("tcp", net.JoinHostPort(cfg.GrpsServ.Host, cfg.GrpsServ.Port))
	if err != nil {
		log.Info(err.Error())
	}
	if err := srv.Srv.Serve(lis); err != nil {
		panic(err)
	}
}
