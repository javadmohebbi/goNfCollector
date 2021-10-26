package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// All routes related to threats
func (api *APIServer) threatRoutes(r *mux.Router) {
	// check if auth token needed
	api.IfTokenNeeded(r)

	// get all
	r.HandleFunc("/get/all", api.getAllThreats).Methods(http.MethodOptions, http.MethodGet)

}
