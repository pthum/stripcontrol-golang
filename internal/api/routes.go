package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	cpapi "github.com/pthum/stripcontrol-golang/internal/api/colorprofile"
	api "github.com/pthum/stripcontrol-golang/internal/api/common"
	lapi "github.com/pthum/stripcontrol-golang/internal/api/led"
)

// NewRouter initializes a new router, setup with all routes
func NewRouter(enableDebug bool) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	var routes []api.Route
	var cproutes = cpapi.ColorProfileRoutes()
	var lroutes = lapi.LEDRoutes()
	routes = append(routes, cproutes...)
	routes = append(routes, lroutes...)

	for _, route := range routes {
		fmt.Printf("appending \"%v\": %v %v \n", route.Name, route.Method, route.Pattern)
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
