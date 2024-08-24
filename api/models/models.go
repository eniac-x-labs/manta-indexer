package models

import (
	"math/big"

	"github.com/eniac-x-labs/manta-indexer/database/event"
	"github.com/eniac-x-labs/manta-indexer/database/worker"
	"github.com/ethereum/go-ethereum/common"
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
	Records []event.OperatorRegistered `json:"records"`
}

type StrategiesListResponse struct {
	ListResponse
	Records []event.Strategies `json:"records"`
}

type OperatorNodeUrlUpdateListResponse struct {
	ListResponse
	Records []event.OperatorNodeUrlUpdate `json:"records"`
}

type StrategyDepositListResponse struct {
	ListResponse
	Records []event.StrategyDeposit `json:"records"`
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
	Records []event.WithdrawalQueued `json:"records"`
}

type WithdrawalCompletedListResponse struct {
	ListResponse
	Records []event.WithdrawalCompleted `json:"records"`
}

type StakerDelegatedListResponse struct {
	ListResponse
	Records []event.StakerDelegated `json:"records"`
}

type StakerUndelegatedListResponse struct {
	ListResponse
	Records []event.StakerUndelegated `json:"records"`
}

type StakeHolderClaimRewardListResponse struct {
	ListResponse
	Records []event.StakeHolderClaimReward `json:"records"`
}

type OperatorSharesDecreasedListResponse struct {
	ListResponse
	Records []event.OperatorSharesDecreased `json:"records"`
}

type OperatorSharesIncreasedListResponse struct {
	ListResponse
	Records []event.OperatorSharesIncreased `json:"records"`
}

type OperatorAndStakeRewardListResponse struct {
	ListResponse
	Records []event.OperatorAndStakeReward `json:"records"`
}

type OperatorClaimRewardListResponse struct {
	ListResponse
	Records []event.OperatorClaimReward `json:"records"`
}
