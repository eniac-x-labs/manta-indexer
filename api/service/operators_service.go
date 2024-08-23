package service

import (
	"github.com/eniac-x-labs/manta-indexer/api/models"
	"github.com/eniac-x-labs/manta-indexer/database/worker"
	"strings"
)

func (h HandlerSvc) GetOperator(operator string) (*worker.Operators, error) {
	operatorAddress := strings.ToLower(operator)
	operatorDetail, err := h.operatorsView.GetOperator(operatorAddress)
	if err != nil {
		return &worker.Operators{}, err
	}
	return operatorDetail, err
}

func (h HandlerSvc) ListOperator(params *models.QueryDTParams) (*models.OperatorListResponse, error) {
	operatorList, total := h.operatorsView.ListOperator(params.Page, params.PageSize, params.Order)
	return &models.OperatorListResponse{
		ListResponse: models.ListResponse{
			Current: params.Page,
			Size:    params.PageSize,
			Total:   total,
		},
		Records: operatorList,
	}, nil
}
