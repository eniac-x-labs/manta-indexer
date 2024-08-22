package event

import (
	"math/big"

	_ "github.com/eniac-x-labs/manta-indexer/database/utils/serializers"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WithdrawalMigrated struct {
	GUID              uuid.UUID   `gorm:"primaryKey" json:"guid"`
	BlockHash         common.Hash `json:"block_hash" gorm:"serializer:bytes"`
	Number            *big.Int    `json:"number" gorm:"serializer:u256"`
	TxHash            common.Hash `json:"tx_hash" gorm:"serializer:bytes"`
	OldWithdrawalRoot common.Hash `json:"old_withdrawal_root" gorm:"serializer:bytes"`
	NewWithdrawalRoot common.Hash `json:"new_withdrawal_root" gorm:"serializer:bytes"`
	IsHandle          uint8       `json:"is_handle"`
	Timestamp         uint64      `json:"timestamp"`
}

type WithdrawalMigratedView interface {
	QueryWithdrawalMigratedList(page int, pageSize int, order string) ([]WithdrawalMigrated, uint64)
}

type WithdrawalMigratedDB interface {
	WithdrawalMigratedView
	StoreWithdrawalMigrated([]WithdrawalMigrated) error
}

type withdrawalMigratedDB struct {
	gorm *gorm.DB
}

func (db withdrawalMigratedDB) QueryWithdrawalMigratedList(page int, pageSize int, order string) ([]WithdrawalMigrated, uint64) {
	panic("implement me")
}

func (db withdrawalMigratedDB) StoreWithdrawalMigrated(withdrawalMigratedList []WithdrawalMigrated) error {
	result := db.gorm.CreateInBatches(&withdrawalMigratedList, len(withdrawalMigratedList))
	return result.Error
}

func NewWithdrawalMigratedDB(db *gorm.DB) WithdrawalMigratedDB {
	return &withdrawalMigratedDB{gorm: db}
}
