package internalhttp

import (
	"bytes"
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/app"
	cfg "github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/server/grpc"
	memorystorage "github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

func runRest(ctx context.Context, cfg *cfg.Config) {
	srv := NewServer(ctx, nil, cfg)
	log.Printf("server listening at %s:%s", cfg.HttpServ.Host, cfg.HttpServ.Port)
	if err := http.ListenAndServe(net.JoinHostPort(cfg.HttpServ.Host, cfg.HttpServ.Port), srv.Srv); err != nil {
		panic(err)
	}
}

func runGrpc(app *app.App, cfg *cfg.Config, log *logger.Logger) {
	srv := internalgrpc.NewServer(log, app, cfg.GrpsServ)

	lis, err := net.Listen("tcp", net.JoinHostPort(cfg.GrpsServ.Host, cfg.GrpsServ.Port))
	if err != nil {
		log.Info(err.Error())
	}
	if err := srv.Srv.Serve(lis); err != nil {
		panic(err)
	}
}

func TestHTTPServerAdd(t *testing.T) {
	logg, err := logger.New("info")
	require.NoError(t, err)

	storage := memorystorage.New()
	calendar := app.New(logg, storage)

	cfg := cfg.Config{
		HttpServ: cfg.HttpServerConf{
			Host: "127.0.0.1",
			Port: "8000",
		},
		GrpsServ: cfg.GrpcServerConf{
			Host: "127.0.0.1",
			Port: "9000",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	go runRest(ctx, &cfg)
	go runGrpc(calendar, &cfg, logg)

	event := bytes.NewBufferString(`{
  "event" : {
		"ID":"2bb0d64e-8f6e-4863-b1d8-8b20018c743d",
		"Title":"Hello",
		"Date":"2023-01-01T16:00:00Z",
		"Duration": "600s",
		"Description":"Hello",
		"UserId":"cc526645-6fad-461e-9ebf-82a7d936a61f",
		"Notify":"300s"
  }
}`)

	t.Run("add", func(t *testing.T) {
		resp, err := http.Post("http://127.0.0.1:8000/v1/event/add", "application/json", event)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, "{\"Error\":\"\"}", string(body))
	})
}

func TestHTTPServerUpdate(t *testing.T) {
	logg, err := logger.New("info")
	require.NoError(t, err)

	storage := memorystorage.New()
	calendar := app.New(logg, storage)

	cfg := cfg.Config{
		HttpServ: cfg.HttpServerConf{
			Host: "127.0.0.1",
			Port: "18001",
		},
		GrpsServ: cfg.GrpcServerConf{
			Host: "127.0.0.1",
			Port: "9001",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	go runRest(ctx, &cfg)
	go runGrpc(calendar, &cfg, logg)

	event := bytes.NewBufferString(`{
  "event" : {
		"ID":"2bb0d64e-8f6e-4863-b1d8-8b20018c743d",
		"Title":"Hello",
		"Date":"2023-01-01T16:00:00Z",
		"Duration": "600s",
		"Description":"Hello",
		"UserId":"cc526645-6fad-461e-9ebf-82a7d936a61f",
		"Notify":"300s"
  }
}`)

	eu := bytes.NewBufferString(`{
  "event" : {
		"ID":"2bb0d64e-8f6e-4863-b1d8-8b20018c743d",
		"Title":"Hello update",
		"Date":"2023-01-01T16:00:00Z",
		"Duration": "600s",
		"Description":"Hello",
		"UserId":"cc526645-6fad-461e-9ebf-82a7d936a61f",
		"Notify":"300s"
  }
}`)

	t.Run("update", func(t *testing.T) {
		resp, err := http.Post("http://127.0.0.1:18001/v1/event/add", "", event)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, "{\"Error\":\"\"}", string(body))

		cli := http.Client{}
		req, err := http.NewRequest(http.MethodPut, "http://127.0.0.1:18001/v1/event/update", eu)
		require.NoError(t, err)
		resp, err = cli.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		body, err = io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, "{\"Error\":\"\"}", string(body))

		resp, err = http.Get("http://127.0.0.1:18001/v1/event/get/2bb0d64e-8f6e-4863-b1d8-8b20018c743d")
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		body, err = io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, "{\"ID\":\"2bb0d64e-8f6e-4863-b1d8-8b20018c743d\",\"Title\":\"Hello update\",\"Date\":\"2023-01-01T16:00:00Z\",\"Duration\":\"600s\",\"Description\":\"Hello\",\"UserId\":\"cc526645-6fad-461e-9ebf-82a7d936a61f\",\"Notify\":\"300s\"}", string(body))
	})
}

func TestHTTPServerDelete(t *testing.T) {
	logg, err := logger.New("info")
	require.NoError(t, err)

	storage := memorystorage.New()
	calendar := app.New(logg, storage)

	cfg := cfg.Config{
		HttpServ: cfg.HttpServerConf{
			Host: "127.0.0.1",
			Port: "18002",
		},
		GrpsServ: cfg.GrpcServerConf{
			Host: "127.0.0.1",
			Port: "9002",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	go runRest(ctx, &cfg)
	go runGrpc(calendar, &cfg, logg)

	event := bytes.NewBufferString(`{
  "event" : {
		"ID":"2bb0d64e-8f6e-4863-b1d8-8b20018c743d",
		"Title":"Hello",
		"Date":"2023-01-01T16:00:00Z",
		"Duration": "600s",
		"Description":"Hello",
		"UserId":"cc526645-6fad-461e-9ebf-82a7d936a61f",
		"Notify":"300s"
  }
}`)

	t.Run("delete", func(t *testing.T) {
		resp, err := http.Post("http://127.0.0.1:18002/v1/event/add", "", event)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, "{\"Error\":\"\"}", string(body))

		cli := http.Client{}
		req, err := http.NewRequest(http.MethodDelete, "http://127.0.0.1:18002/v1/event/delete/2bb0d64e-8f6e-4863-b1d8-8b20018c743d", nil)
		require.NoError(t, err)
		resp, err = cli.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		body, err = io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, "{\"Error\":\"\"}", string(body))

		resp, err = http.Get("http://127.0.0.1:18002/v1/event/get/2bb0d64e-8f6e-4863-b1d8-8b20018c743d")
		defer resp.Body.Close()
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		body, err = io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, "{\"code\":2,\"message\":\"event not exists\",\"details\":[]}", string(body))
	})
}

