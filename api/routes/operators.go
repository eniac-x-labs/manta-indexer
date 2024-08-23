package routes

import (
	"net/http"

	"github.com/ethereum/go-ethereum/log"
)

func (h Routes) GetOperator(w http.ResponseWriter, r *http.Request) {
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
	params, err := h.svc.QueryDTListParams(pageQuery, pageSizeQuery, order)
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
