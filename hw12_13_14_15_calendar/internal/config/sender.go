package config

type SenderConfig struct {
	Logger  LoggerConf
	Storage StorageConf
	Queue   QueueConf
}

type QueueConf struct {
	Dsn          string
	Type         string
	Exchange     string
	ExchangeType string
	Queue        string
}

func NewSenderCfg(cfg string) (SenderConfig, error) {
	var c SenderConfig
	err := Load(cfg, &c)
	return c, err
}
