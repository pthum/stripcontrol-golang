package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const (
	profilePath           = "/api/colorprofile"
	profileIDPath         = "/api/colorprofile/{id}"
	ledstripPath          = "/api/ledstrip"
	ledstripIDPath        = "/api/ledstrip/{id}"
	ledstripIDProfilePath = "/api/ledstrip/{id}/profile"
)

// Route a route definition
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes a route array definition
type Routes []Route

// NewRouter initializes a new router, setup with all routes
func NewRouter(enableDebug bool) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		if enableDebug {
			handler = RequestLogger(handler)
		}

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	// initialize static data
	var handler http.Handler
	handler = http.StripPrefix("/", http.FileServer(http.Dir("static")))
	if enableDebug {
		handler = RequestLogger(handler)
	}

	router.PathPrefix("/").Handler(handler)

	return router
}

// RequestLogger logs the request and duration
func RequestLogger(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf("[ %s ] %s %s", r.Method, r.RequestURI, time.Since(start))
	})
}

var routes = Routes{

	Route{"GetColorprofiles", http.MethodGet, profilePath, GetAllColorProfiles},

	Route{"CreateColorprofile", http.MethodPost, profilePath, CreateColorProfile},

	Route{"GetColorprofile", http.MethodGet, profileIDPath, GetColorProfile},

	Route{"UpdateColorprofile", http.MethodPut, profileIDPath, UpdateColorProfile},

	Route{"DeleteColorprofile", http.MethodDelete, profileIDPath, DeleteColorProfile},

	// Route{ "ApiHealthGet", http.MethodGet, "/api/health", ApiHealthGet },

	Route{"GetLedstrips", http.MethodGet, ledstripPath, GetAllLedStrips},

	Route{"CreateLedstrip", http.MethodPost, ledstripPath, CreateLedStrip},

	Route{"GetLedstrip", http.MethodGet, ledstripIDPath, GetLedStrip},

	Route{"UpdateLedstrip", http.MethodPut, ledstripIDPath, UpdateLedStrip},

	Route{"DeleteLedstripId", http.MethodDelete, ledstripIDPath, DeleteLedStrip},

	Route{"GetLedstripReferencedProfile", http.MethodGet, ledstripIDProfilePath, GetProfileForStrip},

	Route{"UpdateLedstripReferencedProfile", http.MethodPut, ledstripIDProfilePath, UpdateProfileForStrip},

	Route{"DeleteLedstripReferencedProfile", http.MethodDelete, ledstripIDProfilePath, RemoveProfileForStrip},
}
