package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"time"
)

type Logger struct {
	Log *zap.Logger
}

func New(level string) (*Logger, error) {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	en := zapcore.NewJSONEncoder(config)

	al, err := zap.ParseAtomicLevel(level)

	core := zapcore.NewTee(
		zapcore.NewCore(en, zapcore.AddSync(os.Stdout), al),
	)

	l := zap.New(core)

	log := l.Named("calendar")
	//log, err := zap.NewProduction()
	//if err != nil {
	//	return nil, err
	//}
	//defer log.Sync()
	//
	//_, err = zap.ParseAtomicLevel(level)
	//if err != nil {
	//	return nil, err
	//}
	return &Logger{log}, err
}

func (l Logger) Info(msg string) {
	l.Log.Info(msg)
}

func (l Logger) InfoHttp(r *http.Request, httpStatus int, duration time.Duration) {
	l.Log.Info(fmt.Sprintf("%s [%s] %s %s %s %v %s %s", r.RemoteAddr, time.Now().Format(time.RFC822Z),
		r.Method, r.RequestURI, r.URL.Scheme, httpStatus, duration, r.Header.Get("User-Agent")))
}

func (l Logger) Error(msg string) {
	l.Log.Error(msg)
}

func (l Logger) Warn(msg string) {
	l.Log.Warn(msg)
}

func (l Logger) Debug(msg string) {
	l.Log.Debug(msg)
}
