package service

import (
	"strconv"
	"strings"

	"github.com/eniac-x-labs/manta-indexer/api/models"
	"github.com/eniac-x-labs/manta-indexer/database/event"
	"github.com/eniac-x-labs/manta-indexer/database/worker"
)

type Service interface {
	RegisterOperatorList(*models.QueryDTParams) (*models.RegisterOperatorListResponse, error)
	RegisterOperator(operator string) (*event.OperatorRegistered, error)

	GetOperatorNodeUrlUpdate(operator string) (*event.OperatorNodeUrlUpdate, error)
	ListOperatorNodeUrlUpdate(*models.QueryDTParams) (*models.OperatorNodeUrlUpdateListResponse, error)

	GetStrategyDeposit(staker string) (*event.StrategyDeposit, error)
	ListStrategyDeposit(*models.QueryDTParams) (*models.StrategyDepositListResponse, error)

	GetStakeHolder(staker string) (*worker.StakeHolder, error)
	ListStakeHolder(*models.QueryDTParams) (*models.StakeHolderListResponse, error)

	GetOperator(operator string) (*worker.Operators, error)
	ListOperator(*models.QueryDTParams) (*models.OperatorListResponse, error)

	GetWithdrawalQueued(guid string) (*event.WithdrawalQueued, error)
	ListWithdrawalQueued(*models.QueryDTParams) (*models.WithdrawalQueuedListResponse, error)

	GetWithdrawalCompleted(guid string) (*event.WithdrawalCompleted, error)
	ListWithdrawalCompleted(*models.QueryDTParams) (*models.WithdrawalCompletedListResponse, error)

	GetStakerDelegated(guid string) (*event.StakerDelegated, error)
	ListStakerDelegated(*models.QueryDTParams) (*models.StakerDelegatedListResponse, error)

	GetStakerUndelegated(guid string) (*event.StakerUndelegated, error)
	ListStakerUndelegated(*models.QueryDTParams) (*models.StakerUndelegatedListResponse, error)

	GetStakeHolderClaimReward(guid string) (*event.StakeHolderClaimReward, error)
	ListStakeHolderClaimReward(*models.QueryDTParams) (*models.StakeHolderClaimRewardListResponse, error)

	GetOperatorSharesDecreased(guid string) (*event.OperatorSharesDecreased, error)
	ListOperatorSharesDecreased(*models.QueryDTParams) (*models.OperatorSharesDecreasedListResponse, error)

	GetOperatorSharesIncreased(guid string) (*event.OperatorSharesIncreased, error)
	ListOperatorSharesIncreased(*models.QueryDTParams) (*models.OperatorSharesIncreasedListResponse, error)

	GetOperatorAndStakeReward(guid string) (*event.OperatorAndStakeReward, error)
	ListOperatorAndStakeReward(*models.QueryDTParams) (*models.OperatorAndStakeRewardListResponse, error)

	GetOperatorClaimReward(guid string) (*event.OperatorClaimReward, error)
	ListOperatorClaimReward(*models.QueryDTParams) (*models.OperatorClaimRewardListResponse, error)

	QueryDTListParams(page string, pageSize string, order string) (*models.QueryDTParams, error)
}

type HandlerSvc struct {
	v                           *Validator
	operatorRegisteredView      event.OperatorRegisteredView
	operatorNodeUrlUpdateView   event.OperatorNodeUrlUpdateView
	operatorsView               worker.OperatorsView
	stakeHolderView             worker.StakeHolderView
	strategyDepositView         event.StrategyDepositView
	withdrawalQueuedView        event.WithdrawalQueuedView
	withdrawalCompletedView     event.WithdrawalCompletedView
	stakerDelegatedView         event.StakerDelegatedView
	stakerUndelegatedView       event.StakerUndelegatedView
	stakeHolderClaimRewardView  event.StakeHolderClaimRewardView
	operatorSharesDecreasedView event.OperatorSharesDecreasedView
	operatorSharesIncreasedView event.OperatorSharesIncreasedView
	operatorAndStakeRewardView  event.OperatorAndStakeRewardView
	operatorClaimRewardView     event.OperatorClaimRewardView
}

func New(v *Validator,
	rgv event.OperatorRegisteredView,
	onuu event.OperatorNodeUrlUpdateView,
	operatorsView worker.OperatorsView,
	stakeHolderView worker.StakeHolderView,
	strategyDepositView event.StrategyDepositView,
	withdrawalQueuedView event.WithdrawalQueuedView,
	withdrawalCompletedView event.WithdrawalCompletedView,
	stakerDelegatedView event.StakerDelegatedView,
	stakerUndelegatedView event.StakerUndelegatedView,
	stakeHolderClaimRewardView event.StakeHolderClaimRewardView,
	operatorSharesDecreasedView event.OperatorSharesDecreasedView,
	operatorSharesIncreasedView event.OperatorSharesIncreasedView,
	operatorAndStakeRewardView event.OperatorAndStakeRewardView,
	operatorClaimRewardView event.OperatorClaimRewardView,

) Service {
	return &HandlerSvc{
		v:                           v,
		operatorRegisteredView:      rgv,
		operatorNodeUrlUpdateView:   onuu,
		operatorsView:               operatorsView,
		stakeHolderView:             stakeHolderView,
		strategyDepositView:         strategyDepositView,
		withdrawalQueuedView:        withdrawalQueuedView,
		withdrawalCompletedView:     withdrawalCompletedView,
		stakerDelegatedView:         stakerDelegatedView,
		stakerUndelegatedView:       stakerUndelegatedView,
		stakeHolderClaimRewardView:  stakeHolderClaimRewardView,
		operatorSharesDecreasedView: operatorSharesDecreasedView,
		operatorSharesIncreasedView: operatorSharesIncreasedView,
		operatorAndStakeRewardView:  operatorAndStakeRewardView,
		operatorClaimRewardView:     operatorClaimRewardView,
	}
}

func (h HandlerSvc) RegisterOperator(operator string) (*event.OperatorRegistered, error) {
	addressToLower := strings.ToLower(operator)
	operatorRegistered, err := h.operatorRegisteredView.QueryOperatorRegistered(addressToLower)
	if err != nil {
		return &event.OperatorRegistered{}, err
	}
	return operatorRegistered, err
}

func (h HandlerSvc) RegisterOperatorList(params *models.QueryDTParams) (*models.RegisterOperatorListResponse, error) {
	operatorRegisteredList, total := h.operatorRegisteredView.QueryOperatorRegisteredList(params.Page, params.PageSize, params.Order)
	return &models.RegisterOperatorListResponse{
		Current: params.Page,
		Size:    params.PageSize,
		Total:   total,
		Records: operatorRegisteredList,
	}, nil
}

func (h HandlerSvc) QueryDTListParams(page string, pageSize string, order string) (*models.QueryDTParams, error) {
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return nil, err
	}
	pageVal := h.v.ValidatePage(pageInt)

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		return nil, err
	}
	pageSizeVal := h.v.ValidatePageSize(pageSizeInt)
	orderBy := h.v.ValidateOrder(order)

	return &models.QueryDTParams{
		Page:     pageVal,
		PageSize: pageSizeVal,
		Order:    orderBy,
	}, nil
}
