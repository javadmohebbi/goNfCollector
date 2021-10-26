package api

import "github.com/gorilla/mux"

// if auth needed, it will validate tokens
func (a *APIServer) IfTokenNeeded(r *mux.Router) {
	if a.apiConf.NeedAuth {
		r.Use(a.validateToken)
	}
}
