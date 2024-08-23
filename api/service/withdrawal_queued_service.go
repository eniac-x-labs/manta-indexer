package service

import (
	"github.com/eniac-x-labs/manta-indexer/api/models"
	"github.com/eniac-x-labs/manta-indexer/database/event"
)

func (h HandlerSvc) GetWithdrawalQueued(guid string) (*event.WithdrawalQueued, error) {
	withdrawalQueued, err := h.withdrawalQueuedView.GetWithdrawalQueued(guid)
	if err != nil {
		return &event.WithdrawalQueued{}, err
	}
	return withdrawalQueued, err
}

func (h HandlerSvc) ListWithdrawalQueued(params *models.QueryDTParams) (*models.WithdrawalQueuedListResponse, error) {
	withdrawalQueuedList, total := h.withdrawalQueuedView.ListWithdrawalQueued(params.Page, params.PageSize, params.Order)
	return &models.WithdrawalQueuedListResponse{
		ListResponse: models.ListResponse{
			Current: params.Page,
			Size:    params.PageSize,
			Total:   total,
		},
		Records: withdrawalQueuedList,
	}, nil
}
