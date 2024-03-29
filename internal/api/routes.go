package api

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	alog "github.com/pthum/stripcontrol-golang/internal/log"
	"github.com/samber/do"
)

// NewRouter initializes a new router, setup with all routes
func NewRouter(i *do.Injector, enableDebug bool) *mux.Router {
	l := alog.NewLogger("router")
	router := mux.NewRouter().StrictSlash(true)
	cph := do.MustInvoke[CPHandler](i).(*cpHandlerImpl)
	lh := do.MustInvoke[LEDHandler](i).(*ledHandlerImpl)
	var routes []Route
	var cproutes = cph.colorProfileRoutes()
	var lroutes = lh.ledRoutes()
	routes = append(routes, cproutes...)
	routes = append(routes, lroutes...)

	for _, route := range routes {
		l.Info("appending \"%v\": %v %v \n", route.HandlerName(), route.Method, route.Pattern)
		var handler http.Handler
		handler = route.HandlerFunc
		if enableDebug {
			handler = RequestLogger(handler)
		}

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.HandlerName()).
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
	lg := alog.NewLogger("request")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		lg.Debug("[ %s ] %s %s", r.Method, r.RequestURI, time.Since(start))
	})
}
