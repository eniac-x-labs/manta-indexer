package worker

import (
	"math/big"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum/common"
)

type Operator struct {
	GUID                     uuid.UUID      `gorm:"primaryKey"`
	BlockHash                common.Hash    `gorm:"serializer:bytes"`
	Number                   *big.Int       `gorm:"serializer:u256"`
	TxHash                   common.Hash    `gorm:"serializer:bytes"`
	Operator                 common.Address `gorm:"serializer:bytes"`
	Socket                   string         `json:"socket"`
	EarningsReceiver         common.Address `gorm:"serializer:bytes"`
	DelegationApprover       common.Address `gorm:"serializer:bytes"`
	StakerOptoutWindowBlocks *big.Int       `gorm:"serializer:u256"`
	TotalMantaStake          *big.Int       `gorm:"serializer:u256"`
	TotalStakeReward         *big.Int       `gorm:"serializer:u256"`
	RateReturn               string         `json:"rate_return"`
	Status                   uint8
	Timestamp                uint64
}

type OperatorsView interface {
}

type OperatorDB interface {
	OperatorsView
	StoreOperators([]Operator) error
}

type contractEventDB struct {
	gorm *gorm.DB
}

func NewOperatorsDB(db *gorm.DB) OperatorDB {
	return &contractEventDB{gorm: db}
}

func (db *contractEventDB) StoreOperators(events []Operator) error {
	result := db.gorm.CreateInBatches(&events, len(events))
	return result.Error
}
