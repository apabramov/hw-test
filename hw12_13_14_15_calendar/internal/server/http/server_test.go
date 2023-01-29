package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"google.golang.org/grpc"
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

type Event struct {
	Title string
}

type Response struct {
	Ev    Event `json:"event,omitempty"`
	Error string
}

type Ev struct {
	Events []Event
}

func start(t *testing.T, ctx context.Context, cfg *cfg.Config, logg *logger.Logger, a *app.App) *Server {
	l, err := net.Listen("tcp", net.JoinHostPort(cfg.GrpsServ.Host, cfg.GrpsServ.Port))
	require.NoError(t, err)

	clientOptions := []grpc.DialOption{grpc.WithInsecure()}
	_, err = grpc.Dial(l.Addr().String(), clientOptions...)
	require.NoError(t, err)

	g := internalgrpc.NewServer(logg, a, cfg.GrpsServ)
	h, err := NewServer(ctx, logg, cfg)
	require.NoError(t, err)

	go func() {
		err = g.Srv.Serve(l)
		require.NoError(t, err)
		err = h.Start(ctx)
		require.NoError(t, err)
	}()

	go func() {
		<-ctx.Done()
		g.Stop()
		h.Stop(ctx)
	}()

	return h
}

func TestHTTPServer(t *testing.T) {
	logg, err := logger.New("info")
	require.NoError(t, err)

	storage := memorystorage.New()
	calendar := app.New(logg, storage)

	cfg := cfg.Config{
		HttpServ: cfg.HttpServerConf{
			Host: "127.0.0.1",
			Port: "8080",
		},
		GrpsServ: cfg.GrpcServerConf{
			Host: "127.0.0.1",
			Port: "50051",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	s := start(t, ctx, &cfg, logg, calendar)

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

	ts := httptest.NewServer(s.Srv.Handler)
	t.Run("add", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, ts.URL+"/v1/event/add", bytes.NewBufferString(event))
		w := httptest.NewRecorder()
		s.Srv.Handler.ServeHTTP(w, r)

		resp := w.Result()
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var res Response
		err = json.NewDecoder(resp.Body).Decode(&res)
		require.NoError(t, err)
		require.Equal(t, "Hello", res.Ev.Title)
	})

	t.Run("delete", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodDelete, ts.URL+"/v1/event/delete/2bb0d64e-8f6e-4863-b1d8-8b20018c743d", nil)
		w := httptest.NewRecorder()
		s.Srv.Handler.ServeHTTP(w, r)
		resp := w.Result()

		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		var res Response
		err = json.NewDecoder(resp.Body).Decode(&res)
		require.NoError(t, err)
		require.Equal(t, "", res.Error)
	})

	t.Run("update", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, ts.URL+"/v1/event/add", bytes.NewBufferString(event))
		w := httptest.NewRecorder()
		s.Srv.Handler.ServeHTTP(w, r)
		resp := w.Result()

		var res Response
		err = json.NewDecoder(resp.Body).Decode(&res)
		require.NoError(t, err)

		require.Equal(t, "Hello", res.Ev.Title)

		r = httptest.NewRequest(http.MethodPut, ts.URL+"/v1/event/update", eu)
		w = httptest.NewRecorder()
		s.Srv.Handler.ServeHTTP(w, r)
		resp = w.Result()

		err = json.NewDecoder(resp.Body).Decode(&res)
		require.NoError(t, err)

		require.Equal(t, "Hello update", res.Ev.Title)

		r = httptest.NewRequest(http.MethodGet, ts.URL+"/v1/event/get/2bb0d64e-8f6e-4863-b1d8-8b20018c743d", nil)
		w = httptest.NewRecorder()
		s.Srv.Handler.ServeHTTP(w, r)
		resp = w.Result()

		err = json.NewDecoder(resp.Body).Decode(&res)
		require.NoError(t, err)

		require.Equal(t, "Hello update", res.Ev.Title)
	})

	t.Run("list day", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, ts.URL+"/v1/event/list/day", bytes.NewBufferString(`{"Date":"2023-01-01T16:00:00Z"}`))
		w := httptest.NewRecorder()
		s.Srv.Handler.ServeHTTP(w, r)

		resp := w.Result()

		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		var res Response
		err = json.NewDecoder(resp.Body).Decode(&res)
		require.NoError(t, err)
		require.Equal(t, "", res.Error)
	})

	t.Run("list week", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, ts.URL+"/v1/event/list/week", bytes.NewBufferString(`{"Date":"2023-01-01T16:00:00Z"}`))
		w := httptest.NewRecorder()
		s.Srv.Handler.ServeHTTP(w, r)
		resp := w.Result()

		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		var res Response
		err = json.NewDecoder(resp.Body).Decode(&res)
		require.NoError(t, err)
		require.Equal(t, "", res.Error)
	})

	t.Run("list month", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, ts.URL+"/v1/event/list/month", bytes.NewBufferString(`{"bg":"2023-01-01T00:00:00Z","fn":"2023-02-01T00:00:00Z"}`))
		w := httptest.NewRecorder()
		s.Srv.Handler.ServeHTTP(w, r)
		resp := w.Result()

		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		var result Ev
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		require.Equal(t, "Hello update", result.Events[0].Title)
	})
}
