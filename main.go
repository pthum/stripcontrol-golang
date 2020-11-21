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

	"github.com/pthum/stripcontrol-golang/config"
	"github.com/pthum/stripcontrol-golang/database"
	"github.com/pthum/stripcontrol-golang/mappings"
	"github.com/pthum/stripcontrol-golang/messaging"
)

func main() {
	config.InitConfig()

	var enableDebug = config.CONFIG.Server.Mode != "release"

	router := mappings.NewRouter(enableDebug)
	messaging.Init()
	defer messaging.Close()

	database.ConnectDataBase()
	defer database.CloseDB()

	// Listen and serve on 0.0.0.0:8080
	serve := fmt.Sprintf("%s:%s", config.CONFIG.Server.Host, config.CONFIG.Server.Port)
	server := &http.Server{Addr: serve, Handler: router}
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
