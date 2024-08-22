package worker

import (
	"gorm.io/gorm"
	"math/big"

	"github.com/google/uuid"

	"github.com/ethereum/go-ethereum/common"
)

type TotalOperator struct {
	GUID            uuid.UUID      `gorm:"primaryKey"`
	Staker          common.Address `gorm:"serializer:bytes"`
	TotalMantaStake *big.Int       `gorm:"serializer:u256"`
	TotalReward     *big.Int       `gorm:"serializer:u256"`
	ClaimedAmount   *big.Int       `gorm:"serializer:u256"`
	Timestamp       uint64
}

type TotalOperatorView interface {
}

type TotalOperatorDB interface {
	TotalOperatorView
	StoreTotalOperator([]TotalOperator) error
}

type totalOperatorDB struct {
	gorm *gorm.DB
}

func (to totalOperatorDB) StoreTotalOperator(operators []TotalOperator) error {
	result := to.gorm.CreateInBatches(&operators, len(operators))
	return result.Error
}

func NewTotalOperatorDB(db *gorm.DB) TotalOperatorDB {
	return &totalOperatorDB{gorm: db}
}
