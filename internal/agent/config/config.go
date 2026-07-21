package config

import (
	"errors"
	"flag"
	"os"
	"strconv"
)

type Config struct {
	Address        string
	ReportInterval int
	PollInterval   int
	Level          string
}

func ParseFlags() (Config, error) {
	var cfg Config

	cfgFlagset := flag.NewFlagSet("cfg", flag.ContinueOnError)
	cfgFlagset.StringVar(&cfg.Address, "a", "localhost:8080", "endpoint")
	cfgFlagset.IntVar(&cfg.ReportInterval, "r", 10, "report interval")
	cfgFlagset.IntVar(&cfg.PollInterval, "p", 2, "poll interval")
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
		cfg.Level = loglevel
	}

	reportInterval, exists := os.LookupEnv("REPORT_INTERVAL")
	if exists {
		ReportIntervalInt, err := strconv.Atoi(reportInterval)
		if err != nil {
			return Config{}, errors.New("Bad enviromnet value")

		}
		cfg.ReportInterval = ReportIntervalInt
	}

	pollInterval, exists := os.LookupEnv("POLL_INTERVAL")
	if exists {
		PollIntervalInt, err := strconv.Atoi(pollInterval)
		if err != nil {
			return Config{}, errors.New("Bad enviromnet value")

		}
		cfg.PollInterval = PollIntervalInt
	}

	return cfg, nil
}
