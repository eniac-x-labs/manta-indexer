package event

import (
	"math/big"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum/common"

	_ "github.com/eniac-x-labs/manta-indexer/database/utils/serializers"
)

type StrategyDeposit struct {
	GUID       uuid.UUID      `gorm:"primaryKey" json:"guid"`
	BlockHash  common.Hash    `json:"block_hash" gorm:"serializer:bytes"`
	Number     *big.Int       `json:"number" gorm:"serializer:u256"`
	TxHash     common.Hash    `json:"tx_hash" gorm:"serializer:bytes"`
	Staker     common.Address `json:"staker" gorm:"serializer:bytes"`
	MantaToken common.Address `json:"manta_token" gorm:"serializer:bytes"`
	Strategy   common.Address `json:"strategy" gorm:"serializer:bytes"`
	Shares     *big.Int       `json:"shares" gorm:"serializer:u256"`
	IsHandle   uint8          `json:"is_handle"`
	Timestamp  uint64         `json:"timestamp"`
}

type StrategyDepositView interface {
	QueryStrategyDepositList(page int, pageSize int, order string) ([]StrategyDeposit, uint64)
}

type StrategyDepositDB interface {
	StrategyDepositView
	StoreStrategyDeposit([]StrategyDeposit) error
}

type strategyDepositDB struct {
	gorm *gorm.DB
}

func (db strategyDepositDB) QueryStrategyDepositList(page int, pageSize int, order string) ([]StrategyDeposit, uint64) {
	panic("implement me")
}

func (db strategyDepositDB) StoreStrategyDeposit(strategyDepositList []StrategyDeposit) error {
	result := db.gorm.CreateInBatches(&strategyDepositList, len(strategyDepositList))
	return result.Error
}

func NewStrategyDepositDB(db *gorm.DB) StrategyDepositDB {
	return &strategyDepositDB{gorm: db}
}
