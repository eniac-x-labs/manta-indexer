package routes

import (
	"net/http"

	"github.com/ethereum/go-ethereum/log"
)

func (h Routes) StrategyHandler(w http.ResponseWriter, r *http.Request) {
	strategy := r.URL.Query().Get("strategy")
	if strategy == "" {
		http.Error(w, "invalid query params", http.StatusBadRequest)
		log.Error("error strategy reading request params")
		return
	}
	temp, err := h.svc.Strategy(strategy)
	if err != nil {
		http.Error(w, "Internal server error reading GetOperator", http.StatusInternalServerError)
		log.Error("Unable to read strategy from DB", "err", err.Error())
		return
	}
	err = jsonResponse(w, temp, http.StatusOK)
	if err != nil {
		log.Error("Error writing response", "err", err.Error())
	}
}

func (h Routes) StrategyListHandler(w http.ResponseWriter, r *http.Request) {
	pageQuery := r.URL.Query().Get("page")
	pageSizeQuery := r.URL.Query().Get("pageSize")
	order := r.URL.Query().Get("order")
	params, err := h.svc.QueryListParams(pageQuery, pageSizeQuery, order)
	if err != nil {
		http.Error(w, "invalid query params", http.StatusBadRequest)
		log.Error("error reading request params", "err", err.Error())
		return
	}

	tempList, err := h.svc.StrategyList(params)
	if err != nil {
		http.Error(w, "Internal server error reading ListOperatorHandler", http.StatusInternalServerError)
		log.Error("Unable to read ListOperatorHandler from DB", "err", err.Error())
		return
	}

	err = jsonResponse(w, tempList, http.StatusOK)
	if err != nil {
		log.Error("Error writing response", "err", err.Error())
	}
}

func (h Routes) GetOperatorHandler(w http.ResponseWriter, r *http.Request) {
	operator := r.URL.Query().Get("operator")
	if operator == "" {
		http.Error(w, "invalid query params", http.StatusBadRequest)
		log.Error("error GetOperator reading request params")
		return
	}
	temp, err := h.svc.GetOperator(operator)
	if err != nil {
		http.Error(w, "Internal server error reading GetOperator", http.StatusInternalServerError)
		log.Error("Unable to read GetOperator from DB", "err", err.Error())
		return
	}
	err = jsonResponse(w, temp, http.StatusOK)
	if err != nil {
		log.Error("Error writing response", "err", err.Error())
	}
}

func (h Routes) ListOperatorHandler(w http.ResponseWriter, r *http.Request) {
	pageQuery := r.URL.Query().Get("page")
	pageSizeQuery := r.URL.Query().Get("pageSize")
	order := r.URL.Query().Get("order")
	params, err := h.svc.QueryListParams(pageQuery, pageSizeQuery, order)
	if err != nil {
		http.Error(w, "invalid query params", http.StatusBadRequest)
		log.Error("error reading request params", "err", err.Error())
		return
	}

	tempList, err := h.svc.ListOperator(params)
	if err != nil {
		http.Error(w, "Internal server error reading ListOperatorHandler", http.StatusInternalServerError)
		log.Error("Unable to read ListOperatorHandler from DB", "err", err.Error())
		return
	}

	err = jsonResponse(w, tempList, http.StatusOK)
	if err != nil {
		log.Error("Error writing response", "err", err.Error())
	}
}

func (h Routes) RegisterOperatorHandler(w http.ResponseWriter, r *http.Request) {
	operator := r.URL.Query().Get("operator")
	if operator == "" {
		http.Error(w, "invalid query params", http.StatusBadRequest)
		log.Error("error reading request params")
		return
	}
	registerOperatorRet, err := h.svc.RegisterOperator(operator)
	if err != nil {
		http.Error(w, "Internal server error reading register operator", http.StatusInternalServerError)
		log.Error("Unable to read register operator from DB", "err", err.Error())
		return
	}
	err = jsonResponse(w, registerOperatorRet, http.StatusOK)
	if err != nil {
		log.Error("Error writing response", "err", err.Error())
	}
}

func (h Routes) RegisterOperatorListHandler(w http.ResponseWriter, r *http.Request) {
	pageQuery := r.URL.Query().Get("page")
	pageSizeQuery := r.URL.Query().Get("pageSize")
	order := r.URL.Query().Get("order")
	params, err := h.svc.QueryListParams(pageQuery, pageSizeQuery, order)
	if err != nil {
		http.Error(w, "invalid query params", http.StatusBadRequest)
		log.Error("error reading request params", "err", err.Error())
		return
	}

	registerOperatorsRet, err := h.svc.ListRegisterOperator(params)
	if err != nil {
		http.Error(w, "Internal server error reading register operator list", http.StatusInternalServerError)
		log.Error("Unable to read register operator list from DB", "err", err.Error())
		return
	}

	err = jsonResponse(w, registerOperatorsRet, http.StatusOK)
	if err != nil {
		log.Error("Error writing response", "err", err.Error())
	}
}

