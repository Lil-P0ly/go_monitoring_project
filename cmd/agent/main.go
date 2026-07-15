package main

import (
	"fmt"
	"log"

	"github.com/Lil-P0ly/go_monitoring_project/internal/agent"
	"github.com/Lil-P0ly/go_monitoring_project/internal/agent/model"
)

func main() {
	cfg, err := agent.ParseFlags()
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Printf("%s %d %d", cfg.Address, cfg.PollInterval, cfg.ReportInterval)
	m := model.NewMetricsMap()
	agent.New(cfg, m).Run()
}
