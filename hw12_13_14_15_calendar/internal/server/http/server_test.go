package internalhttp

import (
	"bytes"
	"context"
	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/server/pb"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/app"
	cfg "github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/server/grpc"
	memorystorage "github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

const bufSize = 1024 * 1024

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
		panic(err)
	}
	if err := srv.Srv.Serve(lis); err != nil {
		panic(err)
	}
}

func startGRPC(t *testing.T, ctx context.Context, cfg *cfg.Config, logg *logger.Logger, a *app.App) {
	l, err := net.Listen("tcp", ":"+cfg.GrpsServ.Port)
	require.NoError(t, err)

	clientOptions := []grpc.DialOption{grpc.WithInsecure()}
	cc, err := grpc.Dial(l.Addr().String(), clientOptions...)

	srv := internalgrpc.NewServer(logg, a, cfg.GrpsServ)

	client := pb.NewEventServiceClient(cc)

	go func() {
		srv.Srv.Serve(l)
	}()

	_ = client
	//go func() {
	//	<-ctx.Done()
	//	srv.Stop()
	//}()
	//
	//go func() {
	//	if err := srv.Start(); err != nil {
	//		logg.Error("failed to start grpc server: " + err.Error())
	//	}
	//}()
}

func startHTTP(ctx context.Context, config *cfg.Config, logg *logger.Logger) *Server {
	//l, err := net.Listen("tcp", ":8070")

	srv := NewServer(ctx, logg, config)

	//go func() {
	//	if err := http.ListenAndServe(":8070", srv.Srv); err != nil {
	//		panic(err)
	//	}
	//}()

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

	return srv
}

//func TestHTTPServerAdd(t *testing.T) {
//	logg, err := logger.New("info")
//	require.NoError(t, err)
//
//	storage := memorystorage.New()
//	calendar := app.New(logg, storage)
//
//	cfg := cfg.Config{
//		HttpServ: cfg.HttpServerConf{
//			Host: "127.0.0.1",
//			Port: "18000",
//		},
//		GrpsServ: cfg.GrpcServerConf{
//			Host: "127.0.0.1",
//			Port: "9000",
//		},
//	}
//
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
//	defer cancel()
//
//	go runRest(ctx, &cfg)
//	go runGrpc(calendar, &cfg, logg)
//
//	event := bytes.NewBufferString(`{
//  "event" : {
//		"ID":"2bb0d64e-8f6e-4863-b1d8-8b20018c743d",
//		"Title":"Hello",
//		"Date":"2023-01-01T16:00:00Z",
//		"Duration": "600s",
//		"Description":"Hello",
//		"UserId":"cc526645-6fad-461e-9ebf-82a7d936a61f",
//		"Notify":"300s"
//  }
//}`)
//
//	t.Run("add", func(t *testing.T) {
//		resp, err := http.Post("http://127.0.0.1:18000/v1/event/add", "application/json", event)
//		require.NoError(t, err)
//		defer resp.Body.Close()
//		require.Equal(t, http.StatusOK, resp.StatusCode)
//		body, err := io.ReadAll(resp.Body)
//		require.NoError(t, err)
//		require.Equal(t, "{\"Error\":\"\"}", string(body))
//	})
//}

