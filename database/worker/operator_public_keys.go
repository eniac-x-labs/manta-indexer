package worker

import (
	"gorm.io/gorm"
	"math/big"

	"github.com/google/uuid"

	"github.com/ethereum/go-ethereum/common"
)

type OperatorPublicKeys struct {
	GUID       uuid.UUID      `gorm:"primaryKey"`
	Operator   common.Address `gorm:"serializer:bytes"`
	PubkeyHash common.Hash    `gorm:"serializer:bytes"`
	PubkeyG1   *big.Int       `gorm:"serializer:u256"`
	PubkeyG2   *big.Int       `gorm:"serializer:u256"`
	Timestamp  uint64
}

type OperatorPublicKeysView interface {
}

type OperatorPublicKeysDB interface {
	OperatorPublicKeysView
	StoreOperatorPublicKeys([]OperatorPublicKeys) error
}

type operatorPublicKeysDB struct {
	gorm *gorm.DB
}

func NewOperatorPublicKeysDB(db *gorm.DB) OperatorPublicKeysDB {
	return &operatorPublicKeysDB{gorm: db}
}

func (db *operatorPublicKeysDB) StoreOperatorPublicKeys(opPubKeys []OperatorPublicKeys) error {
	result := db.gorm.CreateInBatches(&opPubKeys, len(opPubKeys))
	return result.Error
}
