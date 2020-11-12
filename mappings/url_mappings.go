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
		// mappings for /api/ledstrip
		stripAPI := api.Group("/ledstrip")
		{
			stripAPI.POST("", controllers.CreateLedStrip)
			stripAPI.GET("", controllers.GetAllLedStrips)

			// mappings for /api/ledstrip/:id
			stripIDAPI := stripAPI.Group("/:id")
			{
				stripIDAPI.GET("", controllers.GetLedStrip)
				stripIDAPI.PUT("", controllers.UpdateLedStrip)
				stripIDAPI.DELETE("", controllers.DeleteLedStrip)

				// mappings for /api/ledstrip/:id/profile
				stripIDProfileAPI := stripIDAPI.Group("/profile")
				{
					stripIDProfileAPI.GET("", controllers.GetProfileForStrip)
					stripIDProfileAPI.PUT("", controllers.UpdateProfileForStrip)
					stripIDProfileAPI.DELETE("", controllers.RemoveProfileForStrip)
				}
			}
		}
		// mappings for /api/colorprofile
		profileAPI := api.Group("/colorprofile")
		{
			profileAPI.POST("", controllers.CreateColorProfile)
			profileAPI.GET("", controllers.GetAllColorProfiles)

			// mappings for /api/colorprofile/:id
			profileIDAPI := profileAPI.Group("/:id")
			{
				profileIDAPI.GET("", controllers.GetColorProfile)
				profileIDAPI.PUT("", controllers.UpdateColorProfile)
				profileIDAPI.DELETE("", controllers.DeleteColorProfile)
			}
		}
	}
	// static content
	Router.StaticFile("/", "./static/index.html")
	Router.StaticFile("/favicon.ico", "./static/favicon.ico")
	Router.Static("/static", "./static/static")

}
