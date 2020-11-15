package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// Common handler methods

// BindJSON bind the response body to the object
func BindJSON(r *http.Request, obj interface{}) (err error) {
	byteData, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(byteData, &obj)
	return
}

// GetParam get the specified param
func GetParam(r *http.Request, param string) (paramValue string) {
	vars := mux.Vars(r)
	paramValue = vars[param]
	return
}

// HandleJSON handle json
func HandleJSON(w *http.ResponseWriter, status int, result interface{}) {
	writer := *w

	marshalled, err := json.Marshal(result)

	if err != nil {
		HandleJSON(w, 500, H{"error": err.Error()})
		return
	}

	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(status)
	writer.Write(marshalled)
}

// H is a shortcut for map[string]interface{}
type H map[string]interface{}
