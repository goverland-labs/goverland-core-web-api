package handlers

import (
	"encoding/json"
	"net/http"

	forms "github.com/goverland-labs/goverland-core-web-api/internal/rest/form/ens"

	"github.com/gorilla/mux"
	"github.com/goverland-labs/goverland-core-storage/protocol/storagepb"
	"github.com/rs/zerolog/log"

	"github.com/goverland-labs/goverland-core-web-api/internal/response"
	"github.com/goverland-labs/goverland-core-web-api/internal/rest/models/ens"
)

type Ens struct {
	ec storagepb.EnsClient
}

func NewEnsHandler(ec storagepb.EnsClient) APIHandler {
	return &Ens{
		ec: ec,
	}
}

func (h *Ens) EnrichRoutes(v1, _ *mux.Router) {
	v1.HandleFunc("/ens-name", h.getEnsNamesAction).Methods(http.MethodGet).Name("get_ens_names")
}

func (h *Ens) getEnsNamesAction(w http.ResponseWriter, r *http.Request) {
	form, ferr := forms.NewGetEnsNamesForm().ParseAndValidate(r)
	if ferr != nil {
		response.HandleError(ferr, w)

		return
	}

	params := form.(*forms.GetEnsNames)
	var list, err = h.ec.GetEnsByAddresses(r.Context(), &storagepb.EnsByAddressesRequest{
		Addresses: params.Addresses,
	})
	if err != nil {
		log.Error().Err(err).Fields(params.ConvertToMap()).Msg("get ens names")
		response.HandleError(response.ResolveError(err), w)

		return
	}

	resp := make([]ens.EnsName, len(list.GetEnsNames()))
	for i, info := range list.GetEnsNames() {
		resp[i] = convertToEnsNameFromProto(info)
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func convertToEnsNameFromProto(info *storagepb.EnsName) ens.EnsName {
	return ens.EnsName{
		Address: info.GetAddress(),
		Name:    info.GetName(),
	}
}
