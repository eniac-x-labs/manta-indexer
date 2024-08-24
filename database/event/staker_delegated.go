package event

import (
	"github.com/ethereum/go-ethereum/log"
	"math/big"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum/common"

	_ "github.com/eniac-x-labs/manta-indexer/database/utils/serializers"
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

func (StakerDelegated) TableName() string {
	return "staker_delegated"
}

type StakerDelegatedView interface {
	ListStakerDelegated(address string, page int, pageSize int, order string) ([]StakerDelegated, uint64)
}

type StakerDelegatedDB interface {
	StakerDelegatedView
	StoreStakerDelegated([]StakerDelegated) error
}

type stakerDelegatedDB struct {
	gorm *gorm.DB
}

func (sd stakerDelegatedDB) ListStakerDelegated(address string, page int, pageSize int, order string) ([]StakerDelegated, uint64) {
	var totalRecord int64
	var stakerDelegatedList []StakerDelegated
	queryRoot := sd.gorm.Table("staker_delegated")
	if address != "0x00" {
		err := sd.gorm.Table("staker_delegated").Select("number").Where("operator = ?", address).Count(&totalRecord).Error
		if err != nil {
			log.Error("get staker delegated count fail")
		}
		queryRoot.Where("operator = ?", address).Offset((page - 1) * pageSize).Limit(pageSize)
	} else {
		err := sd.gorm.Table("staker_delegated").Select("number").Count(&totalRecord).Error
		if err != nil {
			log.Error("get staker delegated count fail ")
		}
		queryRoot.Offset((page - 1) * pageSize).Limit(pageSize)
	}
	if strings.ToLower(order) == "asc" {
		queryRoot.Order("number asc")
	} else {
		queryRoot.Order("number desc")
	}
	qErr := queryRoot.Find(&stakerDelegatedList).Error
	if qErr != nil {
		log.Error("list stakerDelegatedDB fail", "err", qErr)
	}
	return stakerDelegatedList, uint64(totalRecord)
}

func (db stakerDelegatedDB) StoreStakerDelegated(stakerDelegatedList []StakerDelegated) error {
	result := db.gorm.Table("staker_delegated").CreateInBatches(&stakerDelegatedList, len(stakerDelegatedList))
	return result.Error
}

func NewStakerDelegatedDB(db *gorm.DB) StakerDelegatedDB {
	return &stakerDelegatedDB{gorm: db}
}
