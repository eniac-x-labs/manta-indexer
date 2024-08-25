package operator

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

type OperatorAndStakeReward struct {
	GUID             uuid.UUID      `gorm:"primaryKey" json:"guid"`
	BlockHash        common.Hash    `json:"block_hash" gorm:"serializer:bytes"`
	Number           *big.Int       `json:"number" gorm:"serializer:u256"`
	TxHash           common.Hash    `json:"tx_hash" gorm:"serializer:bytes"`
	Strategy         common.Address `json:"strategy" gorm:"serializer:bytes"`
	Operator         common.Address `json:"operator" gorm:"serializer:bytes"`
	StakerFee        *big.Int       `json:"staker_fee" gorm:"serializer:u256"`
	OperatorFee      *big.Int       `json:"operator_fee" gorm:"serializer:u256"`
	IsOperatorHandle uint8          `json:"is_operator_handle"`
	IsStakerHandle   uint8          `json:"is_staker_handle"`
	Timestamp        uint64         `json:"timestamp"`
}

func (OperatorAndStakeReward) TableName() string {
	return "operator_and_stake_reward"
}

type OperatorAndStakeRewardView interface {
	QueryUnHandleOperatorAndStakeReward(isOperator bool) ([]OperatorAndStakeReward, error)
	GetOperatorAndStakeReward(string) (*OperatorAndStakeReward, error)
	ListOperatorAndStakeReward(address string, page int, pageSize int, order string) ([]OperatorAndStakeReward, uint64)
}

type OperatorAndStakeRewardDB interface {
	OperatorAndStakeRewardView
	MarkedOperatorAndStakeRewardHandled([]OperatorAndStakeReward, bool) error
	StoreOperatorAndStakeReward([]OperatorAndStakeReward) error
}

type operatorAndStakeRewardDB struct {
	gorm *gorm.DB
}

func (oas operatorAndStakeRewardDB) QueryUnHandleOperatorAndStakeReward(isOperator bool) ([]OperatorAndStakeReward, error) {
	var operatorAndStakeRewardList []OperatorAndStakeReward
	var err error
	if isOperator {
		err = oas.gorm.Table("operator_and_stake_reward").Where("is_operator_handle = ?", 0).Find(&operatorAndStakeRewardList).Error
	} else {
		err = oas.gorm.Table("operator_and_stake_reward").Where("is_staker_handle = ?", 0).Find(&operatorAndStakeRewardList).Error
	}
	if err != nil {
		log.Error("get unhandled operator and staker reward list fail", "err", err)
		return nil, err
	}
	return operatorAndStakeRewardList, nil
}

func (oas operatorAndStakeRewardDB) MarkedOperatorAndStakeRewardHandled(rewards []OperatorAndStakeReward, isOperator bool) error {
	for i := 0; i < len(rewards); i++ {
		var operatorAndStakeReward = OperatorAndStakeReward{}
		result := oas.gorm.Table("operator_and_stake_reward").Where(&OperatorAndStakeReward{GUID: rewards[i].GUID}).Take(&operatorAndStakeReward)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return nil
			}
			return result.Error
		}
		if isOperator {
			operatorAndStakeReward.IsOperatorHandle = 1
		} else {
			operatorAndStakeReward.IsStakerHandle = 1
		}
		err := oas.gorm.Save(operatorAndStakeReward).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (osr operatorAndStakeRewardDB) GetOperatorAndStakeReward(address string) (*OperatorAndStakeReward, error) {
	var operatorAndStakeReward OperatorAndStakeReward
	result := osr.gorm.Table("operator_and_stake_reward").Where("operator = ?", strings.ToLower(address)).Take(&operatorAndStakeReward)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &operatorAndStakeReward, nil
}

func (osr operatorAndStakeRewardDB) ListOperatorAndStakeReward(address string, page int, pageSize int, order string) ([]OperatorAndStakeReward, uint64) {
	address = strings.ToLower(address)

	var totalRecord int64
	var operatorAndStakeRewardList []OperatorAndStakeReward
	queryRoot := osr.gorm.Table("operator_and_stake_reward")
	if address != "0x00" {
		err := osr.gorm.Table("operator_and_stake_reward").Select("number").Where("operator = ?", address).Count(&totalRecord).Error
		if err != nil {
			log.Error("get operator and staker reward count fail")
		}
		queryRoot.Where("operator = ?", address).Offset((page - 1) * pageSize).Limit(pageSize)
	} else {
		err := osr.gorm.Table("operator_and_stake_reward").Select("number").Count(&totalRecord).Error
		if err != nil {
			log.Error("get operator and staker reward count fail ")
		}
		queryRoot.Offset((page - 1) * pageSize).Limit(pageSize)
	}
	if strings.ToLower(order) == "asc" {
		queryRoot.Order("number asc")
	} else {
		queryRoot.Order("number desc")
	}
	qErr := queryRoot.Find(&operatorAndStakeRewardList).Error
	if qErr != nil {
		log.Error("list operatorAndStakeRewardDB fail", "err", qErr)
	}
	return operatorAndStakeRewardList, uint64(totalRecord)
}

func (oas operatorAndStakeRewardDB) StoreOperatorAndStakeReward(operatorAndStakeRewardList []OperatorAndStakeReward) error {
	result := oas.gorm.Table("operator_and_stake_reward").CreateInBatches(&operatorAndStakeRewardList, len(operatorAndStakeRewardList))
	return result.Error
}

func NewOperatorAndStakeRewardDB(db *gorm.DB) OperatorAndStakeRewardDB {
	return &operatorAndStakeRewardDB{gorm: db}
}
