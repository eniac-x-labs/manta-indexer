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

type OperatorSharesDecreased struct {
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

func (OperatorSharesDecreased) TableName() string {
	return "operator_shares_decreased"
}

type OperatorSharesDecreasedView interface {
	QueryUnHandlerOperatorSharesDecreased() ([]OperatorSharesDecreased, error)
	ListOperatorSharesDecreased(address string, page int, pageSize int, order string) ([]OperatorSharesDecreased, uint64)
}

type OperatorSharesDecreasedDB interface {
	OperatorSharesDecreasedView
	MarkedOperatorSharesDecreasedHandled(unHandleOperatorSharesDecreased []OperatorSharesDecreased) error
	StoreOperatorSharesDecreased([]OperatorSharesDecreased) error
}

type operatorSharesDecreasedDB struct {
	gorm *gorm.DB
}

func (osd operatorSharesDecreasedDB) MarkedOperatorSharesDecreasedHandled(unHandleOperatorSharesDecreased []OperatorSharesDecreased) error {
	for i := 0; i < len(unHandleOperatorSharesDecreased); i++ {
		var operatorSharesDecreased = OperatorSharesDecreased{}
		result := osd.gorm.Table("operator_shares_decreased").Where(&OperatorRegistered{GUID: unHandleOperatorSharesDecreased[i].GUID}).Take(&operatorSharesDecreased)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return nil
			}
			return result.Error
		}
		operatorSharesDecreased.IsHandle = 1
		err := osd.gorm.Save(operatorSharesDecreased).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (osd operatorSharesDecreasedDB) QueryUnHandlerOperatorSharesDecreased() ([]OperatorSharesDecreased, error) {
	var operatorSharesDecreasedList []OperatorSharesDecreased
	err := osd.gorm.Table("operator_shares_decreased").Where("is_handle = ?", 0).Find(&operatorSharesDecreasedList).Error
	if err != nil {
		log.Error("get operator share decrease list fail", "err", err)
		return nil, err
	}
	return operatorSharesDecreasedList, nil
}

func (osd operatorSharesDecreasedDB) ListOperatorSharesDecreased(address string, page int, pageSize int, order string) ([]OperatorSharesDecreased, uint64) {
	var totalRecord int64
	var operatorSharesDecreasedList []OperatorSharesDecreased
	queryRoot := osd.gorm.Table("operator_shares_decreased")
	if address != "0x00" {
		err := osd.gorm.Table("operator_shares_decreased").Select("number").Where("operator = ?", address).Count(&totalRecord).Error
		if err != nil {
			log.Error("get operator share decreased count fail")
		}
		queryRoot.Where("operator = ?", address).Offset((page - 1) * pageSize).Limit(pageSize)
	} else {
		err := osd.gorm.Table("operator_shares_decreased").Select("number").Count(&totalRecord).Error
		if err != nil {
			log.Error("get operator share decreased count fail ")
		}
		queryRoot.Offset((page - 1) * pageSize).Limit(pageSize)
	}
	if strings.ToLower(order) == "asc" {
		queryRoot.Order("number asc")
	} else {
		queryRoot.Order("number desc")
	}
	qErr := queryRoot.Find(&operatorSharesDecreasedList).Error
	if qErr != nil {
		log.Error("list operatorSharesDecreasedDB fail", "err", qErr)
	}
	return operatorSharesDecreasedList, uint64(totalRecord)
}

func (osd operatorSharesDecreasedDB) StoreOperatorSharesDecreased(operatorSharesDecreasedList []OperatorSharesDecreased) error {
	result := osd.gorm.Table("operator_shares_decreased").CreateInBatches(&operatorSharesDecreasedList, len(operatorSharesDecreasedList))
	return result.Error
}

func NewOperatorSharesDecreasedDB(db *gorm.DB) OperatorSharesDecreasedDB {
	return &operatorSharesDecreasedDB{gorm: db}
}
