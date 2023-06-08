package handlers

import "github.com/gorilla/mux"

type APIHandler interface {
	EnrichRoutes(*mux.Router)
}
