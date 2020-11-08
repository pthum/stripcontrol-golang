package mappings

import (
	"github.com/gin-gonic/gin"
	"github.com/pthum/stripcontrol-golang/controllers"
)

// Router router var
var Router *gin.Engine

// CreateURLMappings Create Url Mappings
func CreateURLMappings() {
	Router = gin.Default()

	//Router.Use(controllers.Cors())
	api := Router.Group("/api")
	{
		api.POST("/ledstrip", controllers.CreateLedStrip)
		api.GET("/ledstrip", controllers.GetAllLedStrips)
		api.GET("/ledstrip/:id", controllers.GetLedStrip)
		api.PUT("/ledstrip/:id", controllers.UpdateLedStrip)
		api.DELETE("/ledstrip/:id", controllers.DeleteLedStrip)

		api.GET("/ledstrip/:id/profile", controllers.GetProfileForStrip)
		api.PUT("/ledstrip/:id/profile", controllers.UpdateProfileForStrip)
		api.DELETE("/ledstrip/:id/profile", controllers.RemoveProfileForStrip)

		api.POST("/colorprofile", controllers.CreateColorProfile)
		api.GET("/colorprofile/:id", controllers.GetColorProfile)
		api.GET("/colorprofile", controllers.GetAllColorProfiles)
		api.PUT("/colorprofile/:id", controllers.UpdateColorProfile)
		api.DELETE("/colorprofile/:id", controllers.DeleteColorProfile)
	}
	Router.StaticFile("/", "./static/index.html")
	Router.StaticFile("/favicon.ico", "./static/favicon.ico")
	Router.Static("/static", "./static/static")

}
