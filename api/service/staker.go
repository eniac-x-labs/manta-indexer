package service

import (
	"strings"

	"github.com/eniac-x-labs/manta-indexer/api/models"
	"github.com/eniac-x-labs/manta-indexer/database/worker"
)

func (h HandlerSvc) GetStakeHolder(staker string) (*worker.StakeHolder, error) {
	stakerAddress := strings.ToLower(staker)
	stakeHolder, err := h.stakeHolderView.GetStakeHolder(stakerAddress)
	if err != nil {
		return &worker.StakeHolder{}, err
	}
	return stakeHolder, err
}

func (h HandlerSvc) ListStakeHolder(params *models.QueryListParams) (*models.StakeHolderListResponse, error) {
	stakeHolderList, total := h.stakeHolderView.ListStakeHolder(params.Page, params.PageSize, params.Order)
	return &models.StakeHolderListResponse{
		ListResponse: models.ListResponse{
			Current: params.Page,
			Size:    params.PageSize,
			Total:   total,
		},
		Records: stakeHolderList,
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
