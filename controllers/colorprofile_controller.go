package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pthum/stripcontrol-golang/database"
	"github.com/pthum/stripcontrol-golang/models"
	"github.com/pthum/stripcontrol-golang/utils"
)

const (
	profileNotFoundMsg = "Profile not found!"
)

// GetAllColorProfiles get all color profiles
func GetAllColorProfiles(c *gin.Context) {
	var profiles []models.ColorProfile
	database.DB.Find(&profiles)

	c.JSON(http.StatusOK, profiles)
}

// GetColorProfile get a specific color profile
func GetColorProfile(c *gin.Context) {
	// Get model if exist
	var profile, err = database.GetColorProfile(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": profileNotFoundMsg})
		return
	}
	c.JSON(http.StatusOK, profile)
}

// CreateColorProfile create a color profile
func CreateColorProfile(c *gin.Context) {
	// Validate input
	var input models.ColorProfile
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// generate an id
	input.ID = utils.GenerateID()
	if err := database.DB.Create(&input).Error; err != nil {
		log.Printf("Error: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Header("Location", fmt.Sprintf("%s/%d", c.Request.URL.String(), input.ID))
	c.JSON(http.StatusCreated, input)
}

// UpdateColorProfile update a color profile
func UpdateColorProfile(c *gin.Context) {
	// Get model if exist
	var profile, err = database.GetColorProfile(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": profileNotFoundMsg})
		return
	}

	// Validate input
	var input models.ColorProfile
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Model(&profile).Updates(input)

	c.JSON(http.StatusOK, profile)

}

// DeleteColorProfile delete a color profile
func DeleteColorProfile(c *gin.Context) {
	// Get model if exist
	var profile, err = database.GetColorProfile(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": profileNotFoundMsg})
		return
	}
	if err := database.DB.Delete(&profile).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
