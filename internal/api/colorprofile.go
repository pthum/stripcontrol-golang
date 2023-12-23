package api

import (
	"log"
	"net/http"

	"github.com/pthum/stripcontrol-golang/internal/model"
	"github.com/pthum/stripcontrol-golang/internal/service"
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
	CreateColorProfile(w http.ResponseWriter, r *http.Request)
	UpdateColorProfile(w http.ResponseWriter, r *http.Request)
	DeleteColorProfile(w http.ResponseWriter, r *http.Request)
}
type cpHandlerImpl struct {
	cps service.CPService
}

func NewCPHandler(i *do.Injector) (CPHandler, error) {
	cps := do.MustInvoke[service.CPService](i)
	return &cpHandlerImpl{
		cps: cps,
	}, nil
}

func (h *cpHandlerImpl) colorProfileRoutes() []Route {
	return []Route{
		{http.MethodGet, profilePath, h.GetAllColorProfiles},
		{http.MethodPost, profilePath, h.CreateColorProfile},
		{http.MethodGet, profileIDPath, h.GetColorProfile},
		{http.MethodPut, profileIDPath, h.UpdateColorProfile},
		{http.MethodDelete, profileIDPath, h.DeleteColorProfile},
	}
}

// GetAllColorProfiles get all color profiles
func (h *cpHandlerImpl) GetAllColorProfiles(w http.ResponseWriter, r *http.Request) {
	profiles, err := h.cps.GetAll()
	if err != nil {
		handleError(&w, http.StatusNotFound, err.Error())
		return
	}

	handleJSON(&w, http.StatusOK, profiles)
}

// GetColorProfile get a specific color profile
func (h *cpHandlerImpl) GetColorProfile(w http.ResponseWriter, r *http.Request) {
	// Get model if exist
	profile, err := h.cps.GetColorProfile(getParam(r, "id"))
	if err != nil {
		handleError(&w, http.StatusNotFound, profileNotFoundMsg)
		return
	}

	handleJSON(&w, http.StatusOK, profile)
}

// CreateColorProfile create a color profile
func (h *cpHandlerImpl) CreateColorProfile(w http.ResponseWriter, r *http.Request) {
	// Validate input
	var input model.ColorProfile
	if err := bindJSON(r, &input); err != nil {
		handleError(&w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.cps.CreateColorProfile(&input); err != nil {
		log.Printf("Error: %s", err)
		handleError(&w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithCreated(r, w, &input)
}

// UpdateColorProfile update a color profile
func (h *cpHandlerImpl) UpdateColorProfile(w http.ResponseWriter, r *http.Request) {
	// Validate input
	var input model.ColorProfile
	if err := bindJSON(r, &input); err != nil {
		handleErr(&w, model.NewAppErr(http.StatusBadRequest, err))
		return
	}

	if err := h.cps.UpdateColorProfile(getParam(r, "id"), input); err != nil {
		handleErr(&w, err)
		return
	}

	handleJSON(&w, http.StatusOK, input) //FIXME input to pointer
}

// DeleteColorProfile delete a color profile
func (h *cpHandlerImpl) DeleteColorProfile(w http.ResponseWriter, r *http.Request) {
	if err := h.cps.DeleteColorProfile(getParam(r, "id")); err != nil {
		handleErr(&w, err)
		return
	}

	handleJSON(&w, http.StatusNoContent, nil)
}
