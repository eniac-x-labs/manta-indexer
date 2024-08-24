package worker

import (
	"errors"
	"strings"

	"gorm.io/gorm"
	"math/big"

	"github.com/google/uuid"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

type StakeHolder struct {
	GUID            uuid.UUID      `gorm:"primaryKey" json:"guid"`
	Staker          common.Address `gorm:"serializer:bytes" json:"staker"`
	Strategy        common.Address `json:"strategy" gorm:"serializer:bytes"`
	TotalMantaStake *big.Int       `gorm:"serializer:u256" json:"total_manta_stake"`
	TotalReward     *big.Int       `gorm:"serializer:u256" json:"total_reward"`
	ClaimedAmount   *big.Int       `gorm:"serializer:u256" json:"claimed_amount"`
	Timestamp       uint64
}

func (StakeHolder) TableName() string {
	return "staker_holder"
}

type StakeHolderView interface {
	GetStakeHolder(string) (*StakeHolder, error)
	ListStakeHolder(address string, page int, pageSize int, order string) ([]StakeHolder, uint64)
}

type StakeHolderDB interface {
	StakeHolderView
	QueryAndUpdateStakeHolder(string, string, StakeHolderType) error
	StoreStakerHolder([]StakeHolder) error
}

type stakeHolderDB struct {
	gorm *gorm.DB
}

func (sh *stakeHolderDB) QueryAndUpdateStakeHolder(stakeAddress string, strategyAddress string, shType StakeHolderType) error {
	var stakeHolder StakeHolder
	result := sh.gorm.Table("staker_holder").Where("staker = ? and strategy = ?", stakeAddress, strategyAddress).Take(&stakeHolder)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			stHolder := StakeHolder{
				GUID:            uuid.New(),
				Staker:          common.HexToAddress(stakeAddress),
				Strategy:        common.HexToAddress(strategyAddress),
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

func (shv stakeHolderDB) GetStakeHolder(staker string) (*StakeHolder, error) {
	var stakeHolder StakeHolder
	result := shv.gorm.Table("staker_holder").Where("staker = ?", staker).Take(&stakeHolder)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &stakeHolder, nil
}

func (shv stakeHolderDB) ListStakeHolder(address string, page int, pageSize int, order string) ([]StakeHolder, uint64) {
	var totalRecord int64
	var stakeHolderList []StakeHolder
	queryRoot := shv.gorm.Table("staker_holder")
	if address != "0x00" {
		err := shv.gorm.Table("staker_holder").Select("total_manta_stake").Where("staker = ?", address).Count(&totalRecord).Error
		if err != nil {
			log.Error("get staker holder count fail")
		}
		queryRoot.Where("staker = ?", address).Offset((page - 1) * pageSize).Limit(pageSize)
	} else {
		err := shv.gorm.Table("staker_holder").Select("total_manta_stake").Count(&totalRecord).Error
		if err != nil {
			log.Error("get staker holder count fail ")
		}
		queryRoot.Offset((page - 1) * pageSize).Limit(pageSize)
	}
	if strings.ToLower(order) == "asc" {
		queryRoot.Order("timestamp asc")
	} else {
		queryRoot.Order("timestamp desc")
	}
	qErr := queryRoot.Find(&stakeHolderList).Error
	if qErr != nil {
		log.Error("list stake holder db fail", "err", qErr)
	}
	return stakeHolderList, uint64(totalRecord)
}

func NewStakeHolderDB(db *gorm.DB) StakeHolderDB {
	return &stakeHolderDB{gorm: db}
}

func (sh *stakeHolderDB) StoreStakerHolder(stakers []StakeHolder) error {
	result := sh.gorm.CreateInBatches(&stakers, len(stakers))
	return result.Error
}
