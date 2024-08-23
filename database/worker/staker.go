package worker

import (
	"errors"

	"gorm.io/gorm"
	"math/big"

	"github.com/google/uuid"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

type StakeHolder struct {
	GUID            uuid.UUID      `gorm:"primaryKey"`
	Staker          common.Address `gorm:"serializer:bytes"`
	TotalMantaStake *big.Int       `gorm:"serializer:u256"`
	TotalReward     *big.Int       `gorm:"serializer:u256"`
	ClaimedAmount   *big.Int       `gorm:"serializer:u256"`
	Timestamp       uint64
}

func (StakeHolder) TableName() string {
	return "staker_holder"
}

type StakeHolderView interface {
}

type StakeHolderDB interface {
	StakeHolderView
	QueryAndUpdateStakeHolder(common.Address, StakeHolderType) error
	StoreStakerHolder([]StakeHolder) error
}

type stakeHolderDB struct {
	gorm *gorm.DB
}

func (sh *stakeHolderDB) QueryAndUpdateStakeHolder(stakeAddress common.Address, shType StakeHolderType) error {
	var stakeHolder StakeHolder
	result := sh.gorm.Where(&StakeHolder{Staker: stakeAddress}).Take(&stakeHolder)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			stHolder := StakeHolder{
				GUID:            uuid.New(),
				Staker:          stakeAddress,
				TotalMantaStake: shType.MantaStake,
				TotalReward:     shType.Reward,
				ClaimedAmount:   shType.ClaimedAmount,
				Timestamp:       shType.Timestamp,
			}
			err := sh.gorm.Create(stHolder).Error
			if err != nil {
				log.Error("create stake holder fail", "err", err)
			}
		}
		return result.Error
	}
	if shType.MantaStake != nil {
		stakeHolder.TotalMantaStake = new(big.Int).And(stakeHolder.TotalMantaStake, shType.MantaStake)
	}
	if shType.Reward != nil {
		stakeHolder.TotalReward = new(big.Int).And(stakeHolder.TotalReward, shType.Reward)
	}
	if shType.ClaimedAmount != nil {
		stakeHolder.ClaimedAmount = new(big.Int).And(stakeHolder.ClaimedAmount, shType.ClaimedAmount)
	}
	err := sh.gorm.Save(stakeHolder).Error
	if err != nil {
		log.Error("Update node url fail", "err", err)
		return err
	}
	return nil
}

func NewStakeHolderDB(db *gorm.DB) StakeHolderDB {
	return &stakeHolderDB{gorm: db}
}

func (sh *stakeHolderDB) StoreStakerHolder(stakers []StakeHolder) error {
	result := sh.gorm.CreateInBatches(&stakers, len(stakers))
	return result.Error
}
