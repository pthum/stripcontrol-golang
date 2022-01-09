package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pthum/null"
	api "github.com/pthum/stripcontrol-golang/internal/api/common"
	"github.com/pthum/stripcontrol-golang/internal/database"
	"github.com/pthum/stripcontrol-golang/internal/messaging"
	"github.com/pthum/stripcontrol-golang/internal/models"
	"github.com/pthum/stripcontrol-golang/internal/utils"
)

const (
	profileNotFoundMsg = "Profile not found!"
	profilePath        = "/api/colorprofile"
	profileIDPath      = profilePath + "/{id}"
)

func ColorProfileRoutes() []api.Route {
	return []api.Route{

		{"GetColorprofiles", http.MethodGet, profilePath, GetAllColorProfiles},

		{"CreateColorprofile", http.MethodPost, profilePath, CreateColorProfile},

		{"GetColorprofile", http.MethodGet, profileIDPath, GetColorProfile},

		{"UpdateColorprofile", http.MethodPut, profileIDPath, UpdateColorProfile},

		{"DeleteColorprofile", http.MethodDelete, profileIDPath, DeleteColorProfile},
	}
}

// GetAllColorProfiles get all color profiles
func GetAllColorProfiles(w http.ResponseWriter, r *http.Request) {
	var profiles, err = database.GetAllColorProfiles()
	if err != nil {
		api.HandleError(&w, http.StatusNotFound, err.Error())
		return
	}

	api.HandleJSON(&w, http.StatusOK, profiles)
}

// GetColorProfile get a specific color profile
func GetColorProfile(w http.ResponseWriter, r *http.Request) {
	// Get model if exist
	var profile, err = database.GetColorProfile(api.GetParam(r, "id"))
	if err != nil {
		api.HandleError(&w, http.StatusNotFound, profileNotFoundMsg)
		return
	}
	api.HandleJSON(&w, http.StatusOK, profile)
}

// CreateColorProfile create a color profile
func CreateColorProfile(w http.ResponseWriter, r *http.Request) {
	// Validate input
	var input models.ColorProfile
	if err := api.BindJSON(r, &input); err != nil {
		api.HandleError(&w, http.StatusBadRequest, err.Error())
		return
	}

	// generate an id
	input.ID = utils.GenerateID()
	if err := database.DB.Create(&input).Error; err != nil {
		log.Printf("Error: %s", err)
		api.HandleError(&w, http.StatusBadRequest, err.Error())
		return
	}
	w.Header().Add("Location", fmt.Sprintf("%s/%d", r.RequestURI, input.ID))
	api.HandleJSON(&w, http.StatusCreated, input)
}

// UpdateColorProfile update a color profile
func UpdateColorProfile(w http.ResponseWriter, r *http.Request) {
	// Get model if exist
	var profile, err = database.GetColorProfile(api.GetParam(r, "id"))
	if err != nil {
		api.HandleError(&w, http.StatusNotFound, profileNotFoundMsg)
		return
	}

	// Validate input
	var input models.ColorProfile
	if err := api.BindJSON(r, &input); err != nil {
		api.HandleError(&w, http.StatusBadRequest, err.Error())
		return
	}
	if err := database.UpdateProfile(profile, input); err != nil {
		api.HandleError(&w, http.StatusBadRequest, err.Error())
		return
	}

	go messaging.PublishProfileSaveEvent(null.NewInt(input.ID, true), input)

	api.HandleJSON(&w, http.StatusOK, profile)
}

// DeleteColorProfile delete a color profile
func DeleteColorProfile(w http.ResponseWriter, r *http.Request) {
	// Get model if exist
	var profile, err = database.GetColorProfile(api.GetParam(r, "id"))
	if err != nil {
		api.HandleError(&w, http.StatusNotFound, profileNotFoundMsg)
		return
	}
	if err := database.DB.Delete(&profile).Error; err != nil {
		api.HandleError(&w, http.StatusBadRequest, err.Error())
		return
	}

	go messaging.PublishProfileDeleteEvent(null.NewInt(profile.ID, true))
	api.HandleJSON(&w, http.StatusNoContent, nil)
}
