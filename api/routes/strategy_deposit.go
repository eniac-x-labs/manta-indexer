package routes

import (
	"github.com/ethereum/go-ethereum/log"
	"net/http"
)

func (h Routes) GetStrategyDeposit(w http.ResponseWriter, r *http.Request) {
	staker := r.URL.Query().Get("staker")
	if staker == "" {
		http.Error(w, "invalid query params", http.StatusBadRequest)
		log.Error("error GetStrategyDeposit reading request params")
		return
	}
	temp, err := h.svc.GetStrategyDeposit(staker)
	if err != nil {
		http.Error(w, "Internal server error reading GetStrategyDeposit", http.StatusInternalServerError)
		log.Error("Unable to read GetStrategyDeposit from DB", "err", err.Error())
		return
	}
	err = jsonResponse(w, temp, http.StatusOK)
	if err != nil {
		log.Error("Error writing response", "err", err.Error())
	}
}

func (h Routes) ListStrategyDepositHandler(w http.ResponseWriter, r *http.Request) {
	pageQuery := r.URL.Query().Get("page")
	pageSizeQuery := r.URL.Query().Get("pageSize")
	order := r.URL.Query().Get("order")
	params, err := h.svc.QueryDTListParams(pageQuery, pageSizeQuery, order)
	if err != nil {
		http.Error(w, "invalid query params", http.StatusBadRequest)
		log.Error("error reading request params", "err", err.Error())
		return
	}

	tempList, err := h.svc.ListStrategyDeposit(params)
	if err != nil {
		http.Error(w, "Internal server error reading ListStrategyDepositHandler", http.StatusInternalServerError)
		log.Error("Unable to read ListStrategyDepositHandler from DB", "err", err.Error())
		return
	}

	err = jsonResponse(w, tempList, http.StatusOK)
	if err != nil {
		log.Error("Error writing response", "err", err.Error())
	}
}
