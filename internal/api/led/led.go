package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/pthum/null"
	api "github.com/pthum/stripcontrol-golang/internal/api/common"
	"github.com/pthum/stripcontrol-golang/internal/database"
	"github.com/pthum/stripcontrol-golang/internal/messaging"
	"github.com/pthum/stripcontrol-golang/internal/model"
	"github.com/pthum/stripcontrol-golang/internal/utils"
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
	dbh database.DBHandler
	mh  messaging.EventHandler
}

func LEDRoutes(db database.DBHandler, mh messaging.EventHandler) []api.Route {
	lh := LEDHandlerImpl{
		dbh: db,
		mh:  mh,
	}
	return []api.Route{

		api.Route{"GetLedstrips", http.MethodGet, ledstripPath, lh.GetAllLedStrips},

		api.Route{"CreateLedstrip", http.MethodPost, ledstripPath, lh.CreateLedStrip},

		api.Route{"GetLedstrip", http.MethodGet, ledstripIDPath, lh.GetLedStrip},

		api.Route{"UpdateLedstrip", http.MethodPut, ledstripIDPath, lh.UpdateLedStrip},

		api.Route{"DeleteLedstripId", http.MethodDelete, ledstripIDPath, lh.DeleteLedStrip},

		api.Route{"GetLedstripReferencedProfile", http.MethodGet, ledstripIDProfilePath, lh.GetProfileForStrip},

		api.Route{"UpdateLedstripReferencedProfile", http.MethodPut, ledstripIDProfilePath, lh.UpdateProfileForStrip},

		api.Route{"DeleteLedstripReferencedProfile", http.MethodDelete, ledstripIDProfilePath, lh.RemoveProfileForStrip},
	}
}

// GetAllLedStrips get all existing led strips
func (lh *LEDHandlerImpl) GetAllLedStrips(w http.ResponseWriter, r *http.Request) {
	var strips []model.LedStrip
	err := lh.dbh.GetAll(&strips)
	// var strips, err = database.GetAllLedStrips()
	if err != nil {
		api.HandleError(&w, http.StatusNotFound, err.Error())
		return
	}

	api.HandleJSON(&w, http.StatusOK, strips)
}

// GetLedStrip get a single led strip
func (lh *LEDHandlerImpl) GetLedStrip(w http.ResponseWriter, r *http.Request) {
	// Get model if exist
	var strip model.LedStrip
	var err = lh.dbh.Get(api.GetParam(r, "id"), &strip)
	// var strip, err = database.GetLedStrip(api.GetParam(r, "id"))
	if err != nil {
		api.HandleError(&w, http.StatusNotFound, stripNotFoundMsg)
		return
	}
	api.HandleJSON(&w, http.StatusOK, strip)
}

// CreateLedStrip create an LED strip
func (lh *LEDHandlerImpl) CreateLedStrip(w http.ResponseWriter, r *http.Request) {
	// Validate input
	var input model.LedStrip
	if err := api.BindJSON(r, &input); err != nil {
		api.HandleError(&w, http.StatusBadRequest, err.Error())
		return
	}
	// generate an id
	input.ID = utils.GenerateID()
	log.Printf("Generated ID %d", input.ID)

	if err := lh.dbh.Create(&input); err != nil {
		log.Printf("Error: %s", err)
		api.HandleError(&w, http.StatusBadRequest, err.Error())
		return
	}

	go lh.mh.PublishStripSaveEvent(null.NewInt(0, false), input)
	log.Printf("ID after save %d", input.ID)
	w.Header().Add("Location", fmt.Sprintf("%s/%d", r.RequestURI, input.ID))
	api.HandleJSON(&w, http.StatusCreated, input)
}

// UpdateLedStrip update an LED strip
func (lh *LEDHandlerImpl) UpdateLedStrip(w http.ResponseWriter, r *http.Request) {
	// Get model if exist
	var strip model.LedStrip
	var err = lh.dbh.Get(api.GetParam(r, "id"), &strip)
	if err != nil {
		api.HandleError(&w, http.StatusNotFound, stripNotFoundMsg)
		return
	}

	// Validate input
	var input model.LedStrip
	if err := api.BindJSON(r, &input); err != nil {
		api.HandleError(&w, http.StatusBadRequest, err.Error())
		return
	}
	// run db update for strip async, as precondition checks were successful,
	// chances are good that the update will be successful
	// if the db update fails, the hardware won't change (no message sent)
	// in that case the UI would not reflect the current state,
	// which we accept for now
	go lh.updateAndHandle(strip, input)

	api.HandleJSON(&w, http.StatusNoContent, nil)
}

