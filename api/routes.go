package api

import "github.com/gorilla/mux"

// ALL API Routes
func (api *APIServer) routes(r *mux.Router) *mux.Router {

	// all routes related to devices
	deviceRoutes := r.PathPrefix("/devices").Subrouter()
	api.deviceRoutes(deviceRoutes)

	return r
}
