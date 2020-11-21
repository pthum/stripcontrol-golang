package database

import (
	"fmt"
	"log"
	"reflect"

	"github.com/pthum/stripcontrol-golang/models"
)

// GetAllLedStrips get all LED strips
func GetAllLedStrips() (strips []models.LedStrip, err error) {
	err = DB.Find(&strips).Error
	return
}

// GetLedStrip loads a strip from the database
func GetLedStrip(ID string) (strip models.LedStrip, err error) {
	err = DB.Where("id = ?", ID).First(&strip).Error
	return
}

// UpdateStrip updates the strip
func UpdateStrip(strip models.LedStrip, input models.LedStrip) (err error) {
	// calculate the difference, as gorm seem to update too much fields
	fields := FindPartialUpdateFields(strip, input)
	err = DB.Model(&strip).Debug().Select(fields).Updates(input).Error
	return
}

// FindPartialUpdateFields find the fields that need to be updated
func FindPartialUpdateFields(dbObject interface{}, input interface{}) (fields []string) {
	tIn := reflect.TypeOf(input)
	tDb := reflect.TypeOf(dbObject)
	if tIn.Kind() != tDb.Kind() || tIn != tDb {
		log.Println("different kinds no update")
		return
	}
	valIn := reflect.ValueOf(input)
	valDb := reflect.ValueOf(dbObject)

	for i := 0; i < valIn.NumField(); i++ {
		valueFieldIn := valIn.Field(i)
		valueFieldDb := valDb.Field(i)
		typeField := valIn.Type().Field(i)
		if valueFieldIn.Interface() != valueFieldDb.Interface() {
			fields = append(fields, typeField.Name)
		}
	}

	fmt.Printf("Fields to update: %v\n", fields)
	return
}
