package service

import (
	"github.com/eniac-x-labs/manta-indexer/api/models"
	"github.com/eniac-x-labs/manta-indexer/database/event"
)

func (h HandlerSvc) GetOperatorSharesIncreased(guid string) (*event.OperatorSharesIncreased, error) {
	operatorSharesIncreased, err := h.operatorSharesIncreasedView.GetOperatorSharesIncreased(guid)
	if err != nil {
		return &event.OperatorSharesIncreased{}, err
	}
	return operatorSharesIncreased, err
}

func (h HandlerSvc) ListOperatorSharesIncreased(params *models.QueryDTParams) (*models.OperatorSharesIncreasedListResponse, error) {
	operatorSharesIncreasedList, total := h.operatorSharesIncreasedView.ListOperatorSharesIncreased(params.Page, params.PageSize, params.Order)
	return &models.OperatorSharesIncreasedListResponse{
		ListResponse: models.ListResponse{
			Current: params.Page,
			Size:    params.PageSize,
			Total:   total,
		},
		Records: operatorSharesIncreasedList,
	}, nil
}
