package event

import (
	"math/big"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"

	_ "github.com/eniac-x-labs/manta-indexer/database/utils/serializers"
)

type Strategies struct {
	GUID      uuid.UUID      `gorm:"primaryKey" json:"guid"`
	BlockHash common.Hash    `json:"block_hash" gorm:"serializer:bytes"`
	Number    *big.Int       `json:"number" gorm:"serializer:u256"`
	TxHash    common.Hash    `json:"tx_hash" gorm:"serializer:bytes"`
	Strategy  common.Address `json:"strategy" gorm:"serializer:bytes"`
	IsHandle  uint8          `json:"is_handle"`
	Timestamp uint64         `json:"timestamp"`
}

func (Strategies) TableName() string {
	return "strategies"
}

type StrategiesView interface {
	QueryStrategiesList(page int, pageSize int, order string) ([]Strategies, uint64)
}

type StrategiesDB interface {
	StrategiesView
	RemoveStoreStrategies([]Strategies) error
	StoreStrategies([]Strategies) error
}

type strategiesDB struct {
	gorm *gorm.DB
}

func (db strategiesDB) RemoveStoreStrategies(strategies []Strategies) error {
	for _, v := range strategies {
		db.gorm.Table("strategies").Delete(v.Strategy)
	}
	return nil
}

func (db strategiesDB) QueryStrategiesList(page int, pageSize int, order string) ([]Strategies, uint64) {
	var totalRecord int64
	var strategyList []Strategies
	queryStateRoot := db.gorm.Table("operator_registered")
	err := db.gorm.Table("operator_registered").Select("number").Count(&totalRecord).Error
	if err != nil {
		log.Error("get strategies fail", "err", err)
	}
	queryStateRoot.Offset((page - 1) * pageSize).Limit(pageSize)
	if strings.ToLower(order) == "asc" {
		queryStateRoot.Order("number asc")
	} else {
		queryStateRoot.Order("number desc")
	}
	qErr := queryStateRoot.Find(&strategyList).Error
	if qErr != nil {
		log.Error("get strategies fail", "err", qErr)
	}
	return strategyList, uint64(totalRecord)
}

func (db strategiesDB) StoreStrategies(strategiesList []Strategies) error {
	result := db.gorm.Table("strategies").CreateInBatches(&strategiesList, len(strategiesList))
	return result.Error
}

func NewStrategiesDB(db *gorm.DB) StrategiesDB {
	return &strategiesDB{gorm: db}
}
