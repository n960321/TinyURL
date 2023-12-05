package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"
	"tinyurl/internal/handler"
	"tinyurl/pkg/database"
	"tinyurl/pkg/server"
)

func main() {
	db := database.NewDatabase()
	
	svr := server.NewServer(&server.ServerConfig{
		Address: "0.0.0.0",
		Port:    "8000",
	}, handler.NewHandler(db).GetRouter())

	svr.Listen()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	svr.Shutdown(ctx)

	log.Println("shutting down")
	os.Exit(0)

}
