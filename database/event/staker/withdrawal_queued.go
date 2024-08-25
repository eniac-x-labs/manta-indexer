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

type WithdrawalQueued struct {
	GUID           uuid.UUID      `gorm:"primaryKey" json:"guid"`
	BlockHash      common.Hash    `json:"block_hash" gorm:"serializer:bytes"`
	Number         *big.Int       `json:"number" gorm:"serializer:u256"`
	TxHash         common.Hash    `json:"tx_hash" gorm:"serializer:bytes"`
	WithdrawalRoot common.Hash    `json:"withdrawal_root" gorm:"serializer:bytes"`
	Staker         common.Address `json:"staker" gorm:"serializer:bytes"`
	DelegatedTo    common.Address `json:"delegated_to" gorm:"serializer:bytes"`
	Withdrawer     common.Address `json:"withdrawer" gorm:"serializer:bytes"`
	Strategies     common.Address `json:"strategies" gorm:"serializer:bytes"`
	Shares         *big.Int       `json:"shares" gorm:"serializer:u256"`
	Nonce          *big.Int       `json:"nonce" gorm:"serializer:u256"`
	StartBlock     *big.Int       `json:"start_block" gorm:"serializer:u256"`
	IsHandle       uint8          `json:"is_handle"`
	Timestamp      uint64         `json:"timestamp"`
}

func (WithdrawalQueued) TableName() string {
	return "withdrawal_queued"
}

type WithdrawalQueuedView interface {
	ListWithdrawalQueued(address string, page int, pageSize int, order string) ([]WithdrawalQueued, uint64)
}

type WithdrawalQueuedDB interface {
	WithdrawalQueuedView
	MarkedWithdrawalQueuedHandled([]WithdrawalQueuedType) error
	StoreWithdrawalQueued([]WithdrawalQueued) error
}

func (wq withdrawalQueuedDB) GetWithdrawalQueued(address string) (*WithdrawalQueued, error) {
	var withdrawalQueued WithdrawalQueued
	result := wq.gorm.Table("withdrawal_queued").Where("withdrawer = ?", strings.ToLower(address)).Take(&withdrawalQueued)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &withdrawalQueued, nil
}

func (wq withdrawalQueuedDB) ListWithdrawalQueued(address string, page int, pageSize int, order string) ([]WithdrawalQueued, uint64) {
	var totalRecord int64
	var withdrawalQueuedList []WithdrawalQueued
	queryRoot := wq.gorm.Table("withdrawal_queued")
	if address != "0x00" {
		err := wq.gorm.Table("withdrawal_queued").Select("number").Where("staker = ?", address).Count(&totalRecord).Error
		if err != nil {
			log.Error("get withdrawal queued count fail")
		}
		queryRoot.Where("staker = ?", address).Offset((page - 1) * pageSize).Limit(pageSize)
	} else {
		err := wq.gorm.Table("withdrawal_queued").Select("number").Count(&totalRecord).Error
		if err != nil {
			log.Error("get withdrawal queued count fail ")
		}
		queryRoot.Offset((page - 1) * pageSize).Limit(pageSize)
	}
	if strings.ToLower(order) == "asc" {
		queryRoot.Order("number asc")
	} else {
		queryRoot.Order("number desc")
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

func (db withdrawalQueuedDB) MarkedWithdrawalQueuedHandled(withdrawalQueuedTypeList []WithdrawalQueuedType) error {
	for _, withdrawalQueuedType := range withdrawalQueuedTypeList {
		log.Info("withdrawal queued type", "staker", withdrawalQueuedType.Staker, "strategies", withdrawalQueuedType.Strategies)
		var withdrawalQueuedListTemp []WithdrawalQueued
		result := db.gorm.Table("withdrawal_queued").Where("staker = ? and strategies = ?", strings.ToLower(withdrawalQueuedType.Staker), strings.ToLower(withdrawalQueuedType.Strategies)).Find(&withdrawalQueuedListTemp)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				log.Warn("Not found", "err", result.Error)
				return nil
			}
			return result.Error
		}
		for _, withdrawalQueued := range withdrawalQueuedListTemp {
			log.Info("withdrawal queued", "staker", withdrawalQueued.Staker, "strategies", withdrawalQueued.Strategies)
			withdrawalQueued.IsHandle = 1
			err := db.gorm.Table("withdrawal_queued").Save(withdrawalQueued).Error
			if err != nil {
				return err
			}
		}
		return nil
	}
	return nil
}

func (db withdrawalQueuedDB) StoreWithdrawalQueued(withdrawalQueuedList []WithdrawalQueued) error {
	result := db.gorm.Table("withdrawal_queued").CreateInBatches(&withdrawalQueuedList, len(withdrawalQueuedList))
	return result.Error
}

func NewWithdrawalQueuedDB(db *gorm.DB) WithdrawalQueuedDB {
	return &withdrawalQueuedDB{gorm: db}
}
