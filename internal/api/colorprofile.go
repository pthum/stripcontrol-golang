package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pthum/null"
	"github.com/pthum/stripcontrol-golang/internal/database"
	"github.com/pthum/stripcontrol-golang/internal/messaging"
	"github.com/pthum/stripcontrol-golang/internal/model"
	"github.com/pthum/stripcontrol-golang/internal/utils"
)

const (
	profileNotFoundMsg = "Profile not found!"
	profilePath        = "/api/colorprofile"
	profileIDPath      = profilePath + "/{id}"
)

type CPHandler interface {
	GetAllColorProfiles(w http.ResponseWriter, r *http.Request)
	GetColorProfile(w http.ResponseWriter, r *http.Request)
	UpdateLedStrip(w http.ResponseWriter, r *http.Request)
	CreateColorProfile(w http.ResponseWriter, r *http.Request)
}
type CPHandlerImpl struct {
	dbh database.DBHandler
	mh  messaging.EventHandler
}

func colorProfileRoutes(db database.DBHandler, mh messaging.EventHandler) []Route {
	h := CPHandlerImpl{
		dbh: db,
		mh:  mh,
	}
	return []Route{

		{"GetColorprofiles", http.MethodGet, profilePath, h.GetAllColorProfiles},

		{"CreateColorprofile", http.MethodPost, profilePath, h.CreateColorProfile},

		{"GetColorprofile", http.MethodGet, profileIDPath, h.GetColorProfile},

		{"UpdateColorprofile", http.MethodPut, profileIDPath, h.UpdateColorProfile},

		{"DeleteColorprofile", http.MethodDelete, profileIDPath, h.DeleteColorProfile},
	}
}

// GetAllColorProfiles get all color profiles
func (h *CPHandlerImpl) GetAllColorProfiles(w http.ResponseWriter, r *http.Request) {
	var profiles []model.ColorProfile
	var err = h.dbh.GetAll(&profiles)
	if err != nil {
		HandleError(&w, http.StatusNotFound, err.Error())
		return
	}

	HandleJSON(&w, http.StatusOK, profiles)
}

// GetColorProfile get a specific color profile
func (h *CPHandlerImpl) GetColorProfile(w http.ResponseWriter, r *http.Request) {
	// Get model if exist
	var profile model.ColorProfile

	if h.dbh.Get(GetParam(r, "id"), &profile) != nil {
		HandleError(&w, http.StatusNotFound, profileNotFoundMsg)
		return
	}
	HandleJSON(&w, http.StatusOK, profile)
}

// CreateColorProfile create a color profile
func (h *CPHandlerImpl) CreateColorProfile(w http.ResponseWriter, r *http.Request) {
	// Validate input
	var input model.ColorProfile
	if err := BindJSON(r, &input); err != nil {
		HandleError(&w, http.StatusBadRequest, err.Error())
		return
	}

	// generate an id
	input.ID = utils.GenerateID()

	if err := h.dbh.Create(&input); err != nil {
		log.Printf("Error: %s", err)
		HandleError(&w, http.StatusBadRequest, err.Error())
		return
	}
	w.Header().Add("Location", fmt.Sprintf("%s/%d", r.RequestURI, input.ID))
	HandleJSON(&w, http.StatusCreated, input)
}

// UpdateColorProfile update a color profile
func (h *CPHandlerImpl) UpdateColorProfile(w http.ResponseWriter, r *http.Request) {
	// Get model if exist
	var profile model.ColorProfile
	if h.dbh.Get(GetParam(r, "id"), &profile) != nil {
		HandleError(&w, http.StatusNotFound, profileNotFoundMsg)
		return
	}

	// Validate input
	var input model.ColorProfile
	if err := BindJSON(r, &input); err != nil {
		HandleError(&w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.dbh.Update(profile, input); err != nil {
		HandleError(&w, http.StatusBadRequest, err.Error())
		return
	}

	go h.mh.PublishProfileSaveEvent(null.NewInt(input.ID, true), input)

	HandleJSON(&w, http.StatusOK, profile)
}

// DeleteColorProfile delete a color profile
func (h *CPHandlerImpl) DeleteColorProfile(w http.ResponseWriter, r *http.Request) {
	// Get model if exist
	var profile model.ColorProfile
	var err = h.dbh.Get(GetParam(r, "id"), &profile)
	if err != nil {
		HandleError(&w, http.StatusNotFound, profileNotFoundMsg)
		return
	}
	if err := h.dbh.Delete(&profile); err != nil {
		HandleError(&w, http.StatusBadRequest, err.Error())
		return
	}

	go h.mh.PublishProfileDeleteEvent(null.NewInt(profile.ID, true))
	HandleJSON(&w, http.StatusNoContent, nil)
}
