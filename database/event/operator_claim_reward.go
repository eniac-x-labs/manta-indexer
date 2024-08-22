package event

import (
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
	QueryOperatorClaimRewardList(page int, pageSize int, order string) ([]OperatorClaimReward, uint64)
}

type OperatorClaimRewardDB interface {
	OperatorClaimRewardView
	StoreOperatorClaimReward([]OperatorClaimReward) error
}

type operatorClaimRewardDB struct {
	gorm *gorm.DB
}

func (db operatorClaimRewardDB) QueryOperatorClaimRewardList(page int, pageSize int, order string) ([]OperatorClaimReward, uint64) {
	panic("implement me")
}

func (db operatorClaimRewardDB) StoreOperatorClaimReward(operatorClaimRewardList []OperatorClaimReward) error {
	result := db.gorm.CreateInBatches(&operatorClaimRewardList, len(operatorClaimRewardList))
	return result.Error
}

func NewOperatorClaimRewardDB(db *gorm.DB) OperatorClaimRewardDB {
	return &operatorClaimRewardDB{gorm: db}
}
