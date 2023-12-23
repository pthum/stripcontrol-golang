package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/pthum/null"
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

type LEDHandlerImpl struct {
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
	return &LEDHandlerImpl{
		dbh:   lsdb,
		cpDbh: cpdb,
		mh:    mh,
		lsvc:  lsvc,
	}, nil
}

func (lh *LEDHandlerImpl) ledRoutes() []Route {
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
func (lh *LEDHandlerImpl) GetAllLedStrips(w http.ResponseWriter, r *http.Request) {
	strips, err := lh.lsvc.GetAll()
	if err != nil {
		handleError(&w, http.StatusNotFound, err.Error())
		return
	}

	handleJSON(&w, http.StatusOK, strips)
}

// GetLedStrip get a single led strip
func (lh *LEDHandlerImpl) GetLedStrip(w http.ResponseWriter, r *http.Request) {
	// Get model if exist
	strip, err := lh.lsvc.GetLEDStrip(getParam(r, "id"))
	if err != nil {
		handleError(&w, http.StatusNotFound, stripNotFoundMsg)
		return
	}

	handleJSON(&w, http.StatusOK, strip)
}

// CreateLedStrip create an LED strip
func (lh *LEDHandlerImpl) CreateLedStrip(w http.ResponseWriter, r *http.Request) {
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
func (lh *LEDHandlerImpl) UpdateLedStrip(w http.ResponseWriter, r *http.Request) {
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
func (lh *LEDHandlerImpl) DeleteLedStrip(w http.ResponseWriter, r *http.Request) {
	if err := lh.lsvc.DeleteLEDStrip(getParam(r, "id")); err != nil {
		handleErr(&w, err)
		return
	}

	handleJSON(&w, http.StatusNoContent, nil)
}

// UpdateProfileForStrip update which profile is referenced to the strip
func (lh *LEDHandlerImpl) UpdateProfileForStrip(w http.ResponseWriter, r *http.Request) {
	// Validate input
	var input model.ColorProfile
	if err := bindJSON(r, &input); err != nil {
		handleError(&w, http.StatusBadRequest, err.Error())
		return
	}
	// Get model if exist
	strip, err := lh.dbh.Get(getParam(r, "id"))
	if err != nil {
		handleError(&w, http.StatusNotFound, err.Error())
		return
	}

	profile, err := lh.cpDbh.Get(input.GetStringID())
	if err != nil {
		handleError(&w, http.StatusNotFound, err.Error())
		return
	}

	strip.ProfileID = profile.GetNullID()

	if err := lh.dbh.Save(strip); err != nil {
		log.Printf("Error: %s", err)
		handleError(&w, http.StatusInternalServerError, err.Error())
		return
	}

	go lh.publishStripSaveEvent(strip.GetNullID(), *strip, profile)

	handleJSON(&w, http.StatusOK, profile)
}

// GetProfileForStrip get the current profile of a strip
func (lh *LEDHandlerImpl) GetProfileForStrip(w http.ResponseWriter, r *http.Request) {
	// Get model if exist
	strip, err := lh.dbh.Get(getParam(r, "id"))
	if err != nil {
		handleError(&w, http.StatusNotFound, err.Error())
		return
	}

	if !strip.ProfileID.Valid {
		handleError(&w, http.StatusNotFound, "Profile not found!")
		return
	}

	profile, err := lh.cpDbh.Get(strconv.FormatInt(strip.ProfileID.Int64, 10))
	if err != nil {
		handleError(&w, http.StatusNotFound, err.Error())
		return
	}

	handleJSON(&w, http.StatusOK, profile)
}

// RemoveProfileForStrip remove the current referenced profile
func (lh *LEDHandlerImpl) RemoveProfileForStrip(w http.ResponseWriter, r *http.Request) {
	// Get model if exist
	strip, err := lh.dbh.Get(getParam(r, "id"))
	if err != nil {
		handleError(&w, http.StatusNotFound, err.Error())
		return
	}

	strip.ProfileID.Valid = false

	if err := lh.dbh.Save(strip); err != nil {
		handleError(&w, http.StatusInternalServerError, err.Error())
		return
	}

	go lh.publishStripSaveEvent(strip.GetNullID(), *strip, nil)

	handleJSON(&w, http.StatusNoContent, nil)
}

func (lh *LEDHandlerImpl) publishStripSaveEvent(id null.Int, strip model.LedStrip, profile *model.ColorProfile) {
	var event = model.NewStripEvent(id, model.Save).With(&strip)

	if strip.ProfileID.Valid {
		if profile != nil {
			event.Strip.With(*profile)
		}
	}

	if err := lh.mh.PublishStripEvent(event); err != nil {
		log.Printf("error: %s", err.Error())
		return
	}
}