func TestHTTPServerUpdate(t *testing.T) {
	logg, err := logger.New("info")
	require.NoError(t, err)

	storage := memorystorage.New()
	calendar := app.New(logg, storage)

	cfg := cfg.Config{
		HttpServ: cfg.HttpServerConf{
			Host: "",
			Port: "64000",
		},
		GrpsServ: cfg.GrpcServerConf{
			Host: "",
			Port: "9999",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	s := startHTTP(ctx, &cfg, logg)
	startGRPC(t, ctx, &cfg, logg, calendar)

	event := `{
  "event" : {
		"ID":"2bb0d64e-8f6e-4863-b1d8-8b20018c743d",
		"Title":"Hello",
		"Date":"2023-01-01T16:00:00Z",
		"Duration": "600s",
		"Description":"Hello",
		"UserId":"cc526645-6fad-461e-9ebf-82a7d936a61f",
		"Notify":"300s"
  }
}`

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

	_ = eu

	ts := httptest.NewServer(s.Srv)
	t.Run("add", func(t *testing.T) {
		//resp, err := http.Post("http://:64000/v1/event/add", "", bytes.NewBufferString(event))

		r := httptest.NewRequest(http.MethodPost, ts.URL+"/v1/event/add", bytes.NewBufferString(event))
		w := httptest.NewRecorder()
		s.Srv.ServeHTTP(w, r)

		resp := w.Result()
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, "{\"Error\":\"\"}", string(body))
	})

	t.Run("delete", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodDelete, ts.URL+"/v1/event/delete/2bb0d64e-8f6e-4863-b1d8-8b20018c743d", nil)
		w := httptest.NewRecorder()
		s.Srv.ServeHTTP(w, r)

		resp := w.Result()

		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, "{\"Error\":\"\"}", string(body))
	})
	//
	//t.Run("update", func(t *testing.T) {
	//	c := &http.Client{}
	//	req, err := http.NewRequest(http.MethodPost, "http://:64000/v1/event/add", bytes.NewBufferString(event))
	//	require.NoError(t, err)
	//
	//	resp, err := c.Do(req)
	//	require.NoError(t, err)
	//
	//	body, err := io.ReadAll(resp.Body)
	//	require.NoError(t, err)
	//	require.Equal(t, "{\"Error\":\"\"}", string(body))
	//
	//	req, err = http.NewRequest(http.MethodPut, "http://:64000/v1/event/update", eu)
	//	require.NoError(t, err)
	//	resp, err = c.Do(req)
	//	require.NoError(t, err)
	//
	//	body, err = io.ReadAll(resp.Body)
	//	require.NoError(t, err)
	//	require.Equal(t, "{\"Error\":\"\"}", string(body))
	//
	//	req, err = http.NewRequest(http.MethodGet, "http://:64000/v1/event/get/2bb0d64e-8f6e-4863-b1d8-8b20018c743d", nil)
	//	require.NoError(t, err)
	//	resp, err = c.Do(req)
	//	require.NoError(t, err)
	//
	//	var result Event
	//	err = json.NewDecoder(resp.Body).Decode(&result)
	//	require.NoError(t, err)
	//
	//	require.Equal(t, "Hello update", result.Title)
	//})
	//
	t.Run("list day", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, ts.URL+"/v1/event/list/day", bytes.NewBufferString(`{"Date":"2023-01-01T16:00:00Z"}`))
		w := httptest.NewRecorder()
		s.Srv.ServeHTTP(w, r)

		resp := w.Result()

		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, "{\"Events\":[]}", string(body))
	})
	//
	//t.Run("list week", func(t *testing.T) {
	//	resp, err := http.Post("http://:64000/v1/event/list/week", "", bytes.NewBufferString(`{"Date":"2023-01-01T16:00:00Z"}`))
	//	require.NoError(t, err)
	//	defer resp.Body.Close()
	//	require.Equal(t, http.StatusOK, resp.StatusCode)
	//	body, err := io.ReadAll(resp.Body)
	//	require.NoError(t, err)
	//	require.Equal(t, "{\"Events\":[]}", string(body))
	//})
	//
	//t.Run("list month", func(t *testing.T) {
	//	resp, err := http.Post("http://:64000/v1/event/list/month", "", bytes.NewBufferString(`{"Date":"2023-01-01T00:00:00Z"}`))
	//	require.NoError(t, err)
	//	defer resp.Body.Close()
	//	require.Equal(t, http.StatusOK, resp.StatusCode)
	//	var result Ev
	//	err = json.NewDecoder(resp.Body).Decode(&result)
	//	require.NoError(t, err)
	//	require.Equal(t, "Hello update", result.Events[0].Title)
	//})
}

type Event struct {
	Title string
}

