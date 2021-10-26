package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// All routes related to flows
func (api *APIServer) flowsRoutes(r *mux.Router) {
	// check if auth token needed
	api.IfTokenNeeded(r)

	// get all threat's flows
	// by theatID and interval
	r.HandleFunc("/get/all/threat/{threatID}/interval/{interval}", api.getAllFlowsByThreatID).Methods(http.MethodOptions, http.MethodGet)

}
