package database

import (
	"fmt"
	"log"

	"github.com/pthum/stripcontrol-golang/internal/config"
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

// CloseDB closing db
func CloseDB() {
	// Get generic database object sql.DB to be able to close
	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("err closing db connection: %s", err.Error())
	}
	if err = sqlDB.Close(); err != nil {
		log.Printf("err closing db connection: %s", err.Error())
	} else {
		log.Println("db connection gracefully closed")
	}
}
