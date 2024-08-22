package event

import (
	"errors"
	"math/big"

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

type OperatorRegisteredView interface {
	QueryUnHandleOperatorRegistered() ([]OperatorRegistered, error)
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

func (or operatorRegisteredDB) MarkedOperatorRegisteredHandled(unHandledOperatorRegistered []OperatorRegistered) error {
	for i := 0; i < len(unHandledOperatorRegistered); i++ {
		var operatorRegistereds = OperatorRegistered{}
		result := or.gorm.Where(&OperatorRegistered{GUID: unHandledOperatorRegistered[i].GUID}).Take(&operatorRegistereds)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return nil
			}
			return result.Error
		}
		operatorRegistereds.IsHandle = 1
		err := or.gorm.Save(operatorRegistereds).Error
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
	panic("implement me")
}

func (or operatorRegisteredDB) StoreOperatorRegistered(operatorRegisteredList []OperatorRegistered) error {
	result := or.gorm.CreateInBatches(&operatorRegisteredList, len(operatorRegisteredList))
	return result.Error
}

func NewOperatorRegisteredDB(db *gorm.DB) OperatorRegisteredDB {
	return &operatorRegisteredDB{gorm: db}
}
