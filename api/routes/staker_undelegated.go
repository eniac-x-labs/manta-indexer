package routes

import (
	"net/http"

	"github.com/ethereum/go-ethereum/log"
)

func (h Routes) GetStakerUndelegated(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")
	if guid == "" {
		http.Error(w, "invalid query params", http.StatusBadRequest)
		log.Error("error GetStakerUndelegated reading request params")
		return
	}
	temp, err := h.svc.GetStakerUndelegated(guid)
	if err != nil {
		http.Error(w, "Internal server error reading GetStakerUndelegated", http.StatusInternalServerError)
		log.Error("Unable to read GetStakerUndelegated from DB", "err", err.Error())
		return
	}
	err = jsonResponse(w, temp, http.StatusOK)
	if err != nil {
		log.Error("Error writing response", "err", err.Error())
	}
}

func (h Routes) ListStakerUndelegatedHandler(w http.ResponseWriter, r *http.Request) {
	pageQuery := r.URL.Query().Get("page")
	pageSizeQuery := r.URL.Query().Get("pageSize")
	order := r.URL.Query().Get("order")
	params, err := h.svc.QueryDTListParams(pageQuery, pageSizeQuery, order)
	if err != nil {
		http.Error(w, "invalid query params", http.StatusBadRequest)
		log.Error("error reading request params", "err", err.Error())
		return
	}

	tempList, err := h.svc.ListStakerUndelegated(params)
	if err != nil {
		http.Error(w, "Internal server error reading ListStakerUndelegatedHandler", http.StatusInternalServerError)
		log.Error("Unable to read ListStakerUndelegatedHandler from DB", "err", err.Error())
		return
	}

	err = jsonResponse(w, tempList, http.StatusOK)
	if err != nil {
		log.Error("Error writing response", "err", err.Error())
	}
}
