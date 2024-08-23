package models

import (
	"github.com/eniac-x-labs/manta-indexer/database/event"
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
