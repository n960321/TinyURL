package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	http.Server
}

type ServerConfig struct {
	Address string
	Port    string
}

func NewServer(cfg *ServerConfig, handler *mux.Router) *Server {
	return &Server{
		http.Server{
			Addr: cfg.Address + ":" + cfg.Port,
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
