package event

import (
	"math/big"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum/common"

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

type StakeHolderClaimRewardView interface {
	QueryUnHandleStakeHolderClaimReward() ([]StakeHolderClaimReward, error)
	QueryStakeHolderClaimRewardList(page int, pageSize int, order string) ([]StakeHolderClaimReward, uint64)
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
	panic("implement me")
}

func (shc stakeHolderClaimRewardDB) MarkedStakeHolderClaimRewardHandled(stakeHolderClaimReward []StakeHolderClaimReward) error {
	panic("implement me")
}

func (shc stakeHolderClaimRewardDB) QueryStakeHolderClaimRewardList(page int, pageSize int, order string) ([]StakeHolderClaimReward, uint64) {
	panic("implement me")
}

func (shc stakeHolderClaimRewardDB) StoreStakeHolderClaimReward(stakeHolderClaimRewardList []StakeHolderClaimReward) error {
	result := shc.gorm.CreateInBatches(&stakeHolderClaimRewardList, len(stakeHolderClaimRewardList))
	return result.Error
}

func NewStakeHolderClaimRewardDB(db *gorm.DB) StakeHolderClaimRewardDB {
	return &stakeHolderClaimRewardDB{gorm: db}
}
