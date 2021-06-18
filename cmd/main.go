package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/Striker87/members-club"
	"github.com/Striker87/members-club/router"
	"github.com/Striker87/members-club/storage"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error init config: %s", err)
	}

	store := make(map[string]storage.User)
	route := router.Set()

	srv := members.Run(viper.GetString("port"), route, store)
	go func() {
		fmt.Println("Server started at port:", viper.GetString("port"))

		if err := srv.HttpServer.ListenAndServe(); err != nil {
			log.Fatalf("failed to listen and serve: %+v", err)
		}
	}()

	gracefulShutdown(srv)
}

func initConfig() error {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}

func gracefulShutdown(srv members.Server) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	<-quit

	log.Println("received terminate, graceful shutdown")

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	_ = srv.HttpServer.Shutdown(ctx)
}
