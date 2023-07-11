package config

import "flag"

type Config struct {
  Port int
  Env string
  Key string
}

func NewConfig() Config {
	var cfg Config

	flag.IntVar(&cfg.Port, "port", 8080, "API Server Port")
	flag.StringVar(&cfg.Env, "env", "dev", "Environment (dev, prod)")
	flag.StringVar(&cfg.Key, "key", "NO_KEY", "API Key")
	flag.Parse()

	return cfg
}
