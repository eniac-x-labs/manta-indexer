package worker

import (
	"gorm.io/gorm"
	"math/big"

	"github.com/google/uuid"

	"github.com/ethereum/go-ethereum/common"
)

type StakeHolder struct {
	GUID            uuid.UUID      `gorm:"primaryKey"`
	Staker          common.Address `gorm:"serializer:bytes"`
	TotalMantaStake *big.Int       `gorm:"serializer:u256"`
	TotalReward     *big.Int       `gorm:"serializer:u256"`
	ClaimedAmount   *big.Int       `gorm:"serializer:u256"`
	Timestamp       uint64
}

type StakeHolderView interface {
}

type StakeHolderDB interface {
	StakeHolderView
	StoreStakerHolder([]StakeHolder) error
}

type stakeHolderDB struct {
	gorm *gorm.DB
}

func NewStakeHolderDB(db *gorm.DB) StakeHolderDB {
	return &stakeHolderDB{gorm: db}
}

func (db *stakeHolderDB) StoreStakerHolder(stakers []StakeHolder) error {
	result := db.gorm.CreateInBatches(&stakers, len(stakers))
	return result.Error
}
