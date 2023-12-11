package main

import (
	"context"
	"os"
	"os/signal"
	"time"
	"tinyurl/internal/handler"
	"tinyurl/pkg/database"
	"tinyurl/pkg/logger"
	"tinyurl/pkg/redis"
	"tinyurl/pkg/server"

	"github.com/rs/zerolog/log"
)

func main() {
	logger.SetLogger()
	db := database.NewDatabase()
	redis := redis.NewRedisCache()
	svr := server.NewServer(
		&server.ServerConfig{
			Address: "0.0.0.0",
			Port:    "8000",
		},
		handler.NewHandler(db, redis).GetRouter(),
	)

	svr.Listen()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	svr.Shutdown(ctx)
	log.Info().Msg("shutting down")
	os.Exit(0)

}