//func TestHTTPServerDelete(t *testing.T) {
//	logg, err := logger.New("info")
//	require.NoError(t, err)
//
//	storage := memorystorage.New()
//	calendar := app.New(logg, storage)
//
//	cfg := cfg.Config{
//		HttpServ: cfg.HttpServerConf{
//			Host: "127.0.0.1",
//			Port: "18002",
//		},
//		GrpsServ: cfg.GrpcServerConf{
//			Host: "127.0.0.1",
//			Port: "9002",
//		},
//	}
//
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
//	defer cancel()
//
//	go runRest(ctx, &cfg)
//	go runGrpc(calendar, &cfg, logg)
//
//	event := bytes.NewBufferString(`{
//  "event" : {
//		"ID":"2bb0d64e-8f6e-4863-b1d8-8b20018c743d",
//		"Title":"Hello",
//		"Date":"2023-01-01T16:00:00Z",
//		"Duration": "600s",
//		"Description":"Hello",
//		"UserId":"cc526645-6fad-461e-9ebf-82a7d936a61f",
//		"Notify":"300s"
//  }
//}`)
//
//	t.Run("delete", func(t *testing.T) {
//		resp, err := http.Post("http://127.0.0.1:18002/v1/event/add", "", event)
//		require.NoError(t, err)
//		defer resp.Body.Close()
//		require.Equal(t, http.StatusOK, resp.StatusCode)
//		body, err := io.ReadAll(resp.Body)
//		require.NoError(t, err)
//		require.Equal(t, "{\"Error\":\"\"}", string(body))
//
//		cli := http.Client{}
//		req, err := http.NewRequest(http.MethodDelete, "http://127.0.0.1:18002/v1/event/delete/2bb0d64e-8f6e-4863-b1d8-8b20018c743d", nil)
//		require.NoError(t, err)
//
//		resp, err = cli.Do(req)
//		require.NoError(t, err)
//
//		defer resp.Body.Close()
//		require.Equal(t, http.StatusOK, resp.StatusCode)
//		body, err = io.ReadAll(resp.Body)
//		require.NoError(t, err)
//		require.Equal(t, "{\"Error\":\"\"}", string(body))
//
//		resp, err = http.Get("http://127.0.0.1:18002/v1/event/get/2bb0d64e-8f6e-4863-b1d8-8b20018c743d")
//		defer resp.Body.Close()
//		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
//		var result response
//		err = json.NewDecoder(resp.Body).Decode(&result)
//		require.NoError(t, err)
//		require.Equal(t, "event not exists", result.Message)
//	})
//}

type response struct {
	Message string
}

//func TestHTTPServerListByDay(t *testing.T) {
//	logg, err := logger.New("info")
//	require.NoError(t, err)
//
//	storage := memorystorage.New()
//	calendar := app.New(logg, storage)
//
//	cfg := cfg.Config{
//		HttpServ: cfg.HttpServerConf{
//			Host: "127.0.0.1",
//			Port: "18003",
//		},
//		GrpsServ: cfg.GrpcServerConf{
//			Host: "127.0.0.1",
//			Port: "9003",
//		},
//	}
//
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
//	defer cancel()
//
//	go runRest(ctx, &cfg)
//	go runGrpc(calendar, &cfg, logg)
//
//	event := bytes.NewBufferString(`{
//  "event" : {
//		"ID":"2bb0d64e-8f6e-4863-b1d8-8b20018c743d",
//		"Title":"Hello",
//		"Date":"2023-01-01T16:00:00Z",
//		"Duration": "600s",
//		"Description":"Hello",
//		"UserId":"cc526645-6fad-461e-9ebf-82a7d936a61f",
//		"Notify":"300s"
//  }
//}`)
//
//	t.Run("list day", func(t *testing.T) {
//		resp, err := http.Post("http://127.0.0.1:18003/v1/event/add", "", event)
//		require.NoError(t, err)
//		defer resp.Body.Close()
//		require.Equal(t, http.StatusOK, resp.StatusCode)
//		body, err := io.ReadAll(resp.Body)
//		require.NoError(t, err)
//		require.Equal(t, "{\"Error\":\"\"}", string(body))
//
//		resp, err = http.Post("http://127.0.0.1:18003/v1/event/list/day", "", bytes.NewBufferString("{\n  \"Date\":\"2023-01-01T16:00:00Z\"\n}"))
//		require.NoError(t, err)
//		defer resp.Body.Close()
//		require.Equal(t, http.StatusOK, resp.StatusCode)
//		body, err = io.ReadAll(resp.Body)
//		require.NoError(t, err)
//		require.Equal(t, "{\"Events\":[]}", string(body))
//	})
//}

