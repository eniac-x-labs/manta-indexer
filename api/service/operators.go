package service

import (
	"github.com/eniac-x-labs/manta-indexer/database/event"
	"strings"

	"github.com/eniac-x-labs/manta-indexer/api/models"
	"github.com/eniac-x-labs/manta-indexer/database/worker"
)

func (h HandlerSvc) Strategy(strategy string) (*event.Strategies, error) {
	addressToLower := strings.ToLower(strategy)
	strategies, err := h.strategiesView.QueryStrategies(addressToLower)
	if err != nil {
		return &event.Strategies{}, err
	}
	return strategies, err
}

func (h HandlerSvc) StrategyList(params *models.QueryListParams) (*models.StrategiesListResponse, error) {
	strategiesList, total := h.strategiesView.QueryStrategiesList(params.Page, params.PageSize, params.Order)
	return &models.StrategiesListResponse{
		ListResponse: models.ListResponse{
			Current: params.Page,
			Size:    params.PageSize,
			Total:   total,
		},
		Records: strategiesList,
	}, nil
}

func (h HandlerSvc) GetOperator(operator string) (*worker.Operators, error) {
	operatorAddress := strings.ToLower(operator)
	operatorDetail, err := h.operatorsView.GetOperator(operatorAddress)
	if err != nil {
		return &worker.Operators{}, err
	}
	return operatorDetail, err
}

func (h HandlerSvc) ListOperator(params *models.QueryListParams) (*models.OperatorListResponse, error) {
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

func (h HandlerSvc) RegisterOperator(operator string) (*event.OperatorRegistered, error) {
	addressToLower := strings.ToLower(operator)
	operatorRegistered, err := h.operatorRegisteredView.QueryOperatorRegistered(addressToLower)
	if err != nil {
		return &event.OperatorRegistered{}, err
	}
	return operatorRegistered, err
}

func (h HandlerSvc) RegisterOperatorList(params *models.QueryListParams) (*models.RegisterOperatorListResponse, error) {
	operatorRegisteredList, total := h.operatorRegisteredView.QueryOperatorRegisteredList(params.Page, params.PageSize, params.Order)
	return &models.RegisterOperatorListResponse{
		ListResponse: models.ListResponse{
			Current: params.Page,
			Size:    params.PageSize,
			Total:   total,
		},
		Records: operatorRegisteredList,
	}, nil
}

func (h HandlerSvc) ListOperatorNodeUrlUpdate(params *models.QueryAddressListParams) (*models.OperatorNodeUrlUpdateListResponse, error) {
	operatorNodeUrlUpdateList, total := h.operatorNodeUrlUpdateView.ListOperatorNodeUrlUpdate(params.Address, params.Page, params.PageSize, params.Order)
	return &models.OperatorNodeUrlUpdateListResponse{
		ListResponse: models.ListResponse{
			Current: params.Page,
			Size:    params.PageSize,
			Total:   total,
		},
		Records: operatorNodeUrlUpdateList,
	}, nil
}

func (h HandlerSvc) ListOperatorSharesIncreased(params *models.QueryAddressListParams) (*models.OperatorSharesIncreasedListResponse, error) {
	operatorSharesIncreasedList, total := h.operatorSharesIncreasedView.ListOperatorSharesIncreased(params.Address, params.Page, params.PageSize, params.Order)
	return &models.OperatorSharesIncreasedListResponse{
		ListResponse: models.ListResponse{
			Current: params.Page,
			Size:    params.PageSize,
			Total:   total,
		},
		Records: operatorSharesIncreasedList,
	}, nil
}

func (h HandlerSvc) ListOperatorSharesDecreased(params *models.QueryAddressListParams) (*models.OperatorSharesDecreasedListResponse, error) {
	operatorSharesDecreasedList, total := h.operatorSharesDecreasedView.ListOperatorSharesDecreased(params.Address, params.Page, params.PageSize, params.Order)
	return &models.OperatorSharesDecreasedListResponse{
		ListResponse: models.ListResponse{
			Current: params.Page,
			Size:    params.PageSize,
			Total:   total,
		},
		Records: operatorSharesDecreasedList,
	}, nil
}

func (h HandlerSvc) ListOperatorAndStakeReward(params *models.QueryAddressListParams) (*models.OperatorAndStakeRewardListResponse, error) {
	rewardList, total := h.operatorAndStakeRewardView.ListOperatorAndStakeReward(params.Address, params.Page, params.PageSize, params.Order)
	return &models.OperatorAndStakeRewardListResponse{
		ListResponse: models.ListResponse{
			Current: params.Page,
			Size:    params.PageSize,
			Total:   total,
		},
		Records: rewardList,
	}, nil
}

func (h HandlerSvc) ListOperatorClaimReward(params *models.QueryAddressListParams) (*models.OperatorClaimRewardListResponse, error) {
	rewardList, total := h.operatorClaimRewardView.ListOperatorClaimReward(params.Address, params.Page, params.PageSize, params.Order)
	return &models.OperatorClaimRewardListResponse{
		ListResponse: models.ListResponse{
			Current: params.Page,
			Size:    params.PageSize,
			Total:   total,
		},
		Records: rewardList,
	}, nil
}
