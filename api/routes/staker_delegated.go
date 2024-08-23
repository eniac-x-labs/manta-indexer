package routes

import (
	"github.com/ethereum/go-ethereum/log"
	"net/http"
)

func (h Routes) GetStakerDelegated(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")
	if guid == "" {
		http.Error(w, "invalid query params", http.StatusBadRequest)
		log.Error("error GetStakerDelegated reading request params")
		return
	}
	temp, err := h.svc.GetStakerDelegated(guid)
	if err != nil {
		http.Error(w, "Internal server error reading GetStakerDelegated", http.StatusInternalServerError)
		log.Error("Unable to read GetStakerDelegated from DB", "err", err.Error())
		return
	}
	err = jsonResponse(w, temp, http.StatusOK)
	if err != nil {
		log.Error("Error writing response", "err", err.Error())
	}
}

func (h Routes) ListStakerDelegatedHandler(w http.ResponseWriter, r *http.Request) {
	pageQuery := r.URL.Query().Get("page")
	pageSizeQuery := r.URL.Query().Get("pageSize")
	order := r.URL.Query().Get("order")
	params, err := h.svc.QueryDTListParams(pageQuery, pageSizeQuery, order)
	if err != nil {
		http.Error(w, "invalid query params", http.StatusBadRequest)
		log.Error("error reading request params", "err", err.Error())
		return
	}

	tempList, err := h.svc.ListStakerDelegated(params)
	if err != nil {
		http.Error(w, "Internal server error reading ListStakerDelegatedHandler", http.StatusInternalServerError)
		log.Error("Unable to read ListStakerDelegatedHandler from DB", "err", err.Error())
		return
	}

	err = jsonResponse(w, tempList, http.StatusOK)
	if err != nil {
		log.Error("Error writing response", "err", err.Error())
	}
}
