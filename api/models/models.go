package models

import (
	"github.com/eniac-x-labs/manta-indexer/database/event"
	"github.com/eniac-x-labs/manta-indexer/database/worker"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type QueryAddressParams struct {
	Address string
}

type QueryDTParams struct {
	Page     int
	PageSize int
	Order    string
}

type ListResponse struct {
	Current int    `json:"Current"`
	Size    int    `json:"Size"`
	Total   uint64 `json:"Total"`
}

type RegisterOperatorResponse struct {
	BlockHash common.Hash    `json:"block_hash"`
	Number    *big.Int       `json:"number"`
	TxHash    common.Hash    `json:"tx_hash"`
	Operator  common.Address `json:"operator"`
	Timestamp uint64         `json:"timestamp"`
}

type RegisterOperatorListResponse struct {
	Current int    `json:"Current"`
	Size    int    `json:"Size"`
	Total   uint64 `json:"Total"`
	Records []event.OperatorRegistered
}

type OperatorNodeUrlUpdateListResponse struct {
	ListResponse
	Records []event.OperatorNodeUrlUpdate
}

type StrategyDepositListResponse struct {
	ListResponse
	Records []event.StrategyDeposit
}

type StakeHolderListResponse struct {
	ListResponse
	Records []worker.StakeHolder
}

type OperatorListResponse struct {
	ListResponse
	Records []worker.Operators
}

type WithdrawalQueuedListResponse struct {
	ListResponse
	Records []event.WithdrawalQueued
}

type WithdrawalCompletedListResponse struct {
	ListResponse
	Records []event.WithdrawalCompleted
}

type StakerDelegatedListResponse struct {
	ListResponse
	Records []event.StakerDelegated
}

type StakerUndelegatedListResponse struct {
	ListResponse
	Records []event.StakerUndelegated
}

type StakeHolderClaimRewardListResponse struct {
	ListResponse
	Records []event.StakeHolderClaimReward
}

type OperatorSharesDecreasedListResponse struct {
	ListResponse
	Records []event.OperatorSharesDecreased
}

type OperatorSharesIncreasedListResponse struct {
	ListResponse
	Records []event.OperatorSharesIncreased
}

type OperatorAndStakeRewardListResponse struct {
	ListResponse
	Records []event.OperatorAndStakeReward
}

type OperatorClaimRewardListResponse struct {
	ListResponse
	Records []event.OperatorClaimReward
}
