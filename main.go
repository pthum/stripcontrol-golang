package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pthum/stripcontrol-golang/config"
	"github.com/pthum/stripcontrol-golang/database"
	"github.com/pthum/stripcontrol-golang/mappings"
	"github.com/pthum/stripcontrol-golang/messaging"
)

func main() {
	config.InitConfig()
	// set release mode if set
	if config.CONFIG.Server.Mode == "release" {
		// gin.SetMode(gin.ReleaseMode)
	}
	router := mappings.NewRouter()
	messaging.Init()
	database.ConnectDataBase()
	// Listen and serve on 0.0.0.0:8080
	serve := fmt.Sprintf("%s:%s", config.CONFIG.Server.Host, config.CONFIG.Server.Port)
	log.Fatal(http.ListenAndServe(serve, router))

}
