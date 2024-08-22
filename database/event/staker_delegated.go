package event

import (
	"math/big"

	_ "github.com/eniac-x-labs/manta-indexer/database/utils/serializers"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StakerDelegated struct {
	GUID      uuid.UUID      `gorm:"primaryKey" json:"guid"`
	BlockHash common.Hash    `json:"block_hash" gorm:"serializer:bytes"`
	Number    *big.Int       `json:"number" gorm:"serializer:u256"`
	TxHash    common.Hash    `json:"tx_hash" gorm:"serializer:bytes"`
	Operator  common.Address `json:"operator" gorm:"serializer:bytes"`
	Staker    common.Address `json:"staker" gorm:"serializer:bytes"`
	IsHandle  uint8          `json:"is_handle"`
	Timestamp uint64         `json:"timestamp"`
}

type StakerDelegatedView interface {
	QueryStakerDelegatedList(page int, pageSize int, order string) ([]StakerDelegated, uint64)
}

type StakerDelegatedDB interface {
	StakerDelegatedView
	StoreStakerDelegated([]StakerDelegated) error
}

type stakerDelegatedDB struct {
	gorm *gorm.DB
}

func (db stakerDelegatedDB) QueryStakerDelegatedList(page int, pageSize int, order string) ([]StakerDelegated, uint64) {
	panic("implement me")
}

func (db stakerDelegatedDB) StoreStakerDelegated(stakerDelegatedList []StakerDelegated) error {
	result := db.gorm.CreateInBatches(&stakerDelegatedList, len(stakerDelegatedList))
	return result.Error
}

func NewStakerDelegatedDB(db *gorm.DB) StakerDelegatedDB {
	return &stakerDelegatedDB{gorm: db}
}
