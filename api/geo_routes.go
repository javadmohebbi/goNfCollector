package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// All routes related to geolocations
func (api *APIServer) geoRoutes(r *mux.Router) {
	// check if auth token needed
	api.IfTokenNeeded(r)

	// get all geos
	r.HandleFunc("/get/all", api.getAllGeos).Methods(http.MethodOptions, http.MethodGet)

	// get information about the last summary of
	// top geo base on an intervals like 1m, 2h ...
	r.HandleFunc("/get/country/top/{top}/device/{deviceID}/direction/{direction}/interval/{interval}", api.getTopGeoCountryByDeviceByInterval).Methods(http.MethodOptions, http.MethodGet)

	// get geo report only when SOURCE OF DESTINATION report based on interval
	r.HandleFunc("/report/{host}/as/{direction}/top/{top}/interval/{interval}", api.getGeoReportWhenHostSrcOrDst).Methods(http.MethodOptions, http.MethodGet)

}
