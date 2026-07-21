package server

import (
	"net/http"

	"github.com/Lil-P0ly/go_monitoring_project/internal/server/config"
	"github.com/Lil-P0ly/go_monitoring_project/internal/server/handler"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	cfg     config.Config
	handler *handler.MemoryStorageHandler
	router  *chi.Mux
}

func New(cfg config.Config, h *handler.MemoryStorageHandler) *Server {
	s := &Server{
		cfg:     cfg,
		handler: h,
		router:  chi.NewRouter(),
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed"))
	})

	// middleware
	s.router.Use(s.handler.LogRequestMiddleware)

	// handlers
	s.router.Get("/", s.handler.PrintMetricsHTML)
	s.router.Post("/update/{metrics_type}/{metrics_name}/{metrics_value}", s.handler.AddValue)
	s.router.Get("/value/{metrics_type}/{metrics_name}", s.handler.PrintLastValue)
}

func (s *Server) Run() error {
	return http.ListenAndServe(s.cfg.Address, s.router)
}
