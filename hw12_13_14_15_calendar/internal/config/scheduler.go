package config

type SchedulerConfig struct {
	Logger  LoggerConf
	Storage StorageConf
	Queue   QueueConf
	Ticker  TickerConf
}

type TickerConf struct {
	Duration string
}

func NewSchedulerCfg(cfg string) (SchedulerConfig, error) {
	var c SchedulerConfig
	err := Load(cfg, &c)
	return c, err
}
