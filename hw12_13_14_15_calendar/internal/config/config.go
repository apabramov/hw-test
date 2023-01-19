package config

import (
	"github.com/BurntSushi/toml"
	"os"
)

type Config struct {
	Logger  LoggerConf
	Servers ServerConf
	Storage StorageConf
}

type LoggerConf struct {
	Level string
	// TODO
}

type ServerConf struct {
	Host string
	Port int
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
