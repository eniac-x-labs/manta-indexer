package models

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/eniac-x-labs/manta-indexer/database/event/operator"
	"github.com/eniac-x-labs/manta-indexer/database/event/staker"
	"github.com/eniac-x-labs/manta-indexer/database/event/strategies"
	"github.com/eniac-x-labs/manta-indexer/database/worker"
)

type QueryAddressListParams struct {
	Address  string
	Page     int
	PageSize int
	Order    string
}

type QueryListParams struct {
	Page     int
	PageSize int
	Order    string
}

type ListResponse struct {
	Current int    `json:"current"`
	Size    int    `json:"size"`
	Total   uint64 `json:"total"`
}

type RegisterOperatorResponse struct {
	BlockHash common.Hash    `json:"block_hash"`
	Number    *big.Int       `json:"number"`
	TxHash    common.Hash    `json:"tx_hash"`
	Operator  common.Address `json:"operator"`
	Timestamp uint64         `json:"timestamp"`
}

type RegisterOperatorListResponse struct {
	ListResponse
	Records []operator.OperatorRegistered `json:"records"`
}

type StrategiesListResponse struct {
	ListResponse
	Records []strategies.Strategies `json:"records"`
}

type OperatorNodeUrlUpdateListResponse struct {
	ListResponse
	Records []operator.OperatorNodeUrlUpdate `json:"records"`
}

type OperatorReceiveStakerDelegateListResponse struct {
	ListResponse
	Records []staker.StakerDelegated `json:"records"`
}

type StrategyDepositListResponse struct {
	ListResponse
	Records []staker.StrategyDeposit `json:"records"`
}

type StakeHolderListResponse struct {
	ListResponse
	Records []worker.StakeHolder `json:"records"`
}

type OperatorListResponse struct {
	ListResponse
	Records []worker.Operators `json:"records"`
}

type WithdrawalQueuedListResponse struct {
	ListResponse
	Records []staker.WithdrawalQueued `json:"records"`
}

type WithdrawalCompletedListResponse struct {
	ListResponse
	Records []staker.WithdrawalCompleted `json:"records"`
}

type StakerDelegatedListResponse struct {
	ListResponse
	Records []staker.StakerDelegated `json:"records"`
}

type StakerUndelegatedListResponse struct {
	ListResponse
	Records []staker.StakerUndelegated `json:"records"`
}

type StakeHolderClaimRewardListResponse struct {
	ListResponse
	Records []staker.StakeHolderClaimReward `json:"records"`
}

type OperatorSharesDecreasedListResponse struct {
	ListResponse
	Records []operator.OperatorSharesDecreased `json:"records"`
}

type OperatorSharesIncreasedListResponse struct {
	ListResponse
	Records []operator.OperatorSharesIncreased `json:"records"`
}

type OperatorAndStakeRewardListResponse struct {
	ListResponse
	Records []operator.OperatorAndStakeReward `json:"records"`
}

type OperatorClaimRewardListResponse struct {
	ListResponse
	Records []operator.OperatorClaimReward `json:"records"`
}
