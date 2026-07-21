package config

import (
	"errors"
	"flag"
	"os"
)

type Config struct {
	Address string
	Level   string
}

func ParseFlags() (Config, error) {
	var cfg Config

	cfgFlagset := flag.NewFlagSet("cfg", flag.ContinueOnError)
	cfgFlagset.StringVar(&cfg.Address, "a", "localhost:8080", "endpoint")
	cfgFlagset.StringVar(&cfg.Level, "l", "info", "loglevel")

	err := cfgFlagset.Parse(os.Args[1:])

	if err != nil {
		return Config{}, errors.New("Unknow argument")
	}

	addr, exists := os.LookupEnv("ADDRESS")
	if exists {
		cfg.Address = addr
	}

	loglevel, exists := os.LookupEnv("LEVEL")
	if exists {
		cfg.Address = loglevel
	}

	return cfg, nil
}
