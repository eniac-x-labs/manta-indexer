package routes

import (
	"github.com/ethereum/go-ethereum/log"
	"net/http"
)

func (h Routes) GetOperatorAndStakeReward(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")
	if guid == "" {
		http.Error(w, "invalid query params", http.StatusBadRequest)
		log.Error("error GetOperatorAndStakeReward reading request params")
		return
	}
	reward, err := h.svc.GetOperatorAndStakeReward(guid)
	if err != nil {
		http.Error(w, "Internal server error reading GetOperatorAndStakeReward", http.StatusInternalServerError)
		log.Error("Unable to read GetOperatorAndStakeReward from DB", "err", err.Error())
		return
	}
	err = jsonResponse(w, reward, http.StatusOK)
	if err != nil {
		log.Error("Error writing response", "err", err.Error())
	}
}

func (h Routes) ListOperatorAndStakeRewardHandler(w http.ResponseWriter, r *http.Request) {
	pageQuery := r.URL.Query().Get("page")
	pageSizeQuery := r.URL.Query().Get("pageSize")
	order := r.URL.Query().Get("order")
	params, err := h.svc.QueryDTListParams(pageQuery, pageSizeQuery, order)
	if err != nil {
		http.Error(w, "invalid query params", http.StatusBadRequest)
		log.Error("error reading request params", "err", err.Error())
		return
	}

	tempList, err := h.svc.ListOperatorAndStakeReward(params)
	if err != nil {
		http.Error(w, "Internal server error reading ListOperatorAndStakeRewardHandler", http.StatusInternalServerError)
		log.Error("Unable to read ListOperatorAndStakeRewardHandler from DB", "err", err.Error())
		return
	}

	err = jsonResponse(w, tempList, http.StatusOK)
	if err != nil {
		log.Error("Error writing response", "err", err.Error())
	}
}
