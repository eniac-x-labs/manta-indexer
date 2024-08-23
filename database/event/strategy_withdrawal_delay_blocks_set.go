package event

import (
	"math/big"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum/common"

	_ "github.com/eniac-x-labs/manta-indexer/database/utils/serializers"
)

type StrategyWithdrawalDelayBlocksSet struct {
	GUID          uuid.UUID      `gorm:"primaryKey" json:"guid"`
	BlockHash     common.Hash    `json:"block_hash" gorm:"serializer:bytes"`
	Number        *big.Int       `json:"number" gorm:"serializer:u256"`
	TxHash        common.Hash    `json:"tx_hash" gorm:"serializer:bytes"`
	Strategy      common.Address `json:"strategy" gorm:"serializer:bytes"`
	PreviousValue *big.Int       `json:"previous_value" gorm:"serializer:u256"`
	NewValue      *big.Int       `json:"new_value" gorm:"serializer:u256"`
	IsHandle      uint8          `json:"is_handle"`
	Timestamp     uint64         `json:"timestamp"`
}

type StrategyWithdrawalDelayBlocksSetView interface {
	QueryStrategyWithdrawalDelayBlocksSetList(page int, pageSize int, order string) ([]StrategyWithdrawalDelayBlocksSet, uint64)
}

type StrategyWithdrawalDelayBlocksSetDB interface {
	StrategyWithdrawalDelayBlocksSetView
	StoreStrategyWithdrawalDelayBlocksSet([]StrategyWithdrawalDelayBlocksSet) error
}

type strategyWithdrawalDelayBlocksSetDB struct {
	gorm *gorm.DB
}

func (db strategyWithdrawalDelayBlocksSetDB) QueryStrategyWithdrawalDelayBlocksSetList(page int, pageSize int, order string) ([]StrategyWithdrawalDelayBlocksSet, uint64) {
	panic("implement me")
}

func (db strategyWithdrawalDelayBlocksSetDB) StoreStrategyWithdrawalDelayBlocksSet(strategyWithdrawalDelayBlocksSetList []StrategyWithdrawalDelayBlocksSet) error {
	result := db.gorm.Table("strategy_withdrawal_delay_blocks_set").CreateInBatches(&strategyWithdrawalDelayBlocksSetList, len(strategyWithdrawalDelayBlocksSetList))
	return result.Error
}

func NewStrategyWithdrawalDelayBlocksSetDB(db *gorm.DB) StrategyWithdrawalDelayBlocksSetDB {
	return &strategyWithdrawalDelayBlocksSetDB{gorm: db}
}
