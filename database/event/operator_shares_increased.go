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

type OperatorSharesIncreased struct {
	GUID      uuid.UUID      `gorm:"primaryKey" json:"guid"`
	BlockHash common.Hash    `json:"block_hash" gorm:"serializer:bytes"`
	Number    *big.Int       `json:"number" gorm:"serializer:u256"`
	TxHash    common.Hash    `json:"tx_hash" gorm:"serializer:bytes"`
	Operator  common.Address `json:"operator" gorm:"serializer:bytes"`
	Staker    common.Address `json:"staker" gorm:"serializer:bytes"`
	Strategy  common.Address `json:"strategy" gorm:"serializer:bytes"`
	Shares    *big.Int       `json:"shares" gorm:"serializer:u256"`
	IsHandle  uint8          `json:"is_handle"`
	Timestamp uint64         `json:"timestamp"`
}

func (OperatorSharesIncreased) TableName() string {
	return "operator_shares_increased"
}

type OperatorSharesIncreasedView interface {
	QueryUnHandleOperatorSharesIncreased() ([]OperatorSharesIncreased, error)
	ListOperatorSharesIncreased(address string, page int, pageSize int, order string) ([]OperatorSharesIncreased, uint64)
}

type OperatorSharesIncreasedDB interface {
	OperatorSharesIncreasedView
	MarkedOperatorSharesIncreasedHandled(unHandleOperatorSharesIncreased []OperatorSharesIncreased) error
	StoreOperatorSharesIncreased([]OperatorSharesIncreased) error
}

type operatorSharesIncreasedDB struct {
	gorm *gorm.DB
}

func (osi operatorSharesIncreasedDB) MarkedOperatorSharesIncreasedHandled(unHandleOperatorSharesIncreased []OperatorSharesIncreased) error {
	for i := 0; i < len(unHandleOperatorSharesIncreased); i++ {
		var operatorSharesIncreased = OperatorSharesIncreased{}
		result := osi.gorm.Table("operator_shares_increased").Where(&OperatorRegistered{GUID: unHandleOperatorSharesIncreased[i].GUID}).Take(&operatorSharesIncreased)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return nil
			}
			return result.Error
		}
		operatorSharesIncreased.IsHandle = 1
		err := osi.gorm.Save(operatorSharesIncreased).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (osi operatorSharesIncreasedDB) QueryUnHandleOperatorSharesIncreased() ([]OperatorSharesIncreased, error) {
	var operatorSharesIncreasedList []OperatorSharesIncreased
	err := osi.gorm.Table("operator_shares_increased").Where("is_handle = ?", 0).Find(&operatorSharesIncreasedList).Error
	if err != nil {
		log.Error("get operator share increased list fail", "err", err)
		return nil, err
	}
	return operatorSharesIncreasedList, nil
}

func (osi operatorSharesIncreasedDB) ListOperatorSharesIncreased(address string, page int, pageSize int, order string) ([]OperatorSharesIncreased, uint64) {
	var totalRecord int64
	var operatorSharesIncreasedList []OperatorSharesIncreased
	queryRoot := osi.gorm.Table("operator_shares_increased")
	if address != "0x00" {
		err := osi.gorm.Table("operator_shares_increased").Select("number").Where("operator = ?", address).Count(&totalRecord).Error
		if err != nil {
			log.Error("get operator share increased count fail")
		}
		queryRoot.Where("operator = ?", address).Offset((page - 1) * pageSize).Limit(pageSize)
	} else {
		err := osi.gorm.Table("operator_shares_increased").Select("number").Count(&totalRecord).Error
		if err != nil {
			log.Error("get operator share increased count fail ")
		}
		queryRoot.Offset((page - 1) * pageSize).Limit(pageSize)
	}
	if strings.ToLower(order) == "asc" {
		queryRoot.Order("number asc")
	} else {
		queryRoot.Order("number desc")
	}
	qErr := queryRoot.Find(&operatorSharesIncreasedList).Error
	if qErr != nil {
		log.Error("list operatorSharesIncreasedList fail", "err", qErr)
	}
	return operatorSharesIncreasedList, uint64(totalRecord)
}

func (osi operatorSharesIncreasedDB) StoreOperatorSharesIncreased(operatorSharesIncreasedList []OperatorSharesIncreased) error {
	result := osi.gorm.Table("operator_shares_increased").CreateInBatches(&operatorSharesIncreasedList, len(operatorSharesIncreasedList))
	return result.Error
}

func NewOperatorSharesIncreasedDB(db *gorm.DB) OperatorSharesIncreasedDB {
	return &operatorSharesIncreasedDB{gorm: db}
}