func (h Routes) ListOperatorNodeUrlUpdateHandler(w http.ResponseWriter, r *http.Request) {
	operator := r.URL.Query().Get("operator")
	pageQuery := r.URL.Query().Get("page")
	pageSizeQuery := r.URL.Query().Get("pageSize")
	order := r.URL.Query().Get("order")
	params, err := h.svc.QueryAddressListParams(operator, pageQuery, pageSizeQuery, order)
	if err != nil {
		http.Error(w, "invalid query params", http.StatusBadRequest)
		log.Error("error reading request params", "err", err.Error())
		return
	}
	tempList, err := h.svc.ListOperatorNodeUrlUpdate(params)
	if err != nil {
		http.Error(w, "Internal server error reading ListOperatorListHandler", http.StatusInternalServerError)
		log.Error("Unable to read ListOperatorListHandler from DB", "err", err.Error())
		return
	}
	err = jsonResponse(w, tempList, http.StatusOK)
	if err != nil {
		log.Error("Error writing response", "err", err.Error())
	}
}

func (h Routes) ListOperatorSharesIncreasedHandler(w http.ResponseWriter, r *http.Request) {
	operator := r.URL.Query().Get("operator")
	pageQuery := r.URL.Query().Get("page")
	pageSizeQuery := r.URL.Query().Get("pageSize")
	order := r.URL.Query().Get("order")
	params, err := h.svc.QueryAddressListParams(operator, pageQuery, pageSizeQuery, order)
	if err != nil {
		http.Error(w, "invalid query params", http.StatusBadRequest)
		log.Error("error reading request params", "err", err.Error())
		return
	}

	tempList, err := h.svc.ListOperatorSharesIncreased(params)
	if err != nil {
		http.Error(w, "Internal server error reading ListOperatorSharesIncreasedHandler", http.StatusInternalServerError)
		log.Error("Unable to read ListOperatorSharesIncreasedHandler from DB", "err", err.Error())
		return
	}

	err = jsonResponse(w, tempList, http.StatusOK)
	if err != nil {
		log.Error("Error writing response", "err", err.Error())
	}
}

func (h Routes) ListOperatorReceiveStakerDelegateHandler(w http.ResponseWriter, r *http.Request) {
	operator := r.URL.Query().Get("operator")
	pageQuery := r.URL.Query().Get("page")
	pageSizeQuery := r.URL.Query().Get("pageSize")
	order := r.URL.Query().Get("order")
	params, err := h.svc.QueryAddressListParams(operator, pageQuery, pageSizeQuery, order)
	if err != nil {
		http.Error(w, "invalid query params", http.StatusBadRequest)
		log.Error("error reading request params", "err", err.Error())
		return
	}

	tempList, err := h.svc.ListOperatorSharesDecreased(params)
	if err != nil {
		http.Error(w, "Internal server error reading ListOperatorSharesDecreasedHandler", http.StatusInternalServerError)
		log.Error("Unable to read ListOperatorSharesDecreasedHandler from DB", "err", err.Error())
		return
	}

	err = jsonResponse(w, tempList, http.StatusOK)
	if err != nil {
		log.Error("Error writing response", "err", err.Error())
	}

}

func (h Routes) ListOperatorSharesDecreasedHandler(w http.ResponseWriter, r *http.Request) {
	operator := r.URL.Query().Get("operator")
	pageQuery := r.URL.Query().Get("page")
	pageSizeQuery := r.URL.Query().Get("pageSize")
	order := r.URL.Query().Get("order")
	params, err := h.svc.QueryAddressListParams(operator, pageQuery, pageSizeQuery, order)
	if err != nil {
		http.Error(w, "invalid query params", http.StatusBadRequest)
		log.Error("error reading request params", "err", err.Error())
		return
	}

	tempList, err := h.svc.ListOperatorSharesDecreased(params)
	if err != nil {
		http.Error(w, "Internal server error reading ListOperatorSharesDecreasedHandler", http.StatusInternalServerError)
		log.Error("Unable to read ListOperatorSharesDecreasedHandler from DB", "err", err.Error())
		return
	}

	err = jsonResponse(w, tempList, http.StatusOK)
	if err != nil {
		log.Error("Error writing response", "err", err.Error())
	}
}

func (h Routes) ListOperatorAndStakeRewardHandler(w http.ResponseWriter, r *http.Request) {
	operator := r.URL.Query().Get("operator")
	pageQuery := r.URL.Query().Get("page")
	pageSizeQuery := r.URL.Query().Get("pageSize")
	order := r.URL.Query().Get("order")
	params, err := h.svc.QueryAddressListParams(operator, pageQuery, pageSizeQuery, order)
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

func (h Routes) ListOperatorClaimRewardHandler(w http.ResponseWriter, r *http.Request) {
	operator := r.URL.Query().Get("operator")
	pageQuery := r.URL.Query().Get("page")
	pageSizeQuery := r.URL.Query().Get("pageSize")
	order := r.URL.Query().Get("order")
	params, err := h.svc.QueryAddressListParams(operator, pageQuery, pageSizeQuery, order)
	if err != nil {
		http.Error(w, "invalid query params", http.StatusBadRequest)
		log.Error("error reading request params", "err", err.Error())
		return
	}
	tempList, err := h.svc.ListOperatorClaimReward(params)
	if err != nil {
		http.Error(w, "Internal server error reading ListOperatorClaimRewardHandler", http.StatusInternalServerError)
		log.Error("Unable to read ListOperatorClaimRewardHandler from DB", "err", err.Error())
		return
	}

	err = jsonResponse(w, tempList, http.StatusOK)
	if err != nil {
		log.Error("Error writing response", "err", err.Error())
	}
}
