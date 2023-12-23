package api

import (
	"log"
	"net/http"

	"github.com/pthum/stripcontrol-golang/internal/database"
	"github.com/pthum/stripcontrol-golang/internal/messaging"
	"github.com/pthum/stripcontrol-golang/internal/model"
	"github.com/pthum/stripcontrol-golang/internal/service"
	"github.com/samber/do"
)

const (
	stripNotFoundMsg      = "LEDStrip not found!"
	ledstripPath          = "/api/ledstrip"
	ledstripIDPath        = ledstripPath + "/{id}"
	ledstripIDProfilePath = ledstripIDPath + "/profile"
)

type LEDHandler interface {
	GetAllLedStrips(w http.ResponseWriter, r *http.Request)
	GetLedStrip(w http.ResponseWriter, r *http.Request)
	UpdateLedStrip(w http.ResponseWriter, r *http.Request)
}

type ledHandlerImpl struct {
	dbh   database.DBHandler[model.LedStrip]
	cpDbh database.DBHandler[model.ColorProfile]
	mh    messaging.EventHandler
	lsvc  service.LEDService
}

func NewLEDHandler(i *do.Injector) (LEDHandler, error) {
	lsdb := do.MustInvoke[database.DBHandler[model.LedStrip]](i)
	cpdb := do.MustInvoke[database.DBHandler[model.ColorProfile]](i)
	mh := do.MustInvoke[messaging.EventHandler](i)
	lsvc := do.MustInvoke[service.LEDService](i)
	return &ledHandlerImpl{
		dbh:   lsdb,
		cpDbh: cpdb,
		mh:    mh,
		lsvc:  lsvc,
	}, nil
}

func (lh *ledHandlerImpl) ledRoutes() []Route {
	return []Route{
		{http.MethodGet, ledstripPath, lh.GetAllLedStrips},
		{http.MethodPost, ledstripPath, lh.CreateLedStrip},
		{http.MethodGet, ledstripIDPath, lh.GetLedStrip},
		{http.MethodPut, ledstripIDPath, lh.UpdateLedStrip},
		{http.MethodDelete, ledstripIDPath, lh.DeleteLedStrip},
		{http.MethodGet, ledstripIDProfilePath, lh.GetProfileForStrip},
		{http.MethodPut, ledstripIDProfilePath, lh.UpdateProfileForStrip},
		{http.MethodDelete, ledstripIDProfilePath, lh.RemoveProfileForStrip},
	}
}

// GetAllLedStrips get all existing led strips
func (lh *ledHandlerImpl) GetAllLedStrips(w http.ResponseWriter, r *http.Request) {
	strips, err := lh.lsvc.GetAll()
	if err != nil {
		handleError(&w, http.StatusNotFound, err.Error())
		return
	}

	handleJSON(&w, http.StatusOK, strips)
}

// GetLedStrip get a single led strip
func (lh *ledHandlerImpl) GetLedStrip(w http.ResponseWriter, r *http.Request) {
	// Get model if exist
	strip, err := lh.lsvc.GetLEDStrip(getParam(r, "id"))
	if err != nil {
		handleError(&w, http.StatusNotFound, stripNotFoundMsg)
		return
	}

	handleJSON(&w, http.StatusOK, strip)
}

// CreateLedStrip create an LED strip
func (lh *ledHandlerImpl) CreateLedStrip(w http.ResponseWriter, r *http.Request) {
	// Validate input
	var input model.LedStrip
	if err := bindJSON(r, &input); err != nil {
		handleError(&w, http.StatusBadRequest, err.Error())
		return
	}

	if err := lh.lsvc.CreateLEDStrip(&input); err != nil {
		log.Printf("Error: %s", err)
		handleError(&w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithCreated(r, w, &input)
}

// UpdateLedStrip update an LED strip
func (lh *ledHandlerImpl) UpdateLedStrip(w http.ResponseWriter, r *http.Request) {
	// Validate input
	var input model.LedStrip
	if err := bindJSON(r, &input); err != nil {
		handleErr(&w, model.NewAppErr(http.StatusBadRequest, err))
		return
	}

	if err := lh.lsvc.UpdateLEDStrip(getParam(r, "id"), input); err != nil {
		handleErr(&w, err)
		return
	}

	handleJSON(&w, http.StatusOK, input)
}

// DeleteLedStrip delete an LED strip
func (lh *ledHandlerImpl) DeleteLedStrip(w http.ResponseWriter, r *http.Request) {
	if err := lh.lsvc.DeleteLEDStrip(getParam(r, "id")); err != nil {
		handleErr(&w, err)
		return
	}

	handleJSON(&w, http.StatusNoContent, nil)
}

// UpdateProfileForStrip update which profile is referenced to the strip
func (lh *ledHandlerImpl) UpdateProfileForStrip(w http.ResponseWriter, r *http.Request) {
	// Validate input
	var input model.ColorProfile
	if err := bindJSON(r, &input); err != nil {
		handleError(&w, http.StatusBadRequest, err.Error())
		return
	}
	profile, err := lh.lsvc.UpdateProfileForStrip(getParam(r, "id"), input)
	if err != nil {
		handleErr(&w, err)
	}

	handleJSON(&w, http.StatusOK, profile)
}

// GetProfileForStrip get the current profile of a strip
func (lh *ledHandlerImpl) GetProfileForStrip(w http.ResponseWriter, r *http.Request) {
	profile, err := lh.lsvc.GetProfileForStrip(getParam(r, "id"))
	if err != nil {
		handleErr(&w, err)
		return
	}

	handleJSON(&w, http.StatusOK, profile)
}

// RemoveProfileForStrip remove the current referenced profile
func (lh *ledHandlerImpl) RemoveProfileForStrip(w http.ResponseWriter, r *http.Request) {
	if err := lh.lsvc.RemoveProfileForStrip(getParam(r, "id")); err != nil {
		handleErr(&w, err)
		return
	}

	handleJSON(&w, http.StatusNoContent, nil)
}
