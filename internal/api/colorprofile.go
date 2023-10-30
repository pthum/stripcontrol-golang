package api

import (
	"log"
	"net/http"

	"github.com/pthum/null"
	"github.com/pthum/stripcontrol-golang/internal/database"
	"github.com/pthum/stripcontrol-golang/internal/messaging"
	"github.com/pthum/stripcontrol-golang/internal/model"
	"github.com/samber/do"
)

const (
	profileNotFoundMsg = "Profile not found!"
	profilePath        = "/api/colorprofile"
	profileIDPath      = profilePath + "/{id}"
)

type CPHandler interface {
	GetAllColorProfiles(w http.ResponseWriter, r *http.Request)
	GetColorProfile(w http.ResponseWriter, r *http.Request)
	UpdateColorProfile(w http.ResponseWriter, r *http.Request)
	CreateColorProfile(w http.ResponseWriter, r *http.Request)
}
type CPHandlerImpl struct {
	dbh database.DBHandler[model.ColorProfile]
	mh  messaging.EventHandler
}

func NewCPHandler(i *do.Injector) (CPHandler, error) {
	db := do.MustInvoke[database.DBHandler[model.ColorProfile]](i)
	mh := do.MustInvoke[messaging.EventHandler](i)
	return &CPHandlerImpl{
		dbh: db,
		mh:  mh,
	}, nil
}

func (h *CPHandlerImpl) colorProfileRoutes() []Route {
	return []Route{
		{http.MethodGet, profilePath, h.GetAllColorProfiles},
		{http.MethodPost, profilePath, h.CreateColorProfile},
		{http.MethodGet, profileIDPath, h.GetColorProfile},
		{http.MethodPut, profileIDPath, h.UpdateColorProfile},
		{http.MethodDelete, profileIDPath, h.DeleteColorProfile},
	}
}

// GetAllColorProfiles get all color profiles
func (h *CPHandlerImpl) GetAllColorProfiles(w http.ResponseWriter, r *http.Request) {
	profiles, err := h.dbh.GetAll()
	if err != nil {
		handleError(&w, http.StatusNotFound, err.Error())
		return
	}

	handleJSON(&w, http.StatusOK, profiles)
}

// GetColorProfile get a specific color profile
func (h *CPHandlerImpl) GetColorProfile(w http.ResponseWriter, r *http.Request) {
	// Get model if exist
	profile, err := h.dbh.Get(getParam(r, "id"))
	if err != nil {
		handleError(&w, http.StatusNotFound, profileNotFoundMsg)
		return
	}

	handleJSON(&w, http.StatusOK, profile)
}

// CreateColorProfile create a color profile
func (h *CPHandlerImpl) CreateColorProfile(w http.ResponseWriter, r *http.Request) {
	// Validate input
	var input model.ColorProfile
	if err := bindJSON(r, &input); err != nil {
		handleError(&w, http.StatusBadRequest, err.Error())
		return
	}

	// generate an id
	input.GenerateID()

	if err := h.dbh.Create(&input); err != nil {
		log.Printf("Error: %s", err)
		handleError(&w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithCreated(r, w, &input)
}

// UpdateColorProfile update a color profile
func (h *CPHandlerImpl) UpdateColorProfile(w http.ResponseWriter, r *http.Request) {
	// Get model if exist
	profile, err := h.dbh.Get(getParam(r, "id"))
	if err != nil {
		handleError(&w, http.StatusNotFound, profileNotFoundMsg)
		return
	}

	// Validate input
	var input model.ColorProfile
	if err := bindJSON(r, &input); err != nil {
		handleError(&w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.dbh.Update(*profile, input); err != nil {
		handleError(&w, http.StatusBadRequest, err.Error())
		return
	}

	var event = model.NewProfileEvent(null.NewInt(input.ID, true), model.Save).With(input)
	go h.mh.PublishProfileEvent(event)

	handleJSON(&w, http.StatusOK, profile)
}

// DeleteColorProfile delete a color profile
func (h *CPHandlerImpl) DeleteColorProfile(w http.ResponseWriter, r *http.Request) {
	// Get model if exist
	profile, err := h.dbh.Get(getParam(r, "id"))
	if err != nil {
		handleError(&w, http.StatusNotFound, profileNotFoundMsg)
		return
	}

	if err := h.dbh.Delete(profile); err != nil {
		handleError(&w, http.StatusBadRequest, err.Error())
		return
	}

	var event = model.NewProfileEvent(null.NewInt(profile.ID, true), model.Delete)
	go h.mh.PublishProfileEvent(event)

	handleJSON(&w, http.StatusNoContent, nil)
}
