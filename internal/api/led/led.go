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
	"github.com/pthum/stripcontrol-golang/internal/models"
	"github.com/pthum/stripcontrol-golang/internal/utils"
)

const (
	stripNotFoundMsg      = "LEDStrip not found!"
	ledstripPath          = "/api/ledstrip"
	ledstripIDPath        = ledstripPath + "/{id}"
	ledstripIDProfilePath = ledstripIDPath + "/profile"
)

func LEDRoutes() []api.Route {
	return []api.Route{

		api.Route{"GetLedstrips", http.MethodGet, ledstripPath, GetAllLedStrips},

		api.Route{"CreateLedstrip", http.MethodPost, ledstripPath, CreateLedStrip},

		api.Route{"GetLedstrip", http.MethodGet, ledstripIDPath, GetLedStrip},

		api.Route{"UpdateLedstrip", http.MethodPut, ledstripIDPath, UpdateLedStrip},

		api.Route{"DeleteLedstripId", http.MethodDelete, ledstripIDPath, DeleteLedStrip},

		api.Route{"GetLedstripReferencedProfile", http.MethodGet, ledstripIDProfilePath, GetProfileForStrip},

		api.Route{"UpdateLedstripReferencedProfile", http.MethodPut, ledstripIDProfilePath, UpdateProfileForStrip},

		api.Route{"DeleteLedstripReferencedProfile", http.MethodDelete, ledstripIDProfilePath, RemoveProfileForStrip},
	}
}

// GetAllLedStrips get all existing led strips
func GetAllLedStrips(w http.ResponseWriter, r *http.Request) {
	var strips, err = database.GetAllLedStrips()
	if err != nil {
		api.HandleError(&w, http.StatusNotFound, err.Error())
		return
	}

	api.HandleJSON(&w, http.StatusOK, strips)
}

// GetLedStrip get a single led strip
func GetLedStrip(w http.ResponseWriter, r *http.Request) {
	// Get model if exist
	var strip, err = database.GetLedStrip(api.GetParam(r, "id"))
	if err != nil {
		api.HandleError(&w, http.StatusNotFound, stripNotFoundMsg)
		return
	}
	api.HandleJSON(&w, http.StatusOK, strip)
}

// CreateLedStrip create an LED strip
func CreateLedStrip(w http.ResponseWriter, r *http.Request) {
	// Validate input
	var input models.LedStrip
	if err := api.BindJSON(r, &input); err != nil {
		api.HandleError(&w, http.StatusBadRequest, err.Error())
		return
	}
	// generate an id
	input.ID = utils.GenerateID()
	log.Printf("Generated ID %d", input.ID)
	if err := database.DB.Create(&input).Error; err != nil {
		log.Printf("Error: %s", err)
		api.HandleError(&w, http.StatusBadRequest, err.Error())
		return
	}

	go messaging.PublishStripSaveEvent(null.NewInt(0, false), input)
	log.Printf("ID after save %d", input.ID)
	w.Header().Add("Location", fmt.Sprintf("%s/%d", r.RequestURI, input.ID))
	api.HandleJSON(&w, http.StatusCreated, input)
}

// UpdateLedStrip update an LED strip
func UpdateLedStrip(w http.ResponseWriter, r *http.Request) {
	// Get model if exist
	var strip, err = database.GetLedStrip(api.GetParam(r, "id"))
	if err != nil {
		api.HandleError(&w, http.StatusNotFound, stripNotFoundMsg)
		return
	}

	// Validate input
	var input models.LedStrip
	if err := api.BindJSON(r, &input); err != nil {
		api.HandleError(&w, http.StatusBadRequest, err.Error())
		return
	}
	// run db update for strip async, as precondition checks were successful,
	// chances are good that the update will be successful
	// if the db update fails, the hardware won't change (no message sent)
	// in that case the UI would not reflect the current state,
	// which we accept for now
	go updateAndHandle(strip, input)

	api.HandleJSON(&w, http.StatusNoContent, nil)
}

// DeleteLedStrip delete an LED strip
func DeleteLedStrip(w http.ResponseWriter, r *http.Request) {
	// Get model if exist
	var strip, err = database.GetLedStrip(api.GetParam(r, "id"))
	if err != nil {
		api.HandleError(&w, http.StatusNotFound, stripNotFoundMsg)
		return
	}

	if err := database.DB.Delete(&strip).Error; err != nil {
		api.HandleError(&w, http.StatusBadRequest, err.Error())
		return
	}
	go messaging.PublishStripDeleteEvent(null.NewInt(strip.ID, true))
	api.HandleJSON(&w, http.StatusNoContent, nil)
}

// UpdateProfileForStrip update which profile is referenced to the strip
func UpdateProfileForStrip(w http.ResponseWriter, r *http.Request) {
	// Get model if exist
	var strip, err1 = database.GetLedStrip(api.GetParam(r, "id"))
	if err1 != nil {
		api.HandleError(&w, http.StatusBadRequest, err1.Error())
		return
	}

	// Validate input
	var input models.ColorProfile
	if err := api.BindJSON(r, &input); err != nil {
		api.HandleError(&w, http.StatusBadRequest, err.Error())
		return
	}
	var profile, err2 = database.GetColorProfile(strconv.FormatInt(input.ID, 10))
	if err2 != nil {
		api.HandleError(&w, http.StatusBadRequest, err2.Error())
		return
	}
	strip.ProfileID = null.NewInt(profile.ID, true)
	database.DB.Save(strip)

	go messaging.PublishStripSaveEvent(null.NewInt(strip.ID, true), strip)

	api.HandleJSON(&w, http.StatusOK, profile)
}

// GetProfileForStrip get the current profile of a strip
func GetProfileForStrip(w http.ResponseWriter, r *http.Request) {
	// Get model if exist
	var strip, err1 = database.GetLedStrip(api.GetParam(r, "id"))
	if err1 != nil {
		api.HandleError(&w, http.StatusBadRequest, err1.Error())
		return
	}
	if !strip.ProfileID.Valid {
		api.HandleError(&w, http.StatusBadRequest, "Profile not found!")
		return
	}
	var profile, err2 = database.GetColorProfile(strconv.FormatInt(strip.ProfileID.Int64, 10))
	if err2 != nil {
		api.HandleError(&w, http.StatusBadRequest, err2.Error())
		return
	}

	api.HandleJSON(&w, http.StatusOK, profile)
}

// RemoveProfileForStrip remove the current referenced profile
func RemoveProfileForStrip(w http.ResponseWriter, r *http.Request) {
	// Get model if exist
	var strip, err1 = database.GetLedStrip(api.GetParam(r, "id"))
	if err1 != nil {
		api.HandleError(&w, http.StatusBadRequest, err1.Error())
		return
	}
	strip.ProfileID.Valid = false
	database.DB.Save(strip)

	go messaging.PublishStripSaveEvent(null.NewInt(strip.ID, true), strip)

	api.HandleJSON(&w, http.StatusNoContent, nil)
}

func updateAndHandle(strip models.LedStrip, input models.LedStrip) {
	// profile shouldn't be updated through this endpoint
	input.ProfileID = strip.ProfileID

	if err := database.UpdateStrip(strip, input); err != nil {
		log.Printf("error: %s", err.Error())
		return
	}
	messaging.PublishStripSaveEvent(null.NewInt(input.ID, true), input)
}
