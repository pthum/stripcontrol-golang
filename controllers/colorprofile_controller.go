package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pthum/null"
	"github.com/pthum/stripcontrol-golang/database"
	"github.com/pthum/stripcontrol-golang/messaging"
	"github.com/pthum/stripcontrol-golang/models"
	"github.com/pthum/stripcontrol-golang/utils"
)

const (
	profileNotFoundMsg = "Profile not found!"
)

// GetAllColorProfiles get all color profiles
func GetAllColorProfiles(w http.ResponseWriter, r *http.Request) {
	var profiles []models.ColorProfile
	database.DB.Find(&profiles)

	HandleJSON(&w, http.StatusOK, profiles)
}

// GetColorProfile get a specific color profile
func GetColorProfile(w http.ResponseWriter, r *http.Request) {
	// Get model if exist
	var profile, err = database.GetColorProfile(GetParam(r, "id"))
	if err != nil {
		HandleJSON(&w, http.StatusNotFound, H{"error": profileNotFoundMsg})
		return
	}
	HandleJSON(&w, http.StatusOK, profile)
}

// CreateColorProfile create a color profile
func CreateColorProfile(w http.ResponseWriter, r *http.Request) {
	// Validate input
	var input models.ColorProfile
	if err := BindJSON(r, &input); err != nil {
		HandleJSON(&w, http.StatusBadRequest, H{"error": err.Error()})
		return
	}

	// generate an id
	input.ID = utils.GenerateID()
	if err := database.DB.Create(&input).Error; err != nil {
		log.Printf("Error: %s", err)
		HandleJSON(&w, http.StatusBadRequest, H{"error": err.Error()})
		return
	}
	w.Header().Add("Location", fmt.Sprintf("%s/%d", r.RequestURI, input.ID))
	HandleJSON(&w, http.StatusCreated, input)
}

// UpdateColorProfile update a color profile
func UpdateColorProfile(w http.ResponseWriter, r *http.Request) {
	// Get model if exist
	var profile, err = database.GetColorProfile(GetParam(r, "id"))
	if err != nil {
		HandleJSON(&w, http.StatusNotFound, H{"error": profileNotFoundMsg})
		return
	}

	// Validate input
	var input models.ColorProfile
	if err := BindJSON(r, &input); err != nil {
		HandleJSON(&w, http.StatusBadRequest, H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&input).Error; err != nil {
		HandleJSON(&w, http.StatusBadRequest, H{"error": err.Error()})
		return
	}
	go messaging.PublishProfileSaveEvent(null.NewInt(input.ID, true), input)

	HandleJSON(&w, http.StatusOK, profile)
}

// DeleteColorProfile delete a color profile
func DeleteColorProfile(w http.ResponseWriter, r *http.Request) {
	// Get model if exist
	var profile, err = database.GetColorProfile(GetParam(r, "id"))
	if err != nil {
		HandleJSON(&w, http.StatusNotFound, H{"error": profileNotFoundMsg})
		return
	}
	if err := database.DB.Delete(&profile).Error; err != nil {
		HandleJSON(&w, http.StatusBadRequest, H{"error": err.Error()})
		return
	}

	go messaging.PublishProfileDeleteEvent(null.NewInt(profile.ID, true))
	HandleJSON(&w, http.StatusNoContent, nil)
}
