package routes

import (
	"net/http"

	"github.com/ethereum/go-ethereum/log"
)

func (h Routes) GetWithdrawalQueued(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")
	if guid == "" {
		http.Error(w, "invalid query params", http.StatusBadRequest)
		log.Error("error GetWithdrawalQueued reading request params")
		return
	}
	temp, err := h.svc.GetWithdrawalQueued(guid)
	if err != nil {
		http.Error(w, "Internal server error reading GetWithdrawalQueued", http.StatusInternalServerError)
		log.Error("Unable to read GetWithdrawalQueued from DB", "err", err.Error())
		return
	}
	err = jsonResponse(w, temp, http.StatusOK)
	if err != nil {
		log.Error("Error writing response", "err", err.Error())
	}
}

func (h Routes) ListWithdrawalQueuedHandler(w http.ResponseWriter, r *http.Request) {
	pageQuery := r.URL.Query().Get("page")
	pageSizeQuery := r.URL.Query().Get("pageSize")
	order := r.URL.Query().Get("order")
	params, err := h.svc.QueryDTListParams(pageQuery, pageSizeQuery, order)
	if err != nil {
		http.Error(w, "invalid query params", http.StatusBadRequest)
		log.Error("error reading request params", "err", err.Error())
		return
	}

	tempList, err := h.svc.ListWithdrawalQueued(params)
	if err != nil {
		http.Error(w, "Internal server error reading ListWithdrawalQueuedHandler", http.StatusInternalServerError)
		log.Error("Unable to read ListWithdrawalQueuedHandler from DB", "err", err.Error())
		return
	}

	err = jsonResponse(w, tempList, http.StatusOK)
	if err != nil {
		log.Error("Error writing response", "err", err.Error())
	}
}
