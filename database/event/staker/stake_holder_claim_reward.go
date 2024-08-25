package staker

import (
	"errors"
	"math/big"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"

	_ "github.com/eniac-x-labs/manta-indexer/database/utils/serializers"
)

type StakeHolderClaimReward struct {
	GUID        uuid.UUID      `gorm:"primaryKey" json:"guid"`
	BlockHash   common.Hash    `json:"block_hash" gorm:"serializer:bytes"`
	Number      *big.Int       `json:"number" gorm:"serializer:u256"`
	TxHash      common.Hash    `json:"tx_hash" gorm:"serializer:bytes"`
	StakeHolder common.Address `json:"stake_holder" gorm:"serializer:bytes"`
	Strategy    common.Address `json:"strategy" gorm:"serializer:bytes"`
	Amount      *big.Int       `json:"amount" gorm:"serializer:u256"`
	IsHandle    uint8          `json:"is_handle"`
	Timestamp   uint64         `json:"timestamp"`
}

func (StakeHolderClaimReward) TableName() string {
	return "stake_holder_claim_reward"
}

type StakeHolderClaimRewardView interface {
	QueryUnHandleStakeHolderClaimReward() ([]StakeHolderClaimReward, error)
	ListStakeHolderClaimReward(address string, page int, pageSize int, order string) ([]StakeHolderClaimReward, uint64)
}

type StakeHolderClaimRewardDB interface {
	StakeHolderClaimRewardView
	MarkedStakeHolderClaimRewardHandled(stakeHolderClaimReward []StakeHolderClaimReward) error
	StoreStakeHolderClaimReward([]StakeHolderClaimReward) error
}

type stakeHolderClaimRewardDB struct {
	gorm *gorm.DB
}

func (shc stakeHolderClaimRewardDB) QueryUnHandleStakeHolderClaimReward() ([]StakeHolderClaimReward, error) {
	var stakeHolderClaimReward []StakeHolderClaimReward
	err := shc.gorm.Table("stake_holder_claim_reward").Where("is_handle = ?", 0).Find(&stakeHolderClaimReward).Error
	if err != nil {
		log.Error("get stake holder claim reward fail", "err", err)
		return nil, err
	}
	return stakeHolderClaimReward, nil
}

func (shc stakeHolderClaimRewardDB) MarkedStakeHolderClaimRewardHandled(stakeHolderClaimRewards []StakeHolderClaimReward) error {
	for i := 0; i < len(stakeHolderClaimRewards); i++ {
		var stakeHolderClaimReward = StakeHolderClaimReward{}
		result := shc.gorm.Table("stake_holder_claim_reward").Where("guid = ?", stakeHolderClaimRewards[i].GUID).Take(&stakeHolderClaimReward)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return nil
			}
			return result.Error
		}
		stakeHolderClaimReward.IsHandle = 1
		err := shc.gorm.Save(stakeHolderClaimReward).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (shcr stakeHolderClaimRewardDB) ListStakeHolderClaimReward(address string, page int, pageSize int, order string) ([]StakeHolderClaimReward, uint64) {
	var totalRecord int64
	var stakeHolderClaimRewardList []StakeHolderClaimReward
	queryRoot := shcr.gorm.Table("stake_holder_claim_reward")
	if address != "0x00" {
		err := shcr.gorm.Table("stake_holder_claim_reward").Select("number").Where("stake_holder = ?", address).Count(&totalRecord).Error
		if err != nil {
			log.Error("get stakeholder claim reward count fail")
		}
		queryRoot.Where("stake_holder = ?", address).Offset((page - 1) * pageSize).Limit(pageSize)
	} else {
		err := shcr.gorm.Table("stake_holder_claim_reward").Select("number").Count(&totalRecord).Error
		if err != nil {
			log.Error("get stakeholder claim reward count fail ")
		}
		queryRoot.Offset((page - 1) * pageSize).Limit(pageSize)
	}
	if strings.ToLower(order) == "asc" {
		queryRoot.Order("number asc")
	} else {
		queryRoot.Order("number desc")
	}
	qErr := queryRoot.Find(&stakeHolderClaimRewardList).Error
	if qErr != nil {
		log.Error("list stakeHolderClaimRewardDB fail", "err", qErr)
	}
	return stakeHolderClaimRewardList, uint64(totalRecord)
}

func (shc stakeHolderClaimRewardDB) StoreStakeHolderClaimReward(stakeHolderClaimRewardList []StakeHolderClaimReward) error {
	result := shc.gorm.Table("stake_holder_claim_reward").CreateInBatches(&stakeHolderClaimRewardList, len(stakeHolderClaimRewardList))
	return result.Error
}

func NewStakeHolderClaimRewardDB(db *gorm.DB) StakeHolderClaimRewardDB {
	return &stakeHolderClaimRewardDB{gorm: db}
}
