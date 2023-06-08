package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type DAO struct {
}

func NewDaoHandler() APIHandler {
	return &DAO{}
}

func (h *DAO) EnrichRoutes(baseRouter *mux.Router) {
	baseRouter.HandleFunc("/dao/{id}", h.getByIDAction).Methods(http.MethodGet).Name("get_dao_by_id")
}

func (h *DAO) getByIDAction(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "implement me"})
}
