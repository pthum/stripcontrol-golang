package main

import (
	"log"

	"github.com/pthum/stripcontrol-golang/internal/config"
	"github.com/pthum/stripcontrol-golang/internal/database"
	"github.com/pthum/stripcontrol-golang/internal/database/csv"
	"github.com/pthum/stripcontrol-golang/internal/model"

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

	cpDbh := database.New[model.ColorProfile](cfg.Database)
	defer cpDbh.Close()
	cpCsv := csv.NewHandler[model.ColorProfile](&cfg.CSV)
	cpCsv.Export(cpDbh)

	lsDbh := database.New[model.LedStrip](cfg.Database)
	defer lsDbh.Close()
	lsCsv := csv.NewHandler[model.LedStrip](&cfg.CSV)
	lsCsv.Export(lsDbh)
}