func TestHTTPServerListByDay(t *testing.T) {
	logg, err := logger.New("info")
	require.NoError(t, err)

	storage := memorystorage.New()
	calendar := app.New(logg, storage)

	cfg := cfg.Config{
		HttpServ: cfg.HttpServerConf{
			Host: "127.0.0.1",
			Port: "18003",
		},
		GrpsServ: cfg.GrpcServerConf{
			Host: "127.0.0.1",
			Port: "9003",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	go runRest(ctx, &cfg)
	go runGrpc(calendar, &cfg, logg)

	event := bytes.NewBufferString(`{
  "event" : {
		"ID":"2bb0d64e-8f6e-4863-b1d8-8b20018c743d",
		"Title":"Hello",
		"Date":"2023-01-01T16:00:00Z",
		"Duration": "600s",
		"Description":"Hello",
		"UserId":"cc526645-6fad-461e-9ebf-82a7d936a61f",
		"Notify":"300s"
  }
}`)

	t.Run("list day", func(t *testing.T) {
		resp, err := http.Post("http://127.0.0.1:18003/v1/event/add", "", event)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, "{\"Error\":\"\"}", string(body))

		resp, err = http.Post("http://127.0.0.1:18003/v1/event/list/day", "", bytes.NewBufferString("{\n  \"Date\":\"2023-01-01T16:00:00Z\"\n}"))
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		body, err = io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, "{\"Events\":[]}", string(body))
	})
}

func TestHTTPServerListByWeek(t *testing.T) {
	logg, err := logger.New("info")
	require.NoError(t, err)

	storage := memorystorage.New()
	calendar := app.New(logg, storage)

	cfg := cfg.Config{
		HttpServ: cfg.HttpServerConf{
			Host: "127.0.0.1",
			Port: "18004",
		},
		GrpsServ: cfg.GrpcServerConf{
			Host: "127.0.0.1",
			Port: "9004",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	go runRest(ctx, &cfg)
	go runGrpc(calendar, &cfg, logg)

	event := bytes.NewBufferString(`{
  "event" : {
		"ID":"2bb0d64e-8f6e-4863-b1d8-8b20018c743d",
		"Title":"Hello",
		"Date":"2023-01-01T16:00:00Z",
		"Duration": "600s",
		"Description":"Hello",
		"UserId":"cc526645-6fad-461e-9ebf-82a7d936a61f",
		"Notify":"300s"
  }
}`)

	t.Run("list week", func(t *testing.T) {
		resp, err := http.Post("http://127.0.0.1:18004/v1/event/add", "", event)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, "{\"Error\":\"\"}", string(body))

		resp, err = http.Post("http://127.0.0.1:18004/v1/event/list/week", "", bytes.NewBufferString("{\n  \"Date\":\"2023-01-01T16:00:00Z\"\n}"))
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		body, err = io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, "{\"Events\":[]}", string(body))
	})
}

func TestHTTPServerListByMonth(t *testing.T) {
	logg, err := logger.New("info")
	require.NoError(t, err)

	storage := memorystorage.New()
	calendar := app.New(logg, storage)

	cfg := cfg.Config{
		HttpServ: cfg.HttpServerConf{
			Host: "127.0.0.1",
			Port: "18005",
		},
		GrpsServ: cfg.GrpcServerConf{
			Host: "127.0.0.1",
			Port: "9005",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	go runRest(ctx, &cfg)
	go runGrpc(calendar, &cfg, logg)

	event := bytes.NewBufferString(`{
  "event" : {
		"ID":"2bb0d64e-8f6e-4863-b1d8-8b20018c743d",
		"Title":"Hello",
		"Date":"2023-01-01T16:00:00Z",
		"Duration": "600s",
		"Description":"Hello",
		"UserId":"cc526645-6fad-461e-9ebf-82a7d936a61f",
		"Notify":"300s"
  }
}`)

	t.Run("list month", func(t *testing.T) {
		resp, err := http.Post("http://127.0.0.1:18005/v1/event/add", "", event)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, "{\"Error\":\"\"}", string(body))

		resp, err = http.Post("http://127.0.0.1:18005/v1/event/list/month", "", bytes.NewBufferString("{\n  \"Date\":\"2023-01-01T00:00:00Z\"\n}"))
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		body, err = io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, "{\"Events\":[{\"ID\":\"2bb0d64e-8f6e-4863-b1d8-8b20018c743d\",\"Title\":\"Hello\",\"Date\":\"2023-01-01T16:00:00Z\",\"Duration\":\"600s\",\"Description\":\"Hello\",\"UserId\":\"cc526645-6fad-461e-9ebf-82a7d936a61f\",\"Notify\":\"300s\"}]}", string(body))
	})
}