// DeleteLedStrip delete an LED strip
func (lh *LEDHandlerImpl) DeleteLedStrip(w http.ResponseWriter, r *http.Request) {
	// Get model if exist
	var strip model.LedStrip
	var err = lh.dbh.Get(api.GetParam(r, "id"), &strip)
	if err != nil {
		api.HandleError(&w, http.StatusNotFound, stripNotFoundMsg)
		return
	}

	if err := lh.dbh.Delete(&strip); err != nil {
		api.HandleError(&w, http.StatusBadRequest, err.Error())
		return
	}
	go lh.mh.PublishStripDeleteEvent(null.NewInt(strip.ID, true))
	api.HandleJSON(&w, http.StatusNoContent, nil)
}

// UpdateProfileForStrip update which profile is referenced to the strip
func (lh *LEDHandlerImpl) UpdateProfileForStrip(w http.ResponseWriter, r *http.Request) {
	// Get model if exist
	var strip model.LedStrip
	var err1 = lh.dbh.Get(api.GetParam(r, "id"), &strip)
	if err1 != nil {
		api.HandleError(&w, http.StatusBadRequest, err1.Error())
		return
	}

	// Validate input
	var input model.ColorProfile
	if err := api.BindJSON(r, &input); err != nil {
		api.HandleError(&w, http.StatusBadRequest, err.Error())
		return
	}
	var profile model.ColorProfile
	var err2 = lh.dbh.Get(strconv.FormatInt(input.ID, 10), &profile)
	if err2 != nil {
		api.HandleError(&w, http.StatusBadRequest, err2.Error())
		return
	}
	strip.ProfileID = null.NewInt(profile.ID, true)

	lh.dbh.Save(strip)

	go lh.mh.PublishStripSaveEvent(null.NewInt(strip.ID, true), strip)

	api.HandleJSON(&w, http.StatusOK, profile)
}

// GetProfileForStrip get the current profile of a strip
func (lh *LEDHandlerImpl) GetProfileForStrip(w http.ResponseWriter, r *http.Request) {
	// Get model if exist
	var strip model.LedStrip
	var err1 = lh.dbh.Get(api.GetParam(r, "id"), &strip)
	if err1 != nil {
		api.HandleError(&w, http.StatusBadRequest, err1.Error())
		return
	}
	if !strip.ProfileID.Valid {
		api.HandleError(&w, http.StatusBadRequest, "Profile not found!")
		return
	}
	var profile model.ColorProfile
	var err2 = lh.dbh.Get(strconv.FormatInt(strip.ProfileID.Int64, 10), &profile)
	if err2 != nil {
		api.HandleError(&w, http.StatusBadRequest, err2.Error())
		return
	}

	api.HandleJSON(&w, http.StatusOK, profile)
}

// RemoveProfileForStrip remove the current referenced profile
func (lh *LEDHandlerImpl) RemoveProfileForStrip(w http.ResponseWriter, r *http.Request) {
	// Get model if exist
	var strip model.LedStrip
	var err1 = lh.dbh.Get(api.GetParam(r, "id"), &strip)
	if err1 != nil {
		api.HandleError(&w, http.StatusBadRequest, err1.Error())
		return
	}
	strip.ProfileID.Valid = false
	lh.dbh.Save(strip)

	go lh.mh.PublishStripSaveEvent(null.NewInt(strip.ID, true), strip)

	api.HandleJSON(&w, http.StatusNoContent, nil)
}

func (lh *LEDHandlerImpl) updateAndHandle(strip model.LedStrip, input model.LedStrip) {
	// profile shouldn't be updated through this endpoint
	input.ProfileID = strip.ProfileID

	if err := lh.dbh.Update(&strip, &input); err != nil {
		log.Printf("error: %s", err.Error())
		return
	}
	lh.mh.PublishStripSaveEvent(null.NewInt(input.ID, true), input)
}
