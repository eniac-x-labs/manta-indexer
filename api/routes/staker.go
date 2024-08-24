package routes

import (
	"net/http"

	"github.com/ethereum/go-ethereum/log"
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
	params, err := h.svc.QueryListParams(pageQuery, pageSizeQuery, order)
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

func (h Routes) ListStakerDepositStrategyHandler(w http.ResponseWriter, r *http.Request) {
	staker := r.URL.Query().Get("staker")
	pageQuery := r.URL.Query().Get("page")
	pageSizeQuery := r.URL.Query().Get("pageSize")
	order := r.URL.Query().Get("order")
	params, err := h.svc.QueryAddressListParams(staker, pageQuery, pageSizeQuery, order)
	if err != nil {
		http.Error(w, "invalid query params", http.StatusBadRequest)
		log.Error("error reading request params", "err", err.Error())
		return
	}

	tempList, err := h.svc.ListStakerDepositStrategy(params)
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

func (h Routes) ListStakerDelegatedHandler(w http.ResponseWriter, r *http.Request) {
	staker := r.URL.Query().Get("staker")
	pageQuery := r.URL.Query().Get("page")
	pageSizeQuery := r.URL.Query().Get("pageSize")
	order := r.URL.Query().Get("order")
	params, err := h.svc.QueryAddressListParams(staker, pageQuery, pageSizeQuery, order)
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

func (h Routes) ListStakerUndelegatedHandler(w http.ResponseWriter, r *http.Request) {
	staker := r.URL.Query().Get("staker")
	pageQuery := r.URL.Query().Get("page")
	pageSizeQuery := r.URL.Query().Get("pageSize")
	order := r.URL.Query().Get("order")
	params, err := h.svc.QueryAddressListParams(staker, pageQuery, pageSizeQuery, order)
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

func (h Routes) ListStakeHolderWithdrawalQueuedHandler(w http.ResponseWriter, r *http.Request) {
	staker := r.URL.Query().Get("staker")
	pageQuery := r.URL.Query().Get("page")
	pageSizeQuery := r.URL.Query().Get("pageSize")
	order := r.URL.Query().Get("order")
	params, err := h.svc.QueryAddressListParams(staker, pageQuery, pageSizeQuery, order)
	if err != nil {
		http.Error(w, "invalid query params", http.StatusBadRequest)
		log.Error("error reading request params", "err", err.Error())
		return
	}
	tempList, err := h.svc.ListStakerWithdrawalQueued(params)
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

func (h Routes) ListStakeHolderWithdrawalCompletedHandler(w http.ResponseWriter, r *http.Request) {
	staker := r.URL.Query().Get("staker")
	pageQuery := r.URL.Query().Get("page")
	pageSizeQuery := r.URL.Query().Get("pageSize")
	order := r.URL.Query().Get("order")
	params, err := h.svc.QueryAddressListParams(staker, pageQuery, pageSizeQuery, order)
	if err != nil {
		http.Error(w, "invalid query params", http.StatusBadRequest)
		log.Error("error reading request params", "err", err.Error())
		return
	}

	tempList, err := h.svc.ListStakerWithdrawalCompleted(params)
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

func (h Routes) ListStakeHolderClaimRewardHandler(w http.ResponseWriter, r *http.Request) {
	staker := r.URL.Query().Get("staker")
	pageQuery := r.URL.Query().Get("page")
	pageSizeQuery := r.URL.Query().Get("pageSize")
	order := r.URL.Query().Get("order")
	params, err := h.svc.QueryAddressListParams(staker, pageQuery, pageSizeQuery, order)
	if err != nil {
		http.Error(w, "invalid query params", http.StatusBadRequest)
		log.Error("error reading request params", "err", err.Error())
		return
	}

	tempList, err := h.svc.ListStakeHolderClaimReward(params)
	if err != nil {
		http.Error(w, "Internal server error reading ListStakeHolderClaimRewardHandler", http.StatusInternalServerError)
		log.Error("Unable to read ListStakeHolderClaimRewardHandler from DB", "err", err.Error())
		return
	}

	err = jsonResponse(w, tempList, http.StatusOK)
	if err != nil {
		log.Error("Error writing response", "err", err.Error())
	}
}