//func TestHTTPServerListByWeek(t *testing.T) {
//	logg, err := logger.New("info")
//	require.NoError(t, err)
//
//	storage := memorystorage.New()
//	calendar := app.New(logg, storage)
//
//	cfg := cfg.Config{
//		HttpServ: cfg.HttpServerConf{
//			Host: "127.0.0.1",
//			Port: "18004",
//		},
//		GrpsServ: cfg.GrpcServerConf{
//			Host: "127.0.0.1",
//			Port: "9004",
//		},
//	}
//
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
//	defer cancel()
//
//	go runRest(ctx, &cfg)
//	go runGrpc(calendar, &cfg, logg)
//
//	event := bytes.NewBufferString(`{
//  "event" : {
//		"ID":"2bb0d64e-8f6e-4863-b1d8-8b20018c743d",
//		"Title":"Hello",
//		"Date":"2023-01-01T16:00:00Z",
//		"Duration": "600s",
//		"Description":"Hello",
//		"UserId":"cc526645-6fad-461e-9ebf-82a7d936a61f",
//		"Notify":"300s"
//  }
//}`)
//
//	t.Run("list week", func(t *testing.T) {
//		resp, err := http.Post("http://127.0.0.1:18004/v1/event/add", "", event)
//		require.NoError(t, err)
//		defer resp.Body.Close()
//		require.Equal(t, http.StatusOK, resp.StatusCode)
//		body, err := io.ReadAll(resp.Body)
//		require.NoError(t, err)
//		require.Equal(t, "{\"Error\":\"\"}", string(body))
//
//		resp, err = http.Post("http://127.0.0.1:18004/v1/event/list/week", "", bytes.NewBufferString("{\n  \"Date\":\"2023-01-01T16:00:00Z\"\n}"))
//		require.NoError(t, err)
//		defer resp.Body.Close()
//		require.Equal(t, http.StatusOK, resp.StatusCode)
//		body, err = io.ReadAll(resp.Body)
//		require.NoError(t, err)
//		require.Equal(t, "{\"Events\":[]}", string(body))
//	})
//}

//func TestHTTPServerListByMonth(t *testing.T) {
//	logg, err := logger.New("info")
//	require.NoError(t, err)
//
//	storage := memorystorage.New()
//	calendar := app.New(logg, storage)
//
//	cfg := cfg.Config{
//		HttpServ: cfg.HttpServerConf{
//			Host: "127.0.0.1",
//			Port: "18005",
//		},
//		GrpsServ: cfg.GrpcServerConf{
//			Host: "127.0.0.1",
//			Port: "9005",
//		},
//	}
//
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
//	defer cancel()
//
//	go runRest(ctx, &cfg)
//	go runGrpc(calendar, &cfg, logg)
//
//	event := bytes.NewBufferString(`{
//  "event" : {
//		"ID":"2bb0d64e-8f6e-4863-b1d8-8b20018c743d",
//		"Title":"Hello",
//		"Date":"2023-01-01T16:00:00Z",
//		"Duration": "600s",
//		"Description":"Hello",
//		"UserId":"cc526645-6fad-461e-9ebf-82a7d936a61f",
//		"Notify":"300s"
//  }
//}`)
//
//	t.Run("list month", func(t *testing.T) {
//		resp, err := http.Post("http://127.0.0.1:18005/v1/event/add", "", event)
//		require.NoError(t, err)
//		defer resp.Body.Close()
//		require.Equal(t, http.StatusOK, resp.StatusCode)
//		body, err := io.ReadAll(resp.Body)
//		require.NoError(t, err)
//		require.Equal(t, "{\"Error\":\"\"}", string(body))
//
//		resp, err = http.Post("http://127.0.0.1:18005/v1/event/list/month", "", bytes.NewBufferString("{\n  \"Date\":\"2023-01-01T00:00:00Z\"\n}"))
//		require.NoError(t, err)
//		defer resp.Body.Close()
//		require.Equal(t, http.StatusOK, resp.StatusCode)
//		var result Ev
//		err = json.NewDecoder(resp.Body).Decode(&result)
//		require.NoError(t, err)
//		require.Equal(t, "Hello", result.Events[0].Title)
//	})
//}

type Ev struct {
	Events []Event
}
