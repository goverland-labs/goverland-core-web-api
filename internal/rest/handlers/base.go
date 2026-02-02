package handlers

import "github.com/gorilla/mux"

type APIHandler interface {
	EnrichRoutes(v1, v2 *mux.Router)
}
