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

type OperatorRegistered struct {
	GUID                     uuid.UUID      `gorm:"primaryKey" json:"guid"`
	BlockHash                common.Hash    `json:"block_hash" gorm:"serializer:bytes"`
	Number                   *big.Int       `json:"number" gorm:"serializer:u256"`
	TxHash                   common.Hash    `json:"tx_hash" gorm:"serializer:bytes"`
	Operator                 common.Address `json:"operator" gorm:"serializer:bytes"`
	EarningsReceiver         common.Address `json:"earnings_receiver" gorm:"serializer:bytes"`
	DelegationApprover       common.Address `json:"delegation_approver" gorm:"serializer:bytes"`
	StakerOptoutWindowBlocks *big.Int       `json:"staker_optout_window_blocks" gorm:"serializer:u256"`
	IsHandle                 uint8          `json:"is_handle"`
	Timestamp                uint64         `json:"timestamp"`
}

func (OperatorRegistered) TableName() string {
	return "operator_registered"
}

type OperatorRegisteredView interface {
	QueryUnHandleOperatorRegistered() ([]OperatorRegistered, error)
	QueryOperatorRegistered(string) (*OperatorRegistered, error)
	QueryOperatorRegisteredList(page int, pageSize int, order string) ([]OperatorRegistered, uint64)
}

type OperatorRegisteredDB interface {
	OperatorRegisteredView
	MarkedOperatorRegisteredHandled([]OperatorRegistered) error
	StoreOperatorRegistered([]OperatorRegistered) error
}

type operatorRegisteredDB struct {
	gorm *gorm.DB
}

func (or operatorRegisteredDB) QueryOperatorRegistered(operator string) (*OperatorRegistered, error) {
	var operatorRegistered OperatorRegistered
	result := or.gorm.Table("operator_registered").Where(&OperatorRegistered{Operator: common.HexToAddress(operator)}).Take(&operatorRegistered)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &operatorRegistered, nil
}

func (or operatorRegisteredDB) MarkedOperatorRegisteredHandled(unHandledOperatorRegistered []OperatorRegistered) error {
	for i := 0; i < len(unHandledOperatorRegistered); i++ {
		var operatorRegistereds = OperatorRegistered{}
		result := or.gorm.Table("operator_registered").Where(&OperatorRegistered{GUID: unHandledOperatorRegistered[i].GUID}).Take(&operatorRegistereds)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return nil
			}
			return result.Error
		}
		operatorRegistereds.IsHandle = 1
		err := or.gorm.Table("operator_registered").Save(operatorRegistereds).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (or operatorRegisteredDB) QueryUnHandleOperatorRegistered() ([]OperatorRegistered, error) {
	var operatorRegisteredList []OperatorRegistered
	err := or.gorm.Table("operator_registered").Where("is_handle = ?", 0).Find(&operatorRegisteredList).Error
	if err != nil {
		log.Error("get unhandled operator registered fail", "err", err)
		return nil, err
	}
	return operatorRegisteredList, nil
}

func (or operatorRegisteredDB) QueryOperatorRegisteredList(page int, pageSize int, order string) ([]OperatorRegistered, uint64) {
	var totalRecord int64
	var operatorRegisteredList []OperatorRegistered
	queryStateRoot := or.gorm.Table("operator_registered")
	err := or.gorm.Table("operator_registered").Select("number").Count(&totalRecord).Error
	if err != nil {
		log.Error("get operator registered count fail", "err", err)
	}
	queryStateRoot.Offset((page - 1) * pageSize).Limit(pageSize)
	if strings.ToLower(order) == "asc" {
		queryStateRoot.Order("number asc")
	} else {
		queryStateRoot.Order("number desc")
	}
	qErr := queryStateRoot.Find(&operatorRegisteredList).Error
	if qErr != nil {
		log.Error("get operator registered list fail", "err", qErr)
	}
	return operatorRegisteredList, uint64(totalRecord)
}

func (or operatorRegisteredDB) StoreOperatorRegistered(operatorRegisteredList []OperatorRegistered) error {
	result := or.gorm.Table("operator_registered").CreateInBatches(&operatorRegisteredList, len(operatorRegisteredList))
	return result.Error
}

func NewOperatorRegisteredDB(db *gorm.DB) OperatorRegisteredDB {
	return &operatorRegisteredDB{gorm: db}
}
