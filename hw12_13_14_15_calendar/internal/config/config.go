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
	// TODO
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

func NewConfig(cfg string) (Config, error) {
	var conf Config
	f, err := os.ReadFile(cfg)
	if err != nil {
		return Config{}, err
	}

	if _, err := toml.Decode(string(f), &conf); err != nil {
		return Config{}, err
	}
	return conf, nil
}

// TODO
