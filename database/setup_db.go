package database

import (
	"fmt"
	"log"

	"github.com/pthum/stripcontrol-golang/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB database object
var DB *gorm.DB

// ConnectDataBase set up the connection to the database
func ConnectDataBase() {
	var conn gorm.Dialector
	configString := fmt.Sprintf("%s", config.CONFIG.Database.Host)
	log.Printf("Setup %s database with %s", config.CONFIG.Database.Type, configString)
	conn = sqlite.Open(configString)

	db, err := gorm.Open(conn, &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	DB = db
}
