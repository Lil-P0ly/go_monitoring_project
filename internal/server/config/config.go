package config

import (
	"errors"
	"flag"
	"os"
)

type Config struct {
	Address string
}

func ParseFlags() (Config, error) {
	var cfg Config

	cfgFlagset := flag.NewFlagSet("cfg", flag.ContinueOnError)
	cfgFlagset.StringVar(&cfg.Address, "a", "localhost:8080", "endpoint")
	err := cfgFlagset.Parse(os.Args[1:])

	if err != nil {
		return Config{}, errors.New("Unknow argument")
	}
	return cfg, nil
}
