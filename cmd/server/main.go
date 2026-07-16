package main

import (
	"fmt"
	"log"

	"github.com/Lil-P0ly/go_monitoring_project/internal/server"
	"github.com/Lil-P0ly/go_monitoring_project/internal/server/config"
	"github.com/Lil-P0ly/go_monitoring_project/internal/server/handler"
	models "github.com/Lil-P0ly/go_monitoring_project/internal/server/model"
)

func main() {
	cfg, err := config.ParseFlags()
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Printf("%s", cfg.Address)

	storage := models.NewMemStorage()
	h := handler.NewMSHandlerWithStorage(storage)

	srv := server.New(cfg, h)
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
