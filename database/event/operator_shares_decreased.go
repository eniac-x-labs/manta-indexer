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

type OperatorSharesDecreasedView interface {
	QueryUnHandlerOperatorSharesDecreased() ([]OperatorSharesDecreased, error)
	GetOperatorSharesDecreased(string) (*OperatorSharesDecreased, error)
	ListOperatorSharesDecreased(page int, pageSize int, order string) ([]OperatorSharesDecreased, uint64)
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

func (osd operatorSharesDecreasedDB) GetOperatorSharesDecreased(address string) (*OperatorSharesDecreased, error) {
	var operatorSharesDecreased OperatorSharesDecreased
	result := osd.gorm.Table("operator_shares_decreased").Where(&OperatorSharesDecreased{Staker: common.HexToAddress(address)}).Take(&operatorSharesDecreased)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &operatorSharesDecreased, nil
}

func (osd operatorSharesDecreasedDB) ListOperatorSharesDecreased(page int, pageSize int, order string) ([]OperatorSharesDecreased, uint64) {
	var totalRecord int64
	var operatorSharesDecreasedList []OperatorSharesDecreased
	queryRoot := osd.gorm.Table("operator_shares_decreased")
	err := queryRoot.Select("guid").Count(&totalRecord).Error
	if err != nil {
		log.Error("list operatorSharesDecreasedDB count fail", "err", err)
	}

	queryRoot.Offset((page - 1) * pageSize).Limit(pageSize)
	if strings.ToLower(order) == "asc" {
		queryRoot.Order("guid asc")
	} else {
		queryRoot.Order("guid desc")
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
