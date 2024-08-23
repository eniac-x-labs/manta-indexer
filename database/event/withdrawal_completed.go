package event

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

type WithdrawalCompleted struct {
	GUID      uuid.UUID      `gorm:"primaryKey" json:"guid"`
	BlockHash common.Hash    `json:"block_hash" gorm:"serializer:bytes"`
	Number    *big.Int       `json:"number" gorm:"serializer:u256"`
	TxHash    common.Hash    `json:"tx_hash" gorm:"serializer:bytes"`
	Operator  common.Address `json:"staker" gorm:"serializer:bytes"`
	Staker    common.Address `json:"manta_token" gorm:"serializer:bytes"`
	Strategy  common.Address `json:"strategy" gorm:"serializer:bytes"`
	Shares    *big.Int       `json:"shares" gorm:"serializer:u256"`
	IsHandle  uint8          `json:"is_handle"`
	Timestamp uint64         `json:"timestamp"`
}

func (WithdrawalCompleted) TableName() string {
	return "withdrawal_completed"
}

type WithdrawalCompletedView interface {
	QueryUnHandleWithdrawalCompleted() ([]WithdrawalCompleted, error)

	GetWithdrawalCompleted(string) (*WithdrawalCompleted, error)
	ListWithdrawalCompleted(page int, pageSize int, order string) ([]WithdrawalCompleted, uint64)
}

type WithdrawalCompletedDB interface {
	WithdrawalCompletedView
	MarkedWithdrawalCompleted([]WithdrawalCompleted) error
	StoreWithdrawalCompleted([]WithdrawalCompleted) error
}

type withdrawalCompletedDB struct {
	gorm *gorm.DB
}

func (wc withdrawalCompletedDB) QueryUnHandleWithdrawalCompleted() ([]WithdrawalCompleted, error) {
	var withdrawalCompletedList []WithdrawalCompleted
	err := wc.gorm.Table("withdrawal_completed").Where("is_handle = ?", 0).Find(&withdrawalCompletedList).Error
	if err != nil {
		log.Error("get unhandled withdraw completed fail", "err", err)
		return nil, err
	}
	return withdrawalCompletedList, nil
}

func (wc withdrawalCompletedDB) MarkedWithdrawalCompleted(withdrawalCompletedList []WithdrawalCompleted) error {
	for i := 0; i < len(withdrawalCompletedList); i++ {
		var withdrawalCompleted = WithdrawalCompleted{}
		result := wc.gorm.Table("withdrawal_completed").Where(&WithdrawalCompleted{GUID: withdrawalCompletedList[i].GUID}).Take(&withdrawalCompleted)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return nil
			}
			return result.Error
		}
		withdrawalCompleted.IsHandle = 1
		err := wc.gorm.Save(withdrawalCompleted).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (wc withdrawalCompletedDB) GetWithdrawalCompleted(address string) (*WithdrawalCompleted, error) {
	var withdrawalCompleted WithdrawalCompleted
	result := wc.gorm.Table("withdrawal_completed").Where(&WithdrawalCompleted{Staker: common.HexToAddress(address)}).Take(&withdrawalCompleted)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &withdrawalCompleted, nil
}

func (wc withdrawalCompletedDB) ListWithdrawalCompleted(page int, pageSize int, order string) ([]WithdrawalCompleted, uint64) {
	var totalRecord int64
	var withdrawalCompletedList []WithdrawalCompleted
	queryRoot := wc.gorm.Table("withdrawal_completed")
	err := wc.gorm.Table("withdrawal_completed").Select("number").Count(&totalRecord).Error
	if err != nil {
		log.Error("list withdrawalCompletedDB count fail", "err", err)
	}

	queryRoot.Offset((page - 1) * pageSize).Limit(pageSize)
	if strings.ToLower(order) == "asc" {
		queryRoot.Order("number asc")
	} else {
		queryRoot.Order("number desc")
	}
	qErr := queryRoot.Find(&withdrawalCompletedList).Error
	if qErr != nil {
		log.Error("list withdrawalCompletedDB fail", "err", qErr)
	}
	return withdrawalCompletedList, uint64(totalRecord)
}

func (wc withdrawalCompletedDB) StoreWithdrawalCompleted(withdrawalCompletedList []WithdrawalCompleted) error {
	result := wc.gorm.Table("withdrawal_completed").CreateInBatches(&withdrawalCompletedList, len(withdrawalCompletedList))
	return result.Error
}

func NewWithdrawalCompletedDB(db *gorm.DB) WithdrawalCompletedDB {
	return &withdrawalCompletedDB{gorm: db}
}
