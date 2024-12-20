package staker

import (
	"errors"
	"strings"

	"math/big"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"

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

func (StrategyDeposit) TableName() string {
	return "strategy_deposit"
}

type StrategyDepositView interface {
	QueryUnHandleStrategyDeposit() ([]StrategyDeposit, error)

	ListStrategyDeposit(address string, page int, pageSize int, order string) ([]StrategyDeposit, uint64)
}

type StrategyDepositDB interface {
	StrategyDepositView
	MarkedStrategyDepositHandled([]StrategyDeposit) error
	StoreStrategyDeposit([]StrategyDeposit) error
}

type strategyDepositDB struct {
	gorm *gorm.DB
}

func (sd strategyDepositDB) QueryUnHandleStrategyDeposit() ([]StrategyDeposit, error) {
	var strategyDepositList []StrategyDeposit
	err := sd.gorm.Table("strategy_deposit").Where("is_handle = ?", 0).Find(&strategyDepositList).Error
	if err != nil {
		log.Error("get strategy deposit fail", "err", err)
		return nil, err
	}
	return strategyDepositList, nil
}

func (sd strategyDepositDB) MarkedStrategyDepositHandled(strategyDeposits []StrategyDeposit) error {
	for i := 0; i < len(strategyDeposits); i++ {
		var strategyDeposit = StrategyDeposit{}
		result := sd.gorm.Table("strategy_deposit").Where("guid = ?", strategyDeposits[i].GUID).Take(&strategyDeposit)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return nil
			}
			return result.Error
		}
		strategyDeposit.IsHandle = 1
		err := sd.gorm.Save(strategyDeposit).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (sdv strategyDepositDB) ListStrategyDeposit(address string, page int, pageSize int, order string) ([]StrategyDeposit, uint64) {
	var totalRecord int64
	var strategyDepositList []StrategyDeposit
	queryRoot := sdv.gorm.Table("strategy_deposit")
	if address != "0x00" {
		err := sdv.gorm.Table("strategy_deposit").Select("number").Where("staker = ?", address).Count(&totalRecord).Error
		if err != nil {
			log.Error("get list strategy count fail")
		}
		queryRoot.Where("staker = ?", address).Offset((page - 1) * pageSize).Limit(pageSize)
	} else {
		err := sdv.gorm.Table("strategy_deposit").Select("number").Count(&totalRecord).Error
		if err != nil {
			log.Error("get list strategy count fail ")
		}
		queryRoot.Offset((page - 1) * pageSize).Limit(pageSize)
	}
	if strings.ToLower(order) == "asc" {
		queryRoot.Order("number asc")
	} else {
		queryRoot.Order("number desc")
	}
	qErr := queryRoot.Find(&strategyDepositList).Error
	if qErr != nil {
		log.Error("list strategy deposit database fail", "err", qErr)
	}
	return strategyDepositList, uint64(totalRecord)
}

func (sd strategyDepositDB) StoreStrategyDeposit(strategyDepositList []StrategyDeposit) error {
	result := sd.gorm.Table("strategy_deposit").CreateInBatches(&strategyDepositList, len(strategyDepositList))
	return result.Error
}

func NewStrategyDepositDB(db *gorm.DB) StrategyDepositDB {
	return &strategyDepositDB{gorm: db}
}
