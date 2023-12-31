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
	local      bool
)

func main() {
	flag.BoolVar(&local, "local", false, "Run on local.")
	flag.StringVar(&configFile, "config", "configs/config.yaml", "The config file.")
	flag.Parse()

	logger.SetLogger(local)
	config := config.GetConfig(configFile)
	db := database.NewDatabase(config.DB,local)
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
