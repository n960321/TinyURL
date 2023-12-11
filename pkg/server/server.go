package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	http.Server
}

type Config struct {
	Address *string `mapstructure:"address"`
	Port    *string `mapstructure:"port"`
}

func NewServer(cfg *Config, handler *mux.Router) *Server {
	return &Server{
		http.Server{
			Addr:    *cfg.Address + ":" + *cfg.Port,
			Handler: handler,
		},
	}
}

func (svr *Server) Listen() {
	go func() {
		if err := svr.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
}
