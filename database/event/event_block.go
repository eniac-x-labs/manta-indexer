package event

import (
	"errors"
	"gorm.io/gorm"
	"math/big"

	"github.com/google/uuid"

	"github.com/ethereum/go-ethereum/common"

	common2 "github.com/eniac-x-labs/manta-indexer/database/common"
	_ "github.com/eniac-x-labs/manta-indexer/database/utils/serializers"
)

type EventBlocks struct {
	GUID       uuid.UUID   `gorm:"primaryKey"`
	Hash       common.Hash `gorm:"serializer:bytes"`
	ParentHash common.Hash `gorm:"serializer:bytes"`
	Number     *big.Int    `gorm:"serializer:u256"`
	Timestamp  uint64
}

type BlocksView interface {
	LatestEventBlockHeader() (*common2.BlockHeader, error)
}

type EventBlocksDB interface {
	BlocksView
	StoreEventBlocks([]EventBlocks) error
}

type evnetBlocksDB struct {
	gorm *gorm.DB
}

func (e evnetBlocksDB) LatestEventBlockHeader() (*common2.BlockHeader, error) {
	eventQuery := e.gorm.Where("number = (?)", e.gorm.Table("event_blocks").Select("MAX(number)"))
	var header common2.BlockHeader
	result := eventQuery.Take(&header)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &header, nil
}

func (e evnetBlocksDB) StoreEventBlocks(eventBlocks []EventBlocks) error {
	result := e.gorm.CreateInBatches(&eventBlocks, len(eventBlocks))
	return result.Error
}

func NewEventBlocksDB(db *gorm.DB) EventBlocksDB {
	return &evnetBlocksDB{gorm: db}
}
