package event

import (
	"errors"
	"github.com/ethereum/go-ethereum/log"
	"math/big"

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

type OperatorClaimRewardView interface {
	QueryUnHandleOperatorClaimReward() ([]OperatorClaimReward, error)
	QueryOperatorClaimRewardList(page int, pageSize int, order string) ([]OperatorClaimReward, uint64)
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
		result := oc.gorm.Where(&OperatorClaimReward{GUID: rewards[i].GUID}).Take(&operatorClaimReward)
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

func (oc operatorClaimRewardDB) QueryOperatorClaimRewardList(page int, pageSize int, order string) ([]OperatorClaimReward, uint64) {
	panic("implement me")
}

func (oc operatorClaimRewardDB) StoreOperatorClaimReward(operatorClaimRewardList []OperatorClaimReward) error {
	result := oc.gorm.Table("operator_claim_reward").CreateInBatches(&operatorClaimRewardList, len(operatorClaimRewardList))
	return result.Error
}

func NewOperatorClaimRewardDB(db *gorm.DB) OperatorClaimRewardDB {
	return &operatorClaimRewardDB{gorm: db}
}
