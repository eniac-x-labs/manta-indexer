package event

import (
	"math/big"

	_ "github.com/eniac-x-labs/manta-indexer/database/utils/serializers"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StakerUndelegated struct {
	GUID      uuid.UUID      `gorm:"primaryKey" json:"guid"`
	BlockHash common.Hash    `json:"block_hash" gorm:"serializer:bytes"`
	Number    *big.Int       `json:"number" gorm:"serializer:u256"`
	TxHash    common.Hash    `json:"tx_hash" gorm:"serializer:bytes"`
	Operator  common.Address `json:"operator" gorm:"serializer:bytes"`
	Staker    common.Address `json:"staker" gorm:"serializer:bytes"`
	IsHandle  uint8          `json:"is_handle"`
	Timestamp uint64         `json:"timestamp"`
}

type StakerUndelegatedView interface {
	QueryStakerUndelegatedList(page int, pageSize int, order string) ([]StakerUndelegated, uint64)
}

type StakerUndelegatedDB interface {
	StakerUndelegatedView
	StoreStakerUndelegated([]StakerUndelegated) error
}

type stakerUndelegatedDB struct {
	gorm *gorm.DB
}

func (db stakerUndelegatedDB) QueryStakerUndelegatedList(page int, pageSize int, order string) ([]StakerUndelegated, uint64) {
	panic("implement me")
}

func (db stakerUndelegatedDB) StoreStakerUndelegated(stakerUndelegatedList []StakerUndelegated) error {
	result := db.gorm.CreateInBatches(&stakerUndelegatedList, len(stakerUndelegatedList))
	return result.Error
}

func NewStakerUndelegatedDB(db *gorm.DB) StakerUndelegatedDB {
	return &stakerUndelegatedDB{gorm: db}
}
