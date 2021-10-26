package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// All routes related to ethernets
func (api *APIServer) ethernetRoutes(r *mux.Router) {
	// check if auth token needed
	api.IfTokenNeeded(r)

	// get information about the last summary of
	// ethernets base on an intervals like 1m, 2h ...
	r.HandleFunc("/get/device/{deviceID}/interval/{interval}", api.getEthernetsByDeviceByInterval).Methods(http.MethodOptions, http.MethodGet)

}
