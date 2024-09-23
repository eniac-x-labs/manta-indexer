package routes

import (
	"net/http"

	"github.com/ethereum/go-ethereum/log"

	"github.com/eniac-x-labs/manta-indexer/common/bigint"
)

func (h Routes) GetFinalityVerifiedHandler(w http.ResponseWriter, r *http.Request) {
	l2BN := r.URL.Query().Get("l2BlockNum")
	if l2BN == "" {
		http.Error(w, "invalid query params", http.StatusBadRequest)
		log.Error("error strategy reading request params")
		return
	}
	temp, err := h.svc.GetFinalityVerified(bigint.StringToBigInt(l2BN))
	if err != nil {
		http.Error(w, "Internal server error reading GetFinalityVerified", http.StatusInternalServerError)
		log.Error("Unable to read finality verified from DB", "err", err.Error())
		return
	}
	err = jsonResponse(w, temp, http.StatusOK)
	if err != nil {
		log.Error("Error writing response", "err", err.Error())
	}
}
