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

func (WithdrawalQueued) TableName() string {
	return "withdrawal_queued"
}

type WithdrawalQueuedView interface {
	GetWithdrawalQueued(string) (*WithdrawalQueued, error)
	ListWithdrawalQueued(page int, pageSize int, order string) ([]WithdrawalQueued, uint64)
}

type WithdrawalQueuedDB interface {
	WithdrawalQueuedView
	StoreWithdrawalQueued([]WithdrawalQueued) error
}

func (wq withdrawalQueuedDB) GetWithdrawalQueued(address string) (*WithdrawalQueued, error) {
	var withdrawalQueued WithdrawalQueued
	result := wq.gorm.Table("withdrawal_queued").Where(&WithdrawalQueued{Withdrawer: common.HexToAddress(address)}).Take(&withdrawalQueued)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &withdrawalQueued, nil
}

func (wq withdrawalQueuedDB) ListWithdrawalQueued(page int, pageSize int, order string) ([]WithdrawalQueued, uint64) {
	var totalRecord int64
	var withdrawalQueuedList []WithdrawalQueued
	queryRoot := wq.gorm.Table("withdrawal_queued")
	err := queryRoot.Select("guid").Count(&totalRecord).Error
	if err != nil {
		log.Error("list withdrawalQueuedDB count fail", "err", err)
	}

	queryRoot.Offset((page - 1) * pageSize).Limit(pageSize)
	if strings.ToLower(order) == "asc" {
		queryRoot.Order("guid asc")
	} else {
		queryRoot.Order("guid desc")
	}
	qErr := queryRoot.Find(&withdrawalQueuedList).Error
	if qErr != nil {
		log.Error("list withdrawalQueuedDB fail", "err", qErr)
	}
	return withdrawalQueuedList, uint64(totalRecord)
}

type withdrawalQueuedDB struct {
	gorm *gorm.DB
}

func (db withdrawalQueuedDB) StoreWithdrawalQueued(withdrawalQueuedList []WithdrawalQueued) error {
	result := db.gorm.Table("withdrawal_queued").CreateInBatches(&withdrawalQueuedList, len(withdrawalQueuedList))
	return result.Error
}

func NewWithdrawalQueuedDB(db *gorm.DB) WithdrawalQueuedDB {
	return &withdrawalQueuedDB{gorm: db}
}
