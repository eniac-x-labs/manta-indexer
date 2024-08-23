package routes

import (
	"github.com/ethereum/go-ethereum/log"
	"net/http"
)

func (h Routes) GetStakeHolder(w http.ResponseWriter, r *http.Request) {
	staker := r.URL.Query().Get("staker")
	if staker == "" {
		http.Error(w, "invalid query params", http.StatusBadRequest)
		log.Error("error GetStakeHolder reading request params")
		return
	}
	temp, err := h.svc.GetStakeHolder(staker)
	if err != nil {
		http.Error(w, "Internal server error reading GetStakeHolder", http.StatusInternalServerError)
		log.Error("Unable to read GetStakeHolder from DB", "err", err.Error())
		return
	}
	err = jsonResponse(w, temp, http.StatusOK)
	if err != nil {
		log.Error("Error writing response", "err", err.Error())
	}
}

func (h Routes) ListStakeHolderHandler(w http.ResponseWriter, r *http.Request) {
	pageQuery := r.URL.Query().Get("page")
	pageSizeQuery := r.URL.Query().Get("pageSize")
	order := r.URL.Query().Get("order")
	params, err := h.svc.QueryDTListParams(pageQuery, pageSizeQuery, order)
	if err != nil {
		http.Error(w, "invalid query params", http.StatusBadRequest)
		log.Error("error reading request params", "err", err.Error())
		return
	}

	tempList, err := h.svc.ListStakeHolder(params)
	if err != nil {
		http.Error(w, "Internal server error reading ListStakeHolderHandler", http.StatusInternalServerError)
		log.Error("Unable to read ListStakeHolderHandler from DB", "err", err.Error())
		return
	}

	err = jsonResponse(w, tempList, http.StatusOK)
	if err != nil {
		log.Error("Error writing response", "err", err.Error())
	}
}
