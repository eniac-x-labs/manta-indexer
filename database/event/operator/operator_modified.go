package operator

import (
	"math/big"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum/common"

	_ "github.com/eniac-x-labs/manta-indexer/database/utils/serializers"
)

type OperatorModified struct {
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

func (OperatorModified) TableName() string {
	return "operator_modified"
}

type OperatorModifiedView interface {
	QueryOperatorModifiedList(page int, pageSize int, order string) ([]OperatorModified, uint64)
}

type OperatorModifiedDB interface {
	OperatorModifiedView
	StoreOperatorModified([]OperatorModified) error
}

type operatorModifiedDB struct {
	gorm *gorm.DB
}

func (db operatorModifiedDB) QueryOperatorModifiedList(page int, pageSize int, order string) ([]OperatorModified, uint64) {
	panic("implement me")
}

func (db operatorModifiedDB) StoreOperatorModified(operatorModifiedList []OperatorModified) error {
	result := db.gorm.Table("operator_modified").CreateInBatches(&operatorModifiedList, len(operatorModifiedList))
	return result.Error
}

func NewOperatorModifiedDB(db *gorm.DB) OperatorModifiedDB {
	return &operatorModifiedDB{gorm: db}
}
