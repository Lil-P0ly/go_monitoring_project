package main

import (
	"fmt"

	"github.com/Lil-P0ly/go_monitoring_project/internal/agent"
	"github.com/Lil-P0ly/go_monitoring_project/internal/agent/config"
	"github.com/Lil-P0ly/go_monitoring_project/internal/agent/logger"
	"github.com/Lil-P0ly/go_monitoring_project/internal/agent/model"
)

func main() {
	cfg, err := config.ParseFlags()
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := logger.InitLogger("info"); err != nil {
		fmt.Println(err)
		return
	}
	defer logger.Sync()
	logger.Infof("%s %d %d", cfg.Address, cfg.PollInterval, cfg.ReportInterval)
	m := model.NewMetricsMap()
	agent.New(cfg, m).Run()
}
