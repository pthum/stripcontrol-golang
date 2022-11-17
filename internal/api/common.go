package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pthum/stripcontrol-golang/internal/model"
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

// Common handler methods

// BindJSON bind the response body to the object
func bindJSON(r *http.Request, obj interface{}) (err error) {
	byteData, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(byteData, &obj)
	return
}

// GetParam get the specified param
func getParam(r *http.Request, param string) (paramValue string) {
	vars := mux.Vars(r)
	paramValue = vars[param]
	return
}

// HandleError handles an error
func handleError(w *http.ResponseWriter, status int, message string) {
	handleJSON(w, status, H{"error": message})
}

// HandleJSON handle json
func handleJSON(w *http.ResponseWriter, status int, result interface{}) {
	writer := *w

	marshalled, err := json.Marshal(result)

	if err != nil {
		handleJSON(w, http.StatusInternalServerError, H{"error": err.Error()})
		return
	}

	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(status)
	_, err = writer.Write(marshalled)
	if err != nil {
		log.Printf("error writing response: %v", err)
	}
}

func respondWithCreated(r *http.Request, w http.ResponseWriter, input model.IDer) {
	log.Printf("ID after save %d", input.GetID())
	w.Header().Add("Location", fmt.Sprintf("%s/%d", r.RequestURI, input.GetID()))
	handleJSON(&w, http.StatusCreated, input)
}

// H is a shortcut for map[string]interface{}
type H map[string]interface{}
