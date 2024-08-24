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

type OperatorNodeUrlUpdate struct {
	GUID        uuid.UUID      `gorm:"primaryKey" json:"guid"`
	BlockHash   common.Hash    `json:"block_hash" gorm:"serializer:bytes"`
	Number      *big.Int       `json:"number" gorm:"serializer:u256"`
	TxHash      common.Hash    `json:"tx_hash" gorm:"serializer:bytes"`
	Operator    common.Address `json:"operator" gorm:"serializer:bytes"`
	MetadataUri string         `json:"metadata_uri"`
	IsHandle    uint8          `json:"is_handle"`
	Timestamp   uint64         `json:"timestamp"`
}

func (OperatorNodeUrlUpdate) TableName() string {
	return "operator_node_url_update"
}

type OperatorNodeUrlUpdateView interface {
	QueryUnHandleOperatorNodeUrlUpdate() ([]OperatorNodeUrlUpdate, error)
	ListOperatorNodeUrlUpdate(address string, page int, pageSize int, order string) ([]OperatorNodeUrlUpdate, uint64)
}

type OperatorNodeUrlUpdateDB interface {
	OperatorNodeUrlUpdateView
	MarkedOperatorNodeUrlUpdateHandled([]OperatorNodeUrlUpdate) error
	StoreOperatorNodeUrlUpdate(operatorNodeUrlUpdateList []OperatorNodeUrlUpdate) error
}

type operatorNodeUrlUpdateDB struct {
	gorm *gorm.DB
}

func (onuu operatorNodeUrlUpdateDB) ListOperatorNodeUrlUpdate(address string, page int, pageSize int, order string) ([]OperatorNodeUrlUpdate, uint64) {
	address = strings.ToLower(address)

	var totalRecord int64
	var operatorNodeUrlUpdateList []OperatorNodeUrlUpdate
	queryRoot := onuu.gorm.Table("operator_node_url_update")
	if address != "0x00" {
		err := onuu.gorm.Table("operator_node_url_update").Select("number").Where("operator = ?", address).Count(&totalRecord).Error
		if err != nil {
			log.Error("get operator node url update count fail")
		}
		queryRoot.Where("operator = ?", address).Offset((page - 1) * pageSize).Limit(pageSize)
	} else {
		err := onuu.gorm.Table("operator_node_url_update").Select("number").Count(&totalRecord).Error
		if err != nil {
			log.Error("get operator node url update count fail ")
		}
		queryRoot.Offset((page - 1) * pageSize).Limit(pageSize)
	}
	if strings.ToLower(order) == "asc" {
		queryRoot.Order("number asc")
	} else {
		queryRoot.Order("number desc")
	}
	qErr := queryRoot.Find(&operatorNodeUrlUpdateList).Error
	if qErr != nil {
		log.Error("list operator node url update db fail", "err", qErr)
	}
	return operatorNodeUrlUpdateList, uint64(totalRecord)
}

func (onuu operatorNodeUrlUpdateDB) MarkedOperatorNodeUrlUpdateHandled(urlUpdates []OperatorNodeUrlUpdate) error {
	for i := 0; i < len(urlUpdates); i++ {
		var operatorNodeUrlUpdates = OperatorNodeUrlUpdate{}
		result := onuu.gorm.Where(&OperatorRegistered{GUID: urlUpdates[i].GUID}).Take(&operatorNodeUrlUpdates)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return nil
			}
			return result.Error
		}
		operatorNodeUrlUpdates.IsHandle = 1
		err := onuu.gorm.Save(operatorNodeUrlUpdates).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (onuu operatorNodeUrlUpdateDB) QueryUnHandleOperatorNodeUrlUpdate() ([]OperatorNodeUrlUpdate, error) {
	var operatorNodeUrlUpdateList []OperatorNodeUrlUpdate
	err := onuu.gorm.Table("operator_node_url_update").Where("is_handle = ?", 0).Find(&operatorNodeUrlUpdateList).Error
	if err != nil {
		log.Error("get unhandled operator node url fail", "err", err)
		return nil, err
	}
	return operatorNodeUrlUpdateList, nil
}

func (onuu operatorNodeUrlUpdateDB) StoreOperatorNodeUrlUpdate(operatorNodeUrlUpdateList []OperatorNodeUrlUpdate) error {
	result := onuu.gorm.Table("operator_node_url_update").CreateInBatches(&operatorNodeUrlUpdateList, len(operatorNodeUrlUpdateList))
	return result.Error
}

func NewOperatorNodeUrlUpdateDB(db *gorm.DB) OperatorNodeUrlUpdateDB {
	return &operatorNodeUrlUpdateDB{gorm: db}
}
