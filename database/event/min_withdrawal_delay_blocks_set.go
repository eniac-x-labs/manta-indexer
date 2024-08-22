package event

import (
	"math/big"

	_ "github.com/eniac-x-labs/manta-indexer/database/utils/serializers"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MinWithdrawalDelayBlocksSet struct {
	GUID          uuid.UUID   `gorm:"primaryKey" json:"guid"`
	BlockHash     common.Hash `json:"block_hash" gorm:"serializer:bytes"`
	Number        *big.Int    `json:"number" gorm:"serializer:u256"`
	TxHash        common.Hash `json:"tx_hash" gorm:"serializer:bytes"`
	PreviousValue *big.Int    `json:"previous_value" gorm:"serializer:u256"`
	NewValue      *big.Int    `json:"new_value" gorm:"serializer:u256"`
	IsHandle      uint8       `json:"is_handle"`
	Timestamp     uint64      `json:"timestamp"`
}

type MinWithdrawalDelayBlocksSetView interface {
	QueryMinWithdrawalDelayBlocksSetList(page int, pageSize int, order string) ([]MinWithdrawalDelayBlocksSet, uint64)
}

type MinWithdrawalDelayBlocksSetDB interface {
	MinWithdrawalDelayBlocksSetView
	StoreMinWithdrawalDelayBlocksSet([]MinWithdrawalDelayBlocksSet) error
}

type minWithdrawalDelayBlocksSetDB struct {
	gorm *gorm.DB
}

func (db minWithdrawalDelayBlocksSetDB) QueryMinWithdrawalDelayBlocksSetList(page int, pageSize int, order string) ([]MinWithdrawalDelayBlocksSet, uint64) {
	panic("implement me")
}

func (db minWithdrawalDelayBlocksSetDB) StoreMinWithdrawalDelayBlocksSet(minWithdrawalDelayBlocksSetList []MinWithdrawalDelayBlocksSet) error {
	result := db.gorm.CreateInBatches(&minWithdrawalDelayBlocksSetList, len(minWithdrawalDelayBlocksSetList))
	return result.Error
}

func NewMinWithdrawalDelayBlocksSetDB(db *gorm.DB) MinWithdrawalDelayBlocksSetDB {
	return &minWithdrawalDelayBlocksSetDB{gorm: db}
}
