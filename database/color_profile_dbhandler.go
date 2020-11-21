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

// UpdateProfile updates the profile
func UpdateProfile(profile models.ColorProfile, input models.ColorProfile) (err error) {
	// calculate the difference, as gorm seem to update too much fields
	fields := FindPartialUpdateFields(profile, input)
	err = DB.Model(&profile).Debug().Select(fields).Updates(input).Error
	return
}
