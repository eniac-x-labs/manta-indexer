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

type StakeStrategy struct {
	GUID          uuid.UUID      `gorm:"primaryKey" json:"guid"`
	Staker        common.Address `gorm:"serializer:bytes" json:"staker"`
	Strategy      common.Address `gorm:"serializer:bytes" json:"strategy"`
	TotalStake    *big.Int       `gorm:"serializer:u256" json:"total_manta_stake"`
	TotalReward   *big.Int       `gorm:"serializer:u256" json:"total_reward"`
	ClaimedAmount *big.Int       `gorm:"serializer:u256" json:"claimed_amount"`
	Timestamp     uint64
}

func (StakeStrategy) TableName() string {
	return "staker_strategy"
}

type StakeStrategyView interface {
	GetStakeStrategy(string) (*StakeStrategy, error)
	ListStakeStrategy(address string, page int, pageSize int, order string) ([]StakeStrategy, uint64)
}

type StakeStrategyDB interface {
	StakeStrategyView
	QueryAndUpdateStakeStrategy(string, string, StakeStrategyOperatorType) error
	StoreStakeStrategy([]StakeStrategy) error
}

type stakeStrategyDB struct {
	gorm *gorm.DB
}

func (sh *stakeStrategyDB) QueryAndUpdateStakeStrategy(stakeAddress string, strategyAddress string, shType StakeStrategyOperatorType) error {
	var stakeStrategy StakeStrategy
	result := sh.gorm.Table("staker_strategy").Where("staker = ? and strategy = ?", strings.ToLower(stakeAddress), strings.ToLower(strategyAddress)).Take(&stakeStrategy)
	if result.Error != nil {
		log.Warn("staker strategy query warning", "staker", stakeAddress, "strategy", strategyAddress, "err", result.Error)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			stHolder := StakeStrategy{
				GUID:          uuid.New(),
				Staker:        common.HexToAddress(stakeAddress),
				Strategy:      common.HexToAddress(strategyAddress),
				TotalStake:    shType.MantaStake,
				TotalReward:   shType.Reward,
				ClaimedAmount: shType.ClaimedAmount,
				Timestamp:     shType.Timestamp,
			}
			err := sh.gorm.Create(stHolder).Error
			if err != nil {
				log.Error("create stake strategy fail", "err", err)
			}
		}
		return nil
	}
	if shType.MantaStake != nil {
		stakeStrategy.TotalStake = new(big.Int).Add(stakeStrategy.TotalStake, shType.MantaStake)
	}
	if shType.Reward != nil {
		stakeStrategy.TotalReward = new(big.Int).Add(stakeStrategy.TotalReward, shType.Reward)
	}
	if shType.ClaimedAmount != nil {
		stakeStrategy.ClaimedAmount = new(big.Int).Add(stakeStrategy.ClaimedAmount, shType.ClaimedAmount)
	}
	err := sh.gorm.Save(stakeStrategy).Error
	if err != nil {
		log.Error("update stake strategy fail", "err", err)
		return err
	}
	return nil
}

func (shv stakeStrategyDB) GetStakeStrategy(staker string) (*StakeStrategy, error) {
	var stakeStrategy StakeStrategy
	result := shv.gorm.Table("staker_strategy").Where("staker = ?", strings.ToLower(staker)).Take(&stakeStrategy)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &stakeStrategy, nil
}

func (shv stakeStrategyDB) ListStakeStrategy(address string, page int, pageSize int, order string) ([]StakeStrategy, uint64) {
	var totalRecord int64
	var stakeStrategyList []StakeStrategy
	queryRoot := shv.gorm.Table("staker_strategy")
	if address != "0x00" {
		err := shv.gorm.Table("staker_strategy").Select("total_stake").Where("staker = ?", address).Count(&totalRecord).Error
		if err != nil {
			log.Error("get stake strategy count fail")
		}
		queryRoot.Where("staker = ?", address).Offset((page - 1) * pageSize).Limit(pageSize)
	} else {
		err := shv.gorm.Table("staker_strategy").Select("total_manta_stake").Count(&totalRecord).Error
		if err != nil {
			log.Error("get stake strategy count fail ")
		}
		queryRoot.Offset((page - 1) * pageSize).Limit(pageSize)
	}
	if strings.ToLower(order) == "asc" {
		queryRoot.Order("timestamp asc")
	} else {
		queryRoot.Order("timestamp desc")
	}
	qErr := queryRoot.Find(&stakeStrategyList).Error
	if qErr != nil {
		log.Error("list stake strategy db fail", "err", qErr)
	}
	return stakeStrategyList, uint64(totalRecord)
}

func NewStakeStrategyDB(db *gorm.DB) StakeStrategyDB {
	return &stakeStrategyDB{gorm: db}
}

func (sh *stakeStrategyDB) StoreStakeStrategy(stakers []StakeStrategy) error {
	result := sh.gorm.CreateInBatches(&stakers, len(stakers))
	return result.Error
}
