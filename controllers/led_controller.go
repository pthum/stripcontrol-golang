package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/pthum/null"
	"github.com/pthum/stripcontrol-golang/database"
	"github.com/pthum/stripcontrol-golang/messaging"
	"github.com/pthum/stripcontrol-golang/models"
	"github.com/pthum/stripcontrol-golang/utils"
)

const (
	stripNotFoundMsg = "LEDStrip not found!"
)

// GetAllLedStrips get all existing led strips
func GetAllLedStrips(w http.ResponseWriter, r *http.Request) {
	var strips, err = database.GetAllLedStrips()
	if err != nil {
		HandleJSON(&w, http.StatusNotFound, H{"error": err.Error()})
		return
	}

	HandleJSON(&w, http.StatusOK, strips)
}

// GetLedStrip get a single led strip
func GetLedStrip(w http.ResponseWriter, r *http.Request) {
	// Get model if exist
	var strip, err = database.GetLedStrip(GetParam(r, "id"))
	if err != nil {
		HandleJSON(&w, http.StatusNotFound, H{"error": stripNotFoundMsg})
		return
	}
	HandleJSON(&w, http.StatusOK, strip)
}

// CreateLedStrip create an LED strip
func CreateLedStrip(w http.ResponseWriter, r *http.Request) {
	// Validate input
	var input models.LedStrip
	if err := BindJSON(r, &input); err != nil {
		HandleJSON(&w, http.StatusBadRequest, H{"error": err.Error()})
		return
	}
	// generate an id
	input.ID = utils.GenerateID()
	log.Printf("Generated ID %d", input.ID)
	if err := database.DB.Create(&input).Error; err != nil {
		log.Printf("Error: %s", err)
		HandleJSON(&w, http.StatusBadRequest, H{"error": err.Error()})
		return
	}

	go messaging.PublishStripSaveEvent(null.NewInt(0, false), input)
	log.Printf("ID after save %d", input.ID)
	w.Header().Add("Location", fmt.Sprintf("%s/%d", r.RequestURI, input.ID))
	HandleJSON(&w, http.StatusCreated, input)
}

// UpdateLedStrip update an LED strip
func UpdateLedStrip(w http.ResponseWriter, r *http.Request) {
	// Get model if exist
	var strip, err = database.GetLedStrip(GetParam(r, "id"))
	if err != nil {
		HandleJSON(&w, http.StatusNotFound, H{"error": stripNotFoundMsg})
		return
	}

	// Validate input
	var input models.LedStrip
	if err := BindJSON(r, &input); err != nil {
		HandleJSON(&w, http.StatusBadRequest, H{"error": err.Error()})
		return
	}
	// profile shouldn't be updated through this endpoint
	input.ProfileID = strip.ProfileID

	if err := database.DB.Save(&input).Error; err != nil {
		HandleJSON(&w, http.StatusBadRequest, H{"error": err.Error()})
		return
	}

	go messaging.PublishStripSaveEvent(null.NewInt(input.ID, true), input)

	HandleJSON(&w, http.StatusNoContent, nil)
}

// DeleteLedStrip delete an LED strip
func DeleteLedStrip(w http.ResponseWriter, r *http.Request) {
	// Get model if exist
	var strip, err = database.GetLedStrip(GetParam(r, "id"))
	if err != nil {
		HandleJSON(&w, http.StatusNotFound, H{"error": stripNotFoundMsg})
		return
	}

	if err := database.DB.Delete(&strip).Error; err != nil {
		HandleJSON(&w, http.StatusBadRequest, H{"error": err.Error()})
		return
	}
	go messaging.PublishStripDeleteEvent(null.NewInt(strip.ID, true))
	HandleJSON(&w, http.StatusNoContent, nil)
}

// UpdateProfileForStrip update which profile is referenced to the strip
func UpdateProfileForStrip(w http.ResponseWriter, r *http.Request) {
	// Get model if exist
	var strip, err1 = database.GetLedStrip(GetParam(r, "id"))
	if err1 != nil {
		HandleJSON(&w, http.StatusNotFound, H{"error": err1.Error()})
		return
	}

	// Validate input
	var input models.ColorProfile
	if err := BindJSON(r, &input); err != nil {
		HandleJSON(&w, http.StatusBadRequest, H{"error": err.Error()})
		return
	}
	var profile, err2 = database.GetColorProfile(strconv.FormatInt(input.ID, 10))
	if err2 != nil {
		HandleJSON(&w, http.StatusNotFound, H{"error": err2.Error()})
		return
	}
	strip.ProfileID = null.NewInt(profile.ID, true)
	database.DB.Save(strip)

	go messaging.PublishStripSaveEvent(null.NewInt(strip.ID, true), strip)

	HandleJSON(&w, http.StatusOK, profile)
}

// GetProfileForStrip get the current profile of a strip
func GetProfileForStrip(w http.ResponseWriter, r *http.Request) {
	// Get model if exist
	var strip, err1 = database.GetLedStrip(GetParam(r, "id"))
	if err1 != nil {
		HandleJSON(&w, http.StatusNotFound, H{"error": err1.Error()})
		return
	}
	if !strip.ProfileID.Valid {
		HandleJSON(&w, http.StatusNotFound, H{"error": "Record not found!"})
		return
	}
	var profile, err2 = database.GetColorProfile(strconv.FormatInt(strip.ProfileID.Int64, 10))
	if err2 != nil {
		HandleJSON(&w, http.StatusNotFound, H{"error": err2.Error()})
		return
	}

	HandleJSON(&w, http.StatusOK, profile)
}

// RemoveProfileForStrip remove the current referenced profile
func RemoveProfileForStrip(w http.ResponseWriter, r *http.Request) {
	// Get model if exist
	var strip, err1 = database.GetLedStrip(GetParam(r, "id"))
	if err1 != nil {
		HandleJSON(&w, http.StatusNotFound, H{"error": err1.Error()})
		return
	}
	strip.ProfileID.Valid = false
	database.DB.Save(strip)

	go messaging.PublishStripSaveEvent(null.NewInt(strip.ID, true), strip)

	HandleJSON(&w, http.StatusNoContent, nil)
}
