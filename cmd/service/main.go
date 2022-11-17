package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pthum/stripcontrol-golang/internal/api"
	"github.com/pthum/stripcontrol-golang/internal/config"
	"github.com/pthum/stripcontrol-golang/internal/database"
	messagingimpl "github.com/pthum/stripcontrol-golang/internal/messaging/impl"

	flag "github.com/spf13/pflag"
)

var configFile string

func init() {
	flag.StringVarP(&configFile, "config", "c", "", "this is the path and filename to the config file")
}

func main() {

	flag.Parse()
	cfg, err := config.InitConfig(configFile)
	if err != nil {
		log.Fatalf("Error initializing config: %v", err)
	}

	var enableDebug = cfg.Server.Mode != "release"
	dbh := database.New(cfg.Database)
	defer dbh.Close()

	mh := messagingimpl.New(cfg.Messaging)
	defer mh.Close()

	router := api.NewRouter(dbh, mh, enableDebug)

	// Listen and serve on 0.0.0.0:8080
	serve := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         serve,
		Handler:      router,
	}
	go func() {
		// panic(server.ListenAndServe())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen:%+s\n", err)
		}
	}()

	// Create channel for shutdown signals.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)

	//Recieve shutdown signals.
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("error shutting down server %s", err)
	} else {
		log.Println("Server gracefully stopped")
	}
}
