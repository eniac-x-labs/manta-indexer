package service

import (
	"github.com/eniac-x-labs/manta-indexer/api/models"
	"github.com/eniac-x-labs/manta-indexer/database/event"
	"strings"
)

func (h HandlerSvc) GetStrategyDeposit(staker string) (*event.StrategyDeposit, error) {
	addressToLower := strings.ToLower(staker)
	strategyDeposit, err := h.strategyDepositView.GetStrategyDeposit(addressToLower)
	if err != nil {
		return &event.StrategyDeposit{}, err
	}
	return strategyDeposit, err
}

func (h HandlerSvc) ListStrategyDeposit(params *models.QueryDTParams) (*models.StrategyDepositListResponse, error) {
	strategyDepositList, total := h.strategyDepositView.ListStrategyDeposit(params.Page, params.PageSize, params.Order)
	return &models.StrategyDepositListResponse{
		ListResponse: models.ListResponse{
			Current: params.Page,
			Size:    params.PageSize,
			Total:   total,
		},
		Records: strategyDepositList,
	}, nil
}
