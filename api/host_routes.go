package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// All routes related to hosts
func (api *APIServer) hostRoutes(r *mux.Router) {
	// check if auth token needed
	api.IfTokenNeeded(r)

	// get all hosts
	r.HandleFunc("/get/all", api.getAllHosts).Methods(http.MethodOptions, http.MethodGet)

	// get information about the last summary of
	// top hosts base on an intervals and device like 1m, 2h ... and 127.0.0.1
	r.HandleFunc("/get/top/{top}/device/{deviceID}/direction/{direction}/interval/{interval}", api.getTopHostsByDeviceByInterval).Methods(http.MethodOptions, http.MethodGet)

	// get host report based on interval
	r.HandleFunc("/report/{host}/top/{top}/interval/{interval}", api.getHostReport).Methods(http.MethodOptions, http.MethodGet)

	// get host report only when SOURCE OR DESTINATION report based on interval
	r.HandleFunc("/report/{host}/as/{direction}/top/{top}/interval/{interval}", api.getHostReportWhenSrcOrDst).Methods(http.MethodOptions, http.MethodGet)

	// get a host by id to update Name, Info in UI
	r.HandleFunc("/get/by/id/{id}", api.getHostByID).Methods(http.MethodOptions, http.MethodGet)

	r.HandleFunc("/set/by/id/{id}", api.setHostByID).Methods(http.MethodOptions, http.MethodPost)
}
