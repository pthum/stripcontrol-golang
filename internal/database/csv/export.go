package csv

import (
	"log"

	"github.com/pthum/stripcontrol-golang/internal/database"
)

func (c *CSVHandler[T]) Export(dbh database.DBHandler[T]) {
	models, err := dbh.GetAll()
	if err != nil {
		panic(err)
	}
	for i := range models {
		err := c.Create(&models[i])
		if err != nil {
			log.Printf("Couldn't create model, error: %v", err)
		}
	}
	c.persistIfNecessary()
}
