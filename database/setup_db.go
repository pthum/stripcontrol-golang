package database

import (
	"fmt"
	"log"

	"github.com/pthum/stripcontrol-golang/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB database object
var DB *gorm.DB

// ConnectDataBase set up the connection to the database
func ConnectDataBase() {
	configString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable TimeZone=Europe/Berlin", config.CONFIG.Database.Username, config.CONFIG.Database.Password, config.CONFIG.Database.DbName, config.CONFIG.Database.Host, config.CONFIG.Database.Port)
	log.Printf("Setup database with %s", configString)
	db, err := gorm.Open(postgres.Open(configString), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	DB = db
}
