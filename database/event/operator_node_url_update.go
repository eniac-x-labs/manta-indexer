package event

import (
	"math/big"

	_ "github.com/eniac-x-labs/manta-indexer/database/utils/serializers"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"gorm.io/gorm"
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

type OperatorNodeUrlUpdateView interface {
	QueryOperatorNodeUrlUpdateList(page int, pageSize int, order string) ([]OperatorNodeUrlUpdate, uint64)
}

type OperatorNodeUrlUpdateDB interface {
	OperatorNodeUrlUpdateView
	StoreOperatorNodeUrlUpdate([]OperatorNodeUrlUpdate) error
}

type operatorNodeUrlUpdateDB struct {
	gorm *gorm.DB
}

func (db operatorNodeUrlUpdateDB) QueryOperatorNodeUrlUpdateList(page int, pageSize int, order string) ([]OperatorNodeUrlUpdate, uint64) {
	panic("implement me")
}

func (db operatorNodeUrlUpdateDB) StoreOperatorNodeUrlUpdate(operatorNodeUrlUpdateList []OperatorNodeUrlUpdate) error {
	result := db.gorm.CreateInBatches(&operatorNodeUrlUpdateList, len(operatorNodeUrlUpdateList))
	return result.Error
}

func NewOperatorNodeUrlUpdateDB(db *gorm.DB) OperatorNodeUrlUpdateDB {
	return &operatorNodeUrlUpdateDB{gorm: db}
}
