package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// All routes related to devices
func (api *APIServer) deviceRoutes(r *mux.Router) {
	// check if auth token needed
	api.IfTokenNeeded(r)

	// get all devices
	r.HandleFunc("/get/all", api.getAllDevices).Methods(http.MethodOptions, http.MethodGet)

	// get information about the last summary of
	// devices base on an intervals like 1m, 2h ...
	r.HandleFunc("/get/summary/interval/{interval}", api.getAllDevicesSummaryByInterval).Methods(http.MethodOptions, http.MethodGet)

	// only for one device
	// get information about the last summary of
	// devices base on an intervals like 1m, 2h ...
	r.HandleFunc("/get/summary/interval/{interval}/by/{device}", api.getAllDevicesSummaryByIntervalByDevice).Methods(http.MethodOptions, http.MethodGet)

	// get a device by
	// device IP like: 127.0.0.1
	r.HandleFunc("/get/by/{device}", api.getByDevice).Methods(http.MethodOptions, http.MethodGet)

	// update a device by
	// device IP like: 127.0.0.1
	r.HandleFunc("/set/by/{device}", api.setByDevice).Methods(http.MethodOptions, http.MethodPost)

	// return an array of packets, bytes .... group by the interval that a user provides
	// r.HandleFunc("/get/summary/group/{interval}/by/{deviceID}", api.getDeviceSummaryGroupByIntervalByDeviceID).Methods(http.MethodOptions, http.MethodGet)

}
