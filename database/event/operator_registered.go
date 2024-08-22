package event

import (
	"math/big"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum/common"

	_ "github.com/eniac-x-labs/manta-indexer/database/utils/serializers"
)

type OperatorRegistered struct {
	GUID                     uuid.UUID      `gorm:"primaryKey" json:"guid"`
	BlockHash                common.Hash    `json:"block_hash" gorm:"serializer:bytes"`
	Number                   *big.Int       `json:"number" gorm:"serializer:u256"`
	TxHash                   common.Hash    `json:"tx_hash" gorm:"serializer:bytes"`
	Operator                 common.Address `json:"operator" gorm:"serializer:bytes"`
	EarningsReceiver         common.Address `json:"earnings_receiver" gorm:"serializer:bytes"`
	DelegationApprover       common.Address `json:"delegation_approver" gorm:"serializer:bytes"`
	StakerOptoutWindowBlocks *big.Int       `json:"staker_optout_window_blocks" gorm:"serializer:u256"`
	IsHandle                 uint8          `json:"is_handle"`
	Timestamp                uint64         `json:"timestamp"`
}

type OperatorRegisteredView interface {
	QueryOperatorRegisteredList(page int, pageSize int, order string) ([]OperatorRegistered, uint64)
}

type OperatorRegisteredDB interface {
	OperatorRegisteredView
	StoreOperatorRegistered([]OperatorRegistered) error
}

type operatorRegisteredDB struct {
	gorm *gorm.DB
}

func (db operatorRegisteredDB) QueryOperatorRegisteredList(page int, pageSize int, order string) ([]OperatorRegistered, uint64) {
	panic("implement me")
}

func (db operatorRegisteredDB) StoreOperatorRegistered(operatorRegisteredList []OperatorRegistered) error {
	result := db.gorm.CreateInBatches(&operatorRegisteredList, len(operatorRegisteredList))
	return result.Error
}

func NewOperatorRegisteredDB(db *gorm.DB) OperatorRegisteredDB {
	return &operatorRegisteredDB{gorm: db}
}
