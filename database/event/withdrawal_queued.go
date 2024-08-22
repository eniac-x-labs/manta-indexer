package event

import (
	"math/big"

	_ "github.com/eniac-x-labs/manta-indexer/database/utils/serializers"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WithdrawalQueued struct {
	GUID           uuid.UUID      `gorm:"primaryKey" json:"guid"`
	BlockHash      common.Hash    `json:"block_hash" gorm:"serializer:bytes"`
	Number         *big.Int       `json:"number" gorm:"serializer:u256"`
	TxHash         common.Hash    `json:"tx_hash" gorm:"serializer:bytes"`
	WithdrawalRoot common.Hash    `json:"withdrawal_root" gorm:"serializer:bytes"`
	Staker         common.Address `json:"staker" gorm:"serializer:bytes"`
	DelegatedTo    common.Address `json:"delegated_to" gorm:"serializer:bytes"`
	Withdrawer     common.Address `json:"withdrawer" gorm:"serializer:bytes"`
	Nonce          *big.Int       `json:"nonce" gorm:"serializer:u256"`
	StartBlock     *big.Int       `json:"start_block" gorm:"serializer:u256"`
	Strategies     string         `json:"strategies"`
	Shares         string         `json:"shares"`
	IsHandle       uint8          `json:"is_handle"`
	Timestamp      uint64         `json:"timestamp"`
}

type WithdrawalQueuedView interface {
	QueryWithdrawalQueuedList(page int, pageSize int, order string) ([]WithdrawalQueued, uint64)
}

type WithdrawalQueuedDB interface {
	WithdrawalQueuedView
	StoreWithdrawalQueued([]WithdrawalQueued) error
}

type withdrawalQueuedDB struct {
	gorm *gorm.DB
}

func (db withdrawalQueuedDB) QueryWithdrawalQueuedList(page int, pageSize int, order string) ([]WithdrawalQueued, uint64) {
	panic("implement me")
}

func (db withdrawalQueuedDB) StoreWithdrawalQueued(withdrawalQueuedList []WithdrawalQueued) error {
	result := db.gorm.CreateInBatches(&withdrawalQueuedList, len(withdrawalQueuedList))
	return result.Error
}

func NewWithdrawalQueuedDB(db *gorm.DB) WithdrawalQueuedDB {
	return &withdrawalQueuedDB{gorm: db}
}
