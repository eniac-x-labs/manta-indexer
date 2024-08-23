package routes

import (
	"github.com/ethereum/go-ethereum/log"
	"net/http"
)

func (h Routes) GetWithdrawalCompleted(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")
	if guid == "" {
		http.Error(w, "invalid query params", http.StatusBadRequest)
		log.Error("error GetWithdrawalCompleted reading request params")
		return
	}
	temp, err := h.svc.GetWithdrawalCompleted(guid)
	if err != nil {
		http.Error(w, "Internal server error reading GetWithdrawalCompleted", http.StatusInternalServerError)
		log.Error("Unable to read GetWithdrawalCompleted from DB", "err", err.Error())
		return
	}
	err = jsonResponse(w, temp, http.StatusOK)
	if err != nil {
		log.Error("Error writing response", "err", err.Error())
	}
}

func (h Routes) ListWithdrawalCompletedHandler(w http.ResponseWriter, r *http.Request) {
	pageQuery := r.URL.Query().Get("page")
	pageSizeQuery := r.URL.Query().Get("pageSize")
	order := r.URL.Query().Get("order")
	params, err := h.svc.QueryDTListParams(pageQuery, pageSizeQuery, order)
	if err != nil {
		http.Error(w, "invalid query params", http.StatusBadRequest)
		log.Error("error reading request params", "err", err.Error())
		return
	}

	tempList, err := h.svc.ListWithdrawalCompleted(params)
	if err != nil {
		http.Error(w, "Internal server error reading ListWithdrawalCompletedHandler", http.StatusInternalServerError)
		log.Error("Unable to read ListWithdrawalCompletedHandler from DB", "err", err.Error())
		return
	}

	err = jsonResponse(w, tempList, http.StatusOK)
	if err != nil {
		log.Error("Error writing response", "err", err.Error())
	}
}
