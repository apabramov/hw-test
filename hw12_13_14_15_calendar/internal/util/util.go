package util

import (
	"context"
	"fmt"

	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/app"
	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/logger"
	memorystorage "github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/storage/sql"
)

func NewStorage(log *logger.Logger, cfg config.StorageConf) app.Storage {
	var st app.Storage
	switch cfg.Type {
	case "memory":
		st = memorystorage.New()
	case "sql":
		st = sqlstorage.New(log, cfg)
		err := st.Connect(context.Background())
		if err != nil {
			log.Info(fmt.Sprintf("NewStorage - %s", err.Error()))
			return nil
		}
	default:
		log.Error(fmt.Sprintf("storage type not found - %s", cfg.Type))
	}
	return st
}
