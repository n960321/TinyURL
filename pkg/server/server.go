package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type Server struct {
	http.Server
}

type Config struct {
	Address string  `mapstructure:"address"`
	Port    *string `mapstructure:"port"`
}

func NewServer(cfg *Config, handler *mux.Router) *Server {
	return &Server{
		http.Server{
			Addr:    cfg.Address + ":" + *cfg.Port,
			Handler: handler,
		},
	}
}

func (svr *Server) Listen() {
	go func() {
		log.Info().Msgf("Server Start Listening on %s", svr.Addr)
		if err := svr.ListenAndServe(); err != nil {
			log.Fatal().Err(err)
		}
	}()
}
