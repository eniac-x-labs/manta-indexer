package routes

import (
	"net/http"

	"github.com/ethereum/go-ethereum/log"
)

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
	params, err := h.svc.QueryDTListParams(pageQuery, pageSizeQuery, order)
	if err != nil {
		http.Error(w, "invalid query params", http.StatusBadRequest)
		log.Error("error reading request params", "err", err.Error())
		return
	}

	registerOperatorsRet, err := h.svc.RegisterOperatorList(params)
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
