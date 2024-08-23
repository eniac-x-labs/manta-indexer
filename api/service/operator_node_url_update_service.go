package service

import (
	"strings"

	"github.com/eniac-x-labs/manta-indexer/api/models"
	"github.com/eniac-x-labs/manta-indexer/database/event"
)

func (h HandlerSvc) GetOperatorNodeUrlUpdate(operator string) (*event.OperatorNodeUrlUpdate, error) {
	addressToLower := strings.ToLower(operator)
	operatorNodeUrlUpdate, err := h.operatorNodeUrlUpdateView.GetOperatorNodeUrlUpdate(addressToLower)
	if err != nil {
		return &event.OperatorNodeUrlUpdate{}, err
	}
	return operatorNodeUrlUpdate, err
}

func (h HandlerSvc) ListOperatorNodeUrlUpdate(params *models.QueryDTParams) (*models.OperatorNodeUrlUpdateListResponse, error) {
	operatorNodeUrlUpdateList, total := h.operatorNodeUrlUpdateView.ListOperatorNodeUrlUpdate(params.Page, params.PageSize, params.Order)
	return &models.OperatorNodeUrlUpdateListResponse{
		ListResponse: models.ListResponse{
			Current: params.Page,
			Size:    params.PageSize,
			Total:   total,
		},
		Records: operatorNodeUrlUpdateList,
	}, nil
}
