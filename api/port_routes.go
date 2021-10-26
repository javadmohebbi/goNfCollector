package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// All routes related to ports
func (api *APIServer) portRoutes(r *mux.Router) {
	// check if auth token needed
	api.IfTokenNeeded(r)

	// get all ports
	r.HandleFunc("/get/all", api.getAllPorts).Methods(http.MethodOptions, http.MethodGet)

	// get information about the last summary of
	// top ports base on an intervals like 1m, 2h ...
	r.HandleFunc("/get/top/{top}/device/{deviceID}/direction/{direction}/interval/{interval}", api.getTopPortsByDeviceByInterval).Methods(http.MethodOptions, http.MethodGet)

}
