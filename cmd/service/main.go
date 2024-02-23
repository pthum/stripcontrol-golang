package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/pthum/stripcontrol-golang/internal/api"
	"github.com/pthum/stripcontrol-golang/internal/config"
	"github.com/pthum/stripcontrol-golang/internal/database/csv"
	alog "github.com/pthum/stripcontrol-golang/internal/log"
	messagingimpl "github.com/pthum/stripcontrol-golang/internal/messaging/impl"
	"github.com/pthum/stripcontrol-golang/internal/model"
	"github.com/pthum/stripcontrol-golang/internal/service"
	"github.com/pthum/stripcontrol-golang/internal/telegram"
	"github.com/samber/do"

	flag "github.com/spf13/pflag"
)

var configFile string

func init() {
	flag.StringVarP(&configFile, "config", "c", "", "this is the path and filename to the config file")
}

func main() {

	flag.Parse()
	cfg, err := config.InitConfig(configFile)
	l := alog.NewLogger("main")
	if err != nil {
		l.Error("Error initializing config: %v", err)
	}
	var enableDebug = cfg.Server.Mode != "release"

	inj := do.New()
	defer inj.Shutdown()
	do.ProvideValue(inj, cfg)
	do.Provide(inj, newScheduler)
	do.Provide(inj, csv.NewHandlerI[model.ColorProfile])
	do.Provide(inj, csv.NewHandlerI[model.LedStrip])
	do.Provide(inj, messagingimpl.New)
	do.Provide(inj, service.NewCPService)
	do.Provide(inj, service.NewLEDService)
	do.Provide(inj, api.NewCPHandler)
	do.Provide(inj, api.NewLEDHandler)

	tgH := telegram.NewHandler(inj, cfg.Telegram)
	go tgH.Handle()

	router := api.NewRouter(inj, enableDebug)

	// Listen and serve on 0.0.0.0:8080
	serve := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         serve,
		Handler:      router,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Warn("listen:%+s\n", err)
		}
	}()
	scheduleJobs(inj)
	// Create channel for shutdown signals.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)

	//Recieve shutdown signals.
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		l.Error("error shutting down server %s", err)
	} else {
		l.Info("Server gracefully stopped")
	}
}

func newScheduler(inj *do.Injector) (*gocron.Scheduler, error) {
	return gocron.NewScheduler(time.UTC), nil
}

func scheduleJobs(inj *do.Injector) {
	s := do.MustInvoke[*gocron.Scheduler](inj)
	// start scheduler
	s.StartAsync()
}
