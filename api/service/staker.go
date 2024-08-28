package service

import (
	"strings"

	"github.com/eniac-x-labs/manta-indexer/api/models"
	"github.com/eniac-x-labs/manta-indexer/database/worker"
)

func (h HandlerSvc) GetStakeHolder(staker string) (*worker.StakeStrategy, error) {
	stakerAddress := strings.ToLower(staker)
	stakeStrategyList, err := h.stakeStrategyView.GetStakeStrategy(stakerAddress)
	if err != nil {
		return &worker.StakeStrategy{}, err
	}
	return stakeStrategyList, err
}

func (h HandlerSvc) ListStakeHolder(params *models.QueryAddressListParams) (*models.StakeHolderListResponse, error) {
	stakerAddress := strings.ToLower(params.Address)
	stakeHolderList, total := h.stakeStrategyView.ListStakeStrategy(stakerAddress, params.Page, params.PageSize, params.Order)
	return &models.StakeHolderListResponse{
		ListResponse: models.ListResponse{
			Current: params.Page,
			Size:    params.PageSize,
			Total:   total,
		},
		Records: stakeHolderList,
	}, nil
}

func (h HandlerSvc) ListStakeOperator(operatorAddress string, params *models.QueryAddressListParams) (*models.StakeOperatorListResponse, error) {
	operatorAddr := strings.ToLower(operatorAddress)
	stakeAddr := strings.ToLower(params.Address)
	stakerOperatorHolderList, total := h.stakerOperatorView.ListStakerOperator(operatorAddr, stakeAddr, params.Page, params.PageSize, params.Order)
	return &models.StakeOperatorListResponse{
		ListResponse: models.ListResponse{
			Current: params.Page,
			Size:    params.PageSize,
			Total:   total,
		},
		Records: stakerOperatorHolderList,
	}, nil
}

func (h HandlerSvc) ListStakerDepositStrategy(params *models.QueryAddressListParams) (*models.StrategyDepositListResponse, error) {
	stakerAddress := strings.ToLower(params.Address)
	strategyDepositList, total := h.strategyDepositView.ListStrategyDeposit(stakerAddress, params.Page, params.PageSize, params.Order)
	return &models.StrategyDepositListResponse{
		ListResponse: models.ListResponse{
			Current: params.Page,
			Size:    params.PageSize,
			Total:   total,
		},
		Records: strategyDepositList,
	}, nil
}

func (h HandlerSvc) ListStakerDelegated(params *models.QueryAddressListParams) (*models.StakerDelegatedListResponse, error) {
	stakerAddress := strings.ToLower(params.Address)
	stakerDelegatedList, total := h.stakerDelegatedView.ListStakerDelegated(stakerAddress, params.Page, params.PageSize, params.Order)
	return &models.StakerDelegatedListResponse{
		ListResponse: models.ListResponse{
			Current: params.Page,
			Size:    params.PageSize,
			Total:   total,
		},
		Records: stakerDelegatedList,
	}, nil
}

func (h HandlerSvc) ListStakerUndelegated(params *models.QueryAddressListParams) (*models.StakerUndelegatedListResponse, error) {
	stakerAddress := strings.ToLower(params.Address)
	stakerUndelegatedList, total := h.stakerUndelegatedView.ListStakerUndelegated(stakerAddress, params.Page, params.PageSize, params.Order)
	return &models.StakerUndelegatedListResponse{
		ListResponse: models.ListResponse{
			Current: params.Page,
			Size:    params.PageSize,
			Total:   total,
		},
		Records: stakerUndelegatedList,
	}, nil
}

func (h HandlerSvc) ListStakerWithdrawalQueued(params *models.QueryAddressListParams) (*models.WithdrawalQueuedListResponse, error) {
	stakerAddress := strings.ToLower(params.Address)
	withdrawalQueuedList, total := h.withdrawalQueuedView.ListWithdrawalQueued(stakerAddress, params.Page, params.PageSize, params.Order)
	return &models.WithdrawalQueuedListResponse{
		ListResponse: models.ListResponse{
			Current: params.Page,
			Size:    params.PageSize,
			Total:   total,
		},
		Records: withdrawalQueuedList,
	}, nil
}

func (h HandlerSvc) ListStakerWithdrawalCompleted(params *models.QueryAddressListParams) (*models.WithdrawalCompletedListResponse, error) {
	stakerAddress := strings.ToLower(params.Address)
	withdrawalCompletedList, total := h.withdrawalCompletedView.ListWithdrawalCompleted(stakerAddress, params.Page, params.PageSize, params.Order)
	return &models.WithdrawalCompletedListResponse{
		ListResponse: models.ListResponse{
			Current: params.Page,
			Size:    params.PageSize,
			Total:   total,
		},
		Records: withdrawalCompletedList,
	}, nil
}

func (h HandlerSvc) ListStakeHolderClaimReward(params *models.QueryAddressListParams) (*models.StakeHolderClaimRewardListResponse, error) {
	stakerAddress := strings.ToLower(params.Address)
	stakeHolderClaimRewardList, total := h.stakeHolderClaimRewardView.ListStakeHolderClaimReward(stakerAddress, params.Page, params.PageSize, params.Order)
	return &models.StakeHolderClaimRewardListResponse{
		ListResponse: models.ListResponse{
			Current: params.Page,
			Size:    params.PageSize,
			Total:   total,
		},
		Records: stakeHolderClaimRewardList,
	}, nil
}
