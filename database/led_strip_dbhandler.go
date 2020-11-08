package database

import (
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
