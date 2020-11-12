package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
func GetAllLedStrips(c *gin.Context) {
	var strips, err = database.GetAllLedStrips()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, strips)
}

// GetLedStrip get a single led strip
func GetLedStrip(c *gin.Context) {
	// Get model if exist
	var strip, err = database.GetLedStrip(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": stripNotFoundMsg})
		return
	}
	c.JSON(http.StatusOK, strip)
}

// CreateLedStrip create an LED strip
func CreateLedStrip(c *gin.Context) {
	// Validate input
	var input models.LedStrip
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// generate an id
	input.ID = utils.GenerateID()
	log.Printf("Generated ID %d", input.ID)
	if err := database.DB.Create(&input).Error; err != nil {
		log.Printf("Error: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	go messaging.PublishStripSaveEvent(null.NewInt(0, false), input)
	log.Printf("ID after save %d", input.ID)
	c.Header("Location", fmt.Sprintf("%s/%d", c.Request.URL.String(), input.ID))
	c.JSON(http.StatusCreated, input)
}

// UpdateLedStrip update an LED strip
func UpdateLedStrip(c *gin.Context) {
	// Get model if exist
	var strip, err = database.GetLedStrip(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": stripNotFoundMsg})
		return
	}

	// Validate input
	var input models.LedStrip
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// profile shouldn't be updated through this endpoint
	input.ProfileID = strip.ProfileID

	if err := database.DB.Save(&input).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	go messaging.PublishStripSaveEvent(null.NewInt(input.ID, true), input)

	c.JSON(http.StatusNoContent, nil)
}

// DeleteLedStrip delete an LED strip
func DeleteLedStrip(c *gin.Context) {
	// Get model if exist
	var strip, err = database.GetLedStrip(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": stripNotFoundMsg})
		return
	}

	if err := database.DB.Delete(&strip).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	go messaging.PublishStripDeleteEvent(null.NewInt(strip.ID, true))
	c.JSON(http.StatusNoContent, nil)
}

// UpdateProfileForStrip update which profile is referenced to the strip
func UpdateProfileForStrip(c *gin.Context) {
	// Get model if exist
	var strip, err1 = database.GetLedStrip(c.Param("id"))
	if err1 != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err1.Error()})
		return
	}

	// Validate input
	var input models.ColorProfile
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var profile, err2 = database.GetColorProfile(strconv.FormatInt(input.ID, 10))
	if err2 != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err2.Error()})
		return
	}
	strip.ProfileID = null.NewInt(profile.ID, true)
	database.DB.Save(strip)

	go messaging.PublishStripSaveEvent(null.NewInt(strip.ID, true), strip)

	c.JSON(http.StatusOK, profile)
}

// GetProfileForStrip get the current profile of a strip
func GetProfileForStrip(c *gin.Context) {
	// Get model if exist
	var strip, err1 = database.GetLedStrip(c.Param("id"))
	if err1 != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err1.Error()})
		return
	}
	if !strip.ProfileID.Valid {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}
	var profile, err2 = database.GetColorProfile(strconv.FormatInt(strip.ProfileID.Int64, 10))
	if err2 != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err2.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// RemoveProfileForStrip remove the current referenced profile
func RemoveProfileForStrip(c *gin.Context) {
	// Get model if exist
	var strip, err1 = database.GetLedStrip(c.Param("id"))
	if err1 != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err1.Error()})
		return
	}
	strip.ProfileID.Valid = false
	database.DB.Save(strip)

	go messaging.PublishStripSaveEvent(null.NewInt(strip.ID, true), strip)

	c.JSON(http.StatusNoContent, nil)
}
