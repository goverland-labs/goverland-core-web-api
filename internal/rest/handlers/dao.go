package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/goverland-labs/core-api/protobuf/internalapi"
)

type DAO struct {
	dc internalapi.DaoClient
}

func NewDaoHandler(dc internalapi.DaoClient) APIHandler {
	return &DAO{
		dc: dc,
	}
}

func (h *DAO) EnrichRoutes(baseRouter *mux.Router) {
	baseRouter.HandleFunc("/dao/{id}", h.getByIDAction).Methods(http.MethodGet).Name("get_dao_by_id")
}

func (h *DAO) getByIDAction(w http.ResponseWriter, r *http.Request) {
	dao, err := h.dc.GetByID(r.Context(), &internalapi.DaoByIDRequest{DaoId: "7771afe2-d58a-4308-8119-9594f6abb2ee"})
	if err != nil {
		panic(err)
	}

	fmt.Println(dao.GetDao().GetId())

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"dao_id": dao.GetDao().GetId()})
}
