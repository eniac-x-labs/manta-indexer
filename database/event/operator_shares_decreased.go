package event

import (
	"math/big"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum/common"

	_ "github.com/eniac-x-labs/manta-indexer/database/utils/serializers"
)

type OperatorSharesDecreased struct {
	GUID      uuid.UUID      `gorm:"primaryKey" json:"guid"`
	BlockHash common.Hash    `json:"block_hash" gorm:"serializer:bytes"`
	Number    *big.Int       `json:"number" gorm:"serializer:u256"`
	TxHash    common.Hash    `json:"tx_hash" gorm:"serializer:bytes"`
	Operator  common.Address `json:"operator" gorm:"serializer:bytes"`
	Staker    common.Address `json:"staker" gorm:"serializer:bytes"`
	Strategy  common.Address `json:"strategy" gorm:"serializer:bytes"`
	Shares    *big.Int       `json:"shares" gorm:"serializer:u256"`
	IsHandle  uint8          `json:"is_handle"`
	Timestamp uint64         `json:"timestamp"`
}

type OperatorSharesDecreasedView interface {
	QueryOperatorSharesDecreasedList(page int, pageSize int, order string) ([]OperatorSharesDecreased, uint64)
}

type OperatorSharesDecreasedDB interface {
	OperatorSharesDecreasedView
	StoreOperatorSharesDecreased([]OperatorSharesDecreased) error
}

type operatorSharesDecreasedDB struct {
	gorm *gorm.DB
}

func (db operatorSharesDecreasedDB) QueryOperatorSharesDecreasedList(page int, pageSize int, order string) ([]OperatorSharesDecreased, uint64) {
	panic("implement me")
}

func (db operatorSharesDecreasedDB) StoreOperatorSharesDecreased(operatorSharesDecreasedList []OperatorSharesDecreased) error {
	result := db.gorm.CreateInBatches(&operatorSharesDecreasedList, len(operatorSharesDecreasedList))
	return result.Error
}

func NewOperatorSharesDecreasedDB(db *gorm.DB) OperatorSharesDecreasedDB {
	return &operatorSharesDecreasedDB{gorm: db}
}
