package service

import (
	"github.com/eniac-x-labs/manta-indexer/api/models"
	"github.com/eniac-x-labs/manta-indexer/database/event"
)

func (h HandlerSvc) GetOperatorSharesDecreased(guid string) (*event.OperatorSharesDecreased, error) {
	operatorSharesDecreased, err := h.operatorSharesDecreasedView.GetOperatorSharesDecreased(guid)
	if err != nil {
		return &event.OperatorSharesDecreased{}, err
	}
	return operatorSharesDecreased, err
}

func (h HandlerSvc) ListOperatorSharesDecreased(params *models.QueryDTParams) (*models.OperatorSharesDecreasedListResponse, error) {
	operatorSharesDecreasedList, total := h.operatorSharesDecreasedView.ListOperatorSharesDecreased(params.Page, params.PageSize, params.Order)
	return &models.OperatorSharesDecreasedListResponse{
		ListResponse: models.ListResponse{
			Current: params.Page,
			Size:    params.PageSize,
			Total:   total,
		},
		Records: operatorSharesDecreasedList,
	}, nil
}
