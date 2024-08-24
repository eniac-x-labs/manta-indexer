package service

import (
	"strconv"

	"github.com/ethereum/go-ethereum/log"

	"github.com/eniac-x-labs/manta-indexer/api/models"
	"github.com/eniac-x-labs/manta-indexer/database/event"
	"github.com/eniac-x-labs/manta-indexer/database/worker"
)

type Service interface {

	/*
	* ============== Strategy ==============
	 */
	Strategy(strategy string) (*event.Strategies, error)
	StrategyList(*models.QueryListParams) (*models.StrategiesListResponse, error)

	/*
	* ============== Operator ==============
	 */
	GetOperator(operator string) (*worker.Operators, error)
	ListOperator(*models.QueryListParams) (*models.OperatorListResponse, error)
	RegisterOperator(operator string) (*event.OperatorRegistered, error)
	RegisterOperatorList(*models.QueryListParams) (*models.RegisterOperatorListResponse, error)
	ListOperatorNodeUrlUpdate(*models.QueryAddressListParams) (*models.OperatorNodeUrlUpdateListResponse, error)
	ListOperatorSharesDecreased(*models.QueryAddressListParams) (*models.OperatorSharesDecreasedListResponse, error)
	ListOperatorSharesIncreased(*models.QueryAddressListParams) (*models.OperatorSharesIncreasedListResponse, error)
	ListOperatorAndStakeReward(*models.QueryAddressListParams) (*models.OperatorAndStakeRewardListResponse, error)
	ListOperatorClaimReward(params *models.QueryAddressListParams) (*models.OperatorClaimRewardListResponse, error)

	/*
	* ============== stakeholder ==============
	 */
	GetStakeHolder(staker string) (*worker.StakeHolder, error)
	ListStakeHolder(*models.QueryListParams) (*models.StakeHolderListResponse, error)
	ListStakerDepositStrategy(*models.QueryAddressListParams) (*models.StrategyDepositListResponse, error)
	ListStakerDelegated(*models.QueryAddressListParams) (*models.StakerDelegatedListResponse, error)
	ListStakerUndelegated(*models.QueryAddressListParams) (*models.StakerUndelegatedListResponse, error)
	ListStakerWithdrawalQueued(*models.QueryAddressListParams) (*models.WithdrawalQueuedListResponse, error)
	ListStakerWithdrawalCompleted(*models.QueryAddressListParams) (*models.WithdrawalCompletedListResponse, error)
	ListStakeHolderClaimReward(*models.QueryAddressListParams) (*models.StakeHolderClaimRewardListResponse, error)

	/*
	* ============== params check ==============
	 */
	QueryListParams(page string, pageSize string, order string) (*models.QueryListParams, error)
	QueryAddressListParams(address string, page string, pageSize string, order string) (*models.QueryAddressListParams, error)
}

type HandlerSvc struct {
	v                           *Validator
	strategiesView              event.StrategiesView
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
	strategiesView event.StrategiesView,
) Service {
	return &HandlerSvc{
		v:                           v,
		strategiesView:              strategiesView,
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

func (h HandlerSvc) QueryListParams(page string, pageSize string, order string) (*models.QueryListParams, error) {
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

	return &models.QueryListParams{
		Page:     pageVal,
		PageSize: pageSizeVal,
		Order:    orderBy,
	}, nil
}

func (h HandlerSvc) QueryAddressListParams(address string, page string, pageSize string, order string) (*models.QueryAddressListParams, error) {
	var paraAddress string
	if address == "0x00" {
		paraAddress = "0x00"
	} else {
		addr, err := h.v.ParseValidateAddress(address)
		if err != nil {
			log.Error("invalid address param", "address", address, "err", err)
			return nil, err
		}
		paraAddress = addr.String()
	}

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

	return &models.QueryAddressListParams{
		Address:  paraAddress,
		Page:     pageVal,
		PageSize: pageSizeVal,
		Order:    orderBy,
	}, nil
}
