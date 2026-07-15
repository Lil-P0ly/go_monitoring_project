package agent

import (
	"errors"
	"flag"
	"os"
	"time"
)

type Config struct {
	Address        string
	ReportInterval time.Duration
	PollInterval   time.Duration
}

func ParseFlags() (Config, error) {
	var cfg Config

	cfgFlagset := flag.NewFlagSet("cfg", flag.ContinueOnError)
	cfgFlagset.StringVar(&cfg.Address, "a", "localhost:8080", "endpoint")
	cfgFlagset.DurationVar(&cfg.ReportInterval, "r", 10*time.Second, "report interval")
	cfgFlagset.DurationVar(&cfg.PollInterval, "p", 2*time.Second, "poll interval")
	err := cfgFlagset.Parse(os.Args[1:])

	if err != nil {
		return Config{}, errors.New("Unknow argument")
	}
	return cfg, nil
}
