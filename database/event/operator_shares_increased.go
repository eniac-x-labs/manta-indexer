package event

import (
	"math/big"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum/common"

	_ "github.com/eniac-x-labs/manta-indexer/database/utils/serializers"
)

type OperatorSharesIncreased struct {
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

type OperatorSharesIncreasedView interface {
	QueryOperatorSharesIncreasedList(page int, pageSize int, order string) ([]OperatorSharesIncreased, uint64)
}

type OperatorSharesIncreasedDB interface {
	OperatorSharesIncreasedView
	StoreOperatorSharesIncreased([]OperatorSharesIncreased) error
}

type operatorSharesIncreasedDB struct {
	gorm *gorm.DB
}

func (db operatorSharesIncreasedDB) QueryOperatorSharesIncreasedList(page int, pageSize int, order string) ([]OperatorSharesIncreased, uint64) {
	panic("implement me")
}

func (db operatorSharesIncreasedDB) StoreOperatorSharesIncreased(operatorSharesIncreasedList []OperatorSharesIncreased) error {
	result := db.gorm.CreateInBatches(&operatorSharesIncreasedList, len(operatorSharesIncreasedList))
	return result.Error
}

func NewOperatorSharesIncreasedDB(db *gorm.DB) OperatorSharesIncreasedDB {
	return &operatorSharesIncreasedDB{gorm: db}
}
