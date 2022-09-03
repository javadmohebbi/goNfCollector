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

	// get port report only when SOURCE OF DESTINATION report based on interval
	r.HandleFunc("/report/{host}/as/{direction}/top/{top}/interval/{interval}", api.getPortReportWhenHostSrcOrDst).Methods(http.MethodOptions, http.MethodGet)

	// get a port by id to update Info in UI
	r.HandleFunc("/get/by/id/{id}", api.getPortByID).Methods(http.MethodOptions, http.MethodGet)

	// set a port by id to update Info in UI
	r.HandleFunc("/set/by/id/{id}", api.setPortByID).Methods(http.MethodOptions, http.MethodPost)
}
