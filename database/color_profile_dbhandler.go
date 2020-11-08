package database

import "github.com/pthum/stripcontrol-golang/models"

// GetAllColorProfiles get all color profiles
func GetAllColorProfiles() (profiles []models.ColorProfile, err error) {
	err = DB.Find(&profiles).Error
	return
}

// GetColorProfile loads a color profile from the database
func GetColorProfile(ID string) (profile models.ColorProfile, err error) {
	err = DB.Where("id = ?", ID).First(&profile).Error
	return
}
