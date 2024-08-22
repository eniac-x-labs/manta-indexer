package event

import (
	"math/big"

	_ "github.com/eniac-x-labs/manta-indexer/database/utils/serializers"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OperatorAndStakeReward struct {
	GUID        uuid.UUID      `gorm:"primaryKey" json:"guid"`
	BlockHash   common.Hash    `json:"block_hash" gorm:"serializer:bytes"`
	Number      *big.Int       `json:"number" gorm:"serializer:u256"`
	TxHash      common.Hash    `json:"tx_hash" gorm:"serializer:bytes"`
	Strategy    common.Address `json:"strategy" gorm:"serializer:bytes"`
	Operator    common.Address `json:"operator" gorm:"serializer:bytes"`
	StakerFee   *big.Int       `json:"staker_fee" gorm:"serializer:u256"`
	OperatorFee *big.Int       `json:"operator_fee" gorm:"serializer:u256"`
	IsHandle    uint8          `json:"is_handle"`
	Timestamp   uint64         `json:"timestamp"`
}

type OperatorAndStakeRewardView interface {
	QueryOperatorAndStakeRewardList(page int, pageSize int, order string) ([]OperatorAndStakeReward, uint64)
}

type OperatorAndStakeRewardDB interface {
	OperatorAndStakeRewardView
	StoreOperatorAndStakeReward([]OperatorAndStakeReward) error
}

type operatorAndStakeRewardDB struct {
	gorm *gorm.DB
}

func (db operatorAndStakeRewardDB) QueryOperatorAndStakeRewardList(page int, pageSize int, order string) ([]OperatorAndStakeReward, uint64) {
	panic("implement me")
}

func (db operatorAndStakeRewardDB) StoreOperatorAndStakeReward(operatorAndStakeRewardList []OperatorAndStakeReward) error {
	result := db.gorm.CreateInBatches(&operatorAndStakeRewardList, len(operatorAndStakeRewardList))
	return result.Error
}

func NewOperatorAndStakeRewardDB(db *gorm.DB) OperatorAndStakeRewardDB {
	return &operatorAndStakeRewardDB{gorm: db}
}
