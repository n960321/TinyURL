package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"time"
	"tinyurl/internal/config"
	"tinyurl/internal/handler"
	"tinyurl/pkg/database"
	"tinyurl/pkg/logger"
	"tinyurl/pkg/redis"
	"tinyurl/pkg/server"

	"github.com/rs/zerolog/log"
)

var (
	configFile string
)

func main() {
	flag.StringVar(&configFile, "config", "configs/config.yaml", "The config file.")
	flag.Parse()

	config := config.GetConfig(configFile)
	logger.SetLogger(*config.Local)
	db := database.NewDatabase(config.DB)
	redis := redis.NewRedisCache(config.Cache)
	svr := server.NewServer(config.Http, handler.NewHandler(db, redis).GetRouter())

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
