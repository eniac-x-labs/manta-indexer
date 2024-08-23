package service

import (
	"github.com/eniac-x-labs/manta-indexer/api/models"
	"github.com/eniac-x-labs/manta-indexer/database/event"
)

func (h HandlerSvc) GetWithdrawalCompleted(guid string) (*event.WithdrawalCompleted, error) {
	withdrawalCompleted, err := h.withdrawalCompletedView.GetWithdrawalCompleted(guid)
	if err != nil {
		return &event.WithdrawalCompleted{}, err
	}
	return withdrawalCompleted, err
}

func (h HandlerSvc) ListWithdrawalCompleted(params *models.QueryDTParams) (*models.WithdrawalCompletedListResponse, error) {
	withdrawalCompletedList, total := h.withdrawalCompletedView.ListWithdrawalCompleted(params.Page, params.PageSize, params.Order)
	return &models.WithdrawalCompletedListResponse{
		ListResponse: models.ListResponse{
			Current: params.Page,
			Size:    params.PageSize,
			Total:   total,
		},
		Records: withdrawalCompletedList,
	}, nil
}
