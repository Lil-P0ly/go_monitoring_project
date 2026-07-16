package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Lil-P0ly/go_monitoring_project/internal/server/config"
	"github.com/Lil-P0ly/go_monitoring_project/internal/server/handler"
	models "github.com/Lil-P0ly/go_monitoring_project/internal/server/model"
	"github.com/go-chi/chi/v5"
)

func main() {

	cfg, err := config.ParseFlags()

	if err != nil {
		fmt.Println(err)
		return
	}
	log.Printf("%s", cfg.Address)

	r := chi.NewRouter()

	memStorage := models.NewMemStorage()

	msh := handler.NewMSHandlerWithStorage(memStorage)

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed"))
	})

	r.Get("/", msh.PrintMetricsHTML)
	r.Post("/update/{metrics_type}/{metrics_name}/{metrics_value}", msh.AddValue)
	r.Get("/value/{metrics_type}/{metrics_name}", msh.PrintLastValue)

	if err := http.ListenAndServe(cfg.Address, r); err != nil {
		log.Fatal(err)
	}

}
