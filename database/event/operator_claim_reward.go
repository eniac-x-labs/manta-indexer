package event

import (
	"errors"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum/common"

	_ "github.com/eniac-x-labs/manta-indexer/database/utils/serializers"
)

type OperatorClaimReward struct {
	GUID      uuid.UUID      `gorm:"primaryKey" json:"guid"`
	BlockHash common.Hash    `json:"block_hash" gorm:"serializer:bytes"`
	Number    *big.Int       `json:"number" gorm:"serializer:u256"`
	TxHash    common.Hash    `json:"tx_hash" gorm:"serializer:bytes"`
	Operator  common.Address `json:"operator" gorm:"serializer:bytes"`
	Amount    *big.Int       `json:"amount" gorm:"serializer:u256"`
	IsHandle  uint8          `json:"is_handle"`
	Timestamp uint64         `json:"timestamp"`
}

func (OperatorClaimReward) TableName() string {
	return "operator_claim_reward"
}

type OperatorClaimRewardView interface {
	QueryUnHandleOperatorClaimReward() ([]OperatorClaimReward, error)
	GetOperatorClaimReward(string) (*OperatorClaimReward, error)
	ListOperatorClaimReward(address string, page int, pageSize int, order string) ([]OperatorClaimReward, uint64)
}

type OperatorClaimRewardDB interface {
	OperatorClaimRewardView
	MarkedOperatorClaimRewardHandled([]OperatorClaimReward) error
	StoreOperatorClaimReward([]OperatorClaimReward) error
}

type operatorClaimRewardDB struct {
	gorm *gorm.DB
}

func (oc operatorClaimRewardDB) QueryUnHandleOperatorClaimReward() ([]OperatorClaimReward, error) {
	var operatorClaimRewardList []OperatorClaimReward
	err := oc.gorm.Table("operator_claim_reward").Where("is_handle = ?", 0).Find(&operatorClaimRewardList).Error
	if err != nil {
		log.Error("get unhandled operator and staker reward list fail", "err", err)
		return nil, err
	}
	return operatorClaimRewardList, nil
}

func (oc operatorClaimRewardDB) MarkedOperatorClaimRewardHandled(rewards []OperatorClaimReward) error {
	for i := 0; i < len(rewards); i++ {
		var operatorClaimReward = OperatorClaimReward{}
		result := oc.gorm.Table("operator_claim_reward").Where(&OperatorClaimReward{GUID: rewards[i].GUID}).Take(&operatorClaimReward)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return nil
			}
			return result.Error
		}
		operatorClaimReward.IsHandle = 1
		err := oc.gorm.Save(operatorClaimReward).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (ocr operatorClaimRewardDB) GetOperatorClaimReward(address string) (*OperatorClaimReward, error) {
	var operatorClaimReward OperatorClaimReward
	result := ocr.gorm.Table("operator_claim_reward").Where(&OperatorClaimReward{Operator: common.HexToAddress(address)}).Take(&operatorClaimReward)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &operatorClaimReward, nil
}

func (ocr operatorClaimRewardDB) ListOperatorClaimReward(address string, page int, pageSize int, order string) ([]OperatorClaimReward, uint64) {
	address = strings.ToLower(address)

	var totalRecord int64
	var operatorClaimRewardList []OperatorClaimReward
	queryRoot := ocr.gorm.Table("operator_claim_reward")
	if address != "0x00" {
		err := ocr.gorm.Table("operator_claim_reward").Select("number").Where("operator = ?", address).Count(&totalRecord).Error
		if err != nil {
			log.Error("get operator claim reward count fail")
		}
		queryRoot.Where("operator = ?", address).Offset((page - 1) * pageSize).Limit(pageSize)
	} else {
		err := ocr.gorm.Table("operator_claim_reward").Select("number").Count(&totalRecord).Error
		if err != nil {
			log.Error("get operator claim reward count fail ")
		}
		queryRoot.Offset((page - 1) * pageSize).Limit(pageSize)
	}
	if strings.ToLower(order) == "asc" {
		queryRoot.Order("number asc")
	} else {
		queryRoot.Order("number desc")
	}
	qErr := queryRoot.Find(&operatorClaimRewardList).Error
	if qErr != nil {
		log.Error("list operatorClaimRewardDB fail", "err", qErr)
	}
	return operatorClaimRewardList, uint64(totalRecord)
}

func (oc operatorClaimRewardDB) StoreOperatorClaimReward(operatorClaimRewardList []OperatorClaimReward) error {
	result := oc.gorm.Table("operator_claim_reward").CreateInBatches(&operatorClaimRewardList, len(operatorClaimRewardList))
	return result.Error
}

func NewOperatorClaimRewardDB(db *gorm.DB) OperatorClaimRewardDB {
	return &operatorClaimRewardDB{gorm: db}
}
