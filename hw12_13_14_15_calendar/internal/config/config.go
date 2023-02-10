package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Logger   LoggerConf
	HttpServ HttpServerConf
	GrpsServ GrpcServerConf
	Storage  StorageConf
}

type LoggerConf struct {
	Level string
}

type HttpServerConf struct {
	Host string
	Port string
}

type GrpcServerConf struct {
	Host string
	Port string
}

type StorageConf struct {
	Type string
	Dsn  string
}

func NewCalenderCfg(cfg string) (Config, error) {
	c := Config{}
	err := Load(cfg, &c)
	return c, err
}

func Load(cfg string, conf interface{}) error {
	f, err := os.ReadFile(cfg)
	if err != nil {
		return err
	}
	if err := toml.Unmarshal(f, conf); err != nil {
		return err
	}
	return nil
}
