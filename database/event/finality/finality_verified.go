package finality

import (
	"errors"
	"math/big"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum/common"

	_ "github.com/eniac-x-labs/manta-indexer/database/utils/serializers"
)

type FinalityVerified struct {
	GUID          uuid.UUID      `gorm:"primaryKey" json:"guid"`
	Proposer      common.Address `json:"proposer" gorm:"serializer:bytes"`
	TxBlockNumber *big.Int       `json:"tx_block_number" gorm:"serializer:u256"`
	L1BlockNumber *big.Int       `json:"L1_block_number" gorm:"serializer:u256"`
	L2BlockNumber *big.Int       `json:"L2_block_number" gorm:"serializer:u256"`
	L1BlockHash   string         `json:"L1_block_hash" gorm:"serializer:bytes"`
	OutputRoot    string         `json:"output_root"`
	Timestamp     uint64         `json:"timestamp"`
}

func (FinalityVerified) TableName() string {
	return "finality_verified"
}

type FinalityVerifiedView interface {
	GetFinalityVerifiedByL2Block(*big.Int) (*FinalityVerified, error)
}

type FinalityVerifiedDB interface {
	FinalityVerifiedView
	StoreFinalityVerified([]FinalityVerified) error
}

type finalityVerifiedDB struct {
	gorm *gorm.DB
}

func NewFinalityVerifiedDB(db *gorm.DB) FinalityVerifiedDB {
	return &finalityVerifiedDB{gorm: db}
}

func (fv finalityVerifiedDB) StoreFinalityVerified(finalityVerifiedList []FinalityVerified) error {
	result := fv.gorm.Table("finality_verified").CreateInBatches(&finalityVerifiedList, len(finalityVerifiedList))
	return result.Error
}

func (fv finalityVerifiedDB) GetFinalityVerifiedByL2Block(l2BlockNumber *big.Int) (*FinalityVerified, error) {
	var finalityVerified FinalityVerified
	result := fv.gorm.Table("finality_verified").Where("L2_block_number = ?", l2BlockNumber.String()).Take(&finalityVerified)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &finalityVerified, nil
}
