package main

import (
	"log"
	"net/http"

	"github.com/Lil-P0ly/go_monitoring_project/internal/server/handler"
	models "github.com/Lil-P0ly/go_monitoring_project/internal/server/model"
)

func main() {
	mux := http.NewServeMux()

	memStorage := models.NewMemStorage()

	msh := handler.NewMSHandlerWithStorage(memStorage)

	mux.HandleFunc("/update/{metrics_type}/{metrics_name}/{metrics_value}", msh.AddValue)

	mux.HandleFunc("/update/{metrics_type}/", msh.NotFound)

	mux.HandleFunc("/metrics", msh.PrintMetrics)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}

}
