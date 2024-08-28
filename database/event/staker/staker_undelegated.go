package staker

import (
	"errors"
	"math/big"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"

	_ "github.com/eniac-x-labs/manta-indexer/database/utils/serializers"
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

func (StakerUndelegated) TableName() string {
	return "staker_undelegated"
}

type StakerUndelegatedView interface {
	QueryUnHandleStakerUndelegated() ([]StakerUndelegated, error)
	ListStakerUndelegated(address string, page int, pageSize int, order string) ([]StakerUndelegated, uint64)
}

type StakerUndelegatedDB interface {
	StakerUndelegatedView
	MarkedStakerUnDelegated(stakerUnDelegatedList []StakerUndelegated) error
	StoreStakerUndelegated([]StakerUndelegated) error
}

type stakerUndelegatedDB struct {
	gorm *gorm.DB
}

func (su stakerUndelegatedDB) QueryUnHandleStakerUndelegated() ([]StakerUndelegated, error) {
	var stakerUndelegatedList []StakerUndelegated
	err := su.gorm.Table("staker_undelegated").Where("is_handle = ?", 0).Find(&stakerUndelegatedList).Error
	if err != nil {
		log.Error("get strategy undelegated fail", "err", err)
		return nil, err
	}
	return stakerUndelegatedList, nil
}

func (su stakerUndelegatedDB) MarkedStakerUnDelegated(stakerUnDelegatedList []StakerUndelegated) error {
	for i := 0; i < len(stakerUnDelegatedList); i++ {
		var stakerUndelegated = StakerUndelegated{}
		result := su.gorm.Table("staker_undelegated").Where("guid = ?", stakerUnDelegatedList[i].GUID).Take(&stakerUndelegated)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return nil
			}
			return result.Error
		}
		stakerUndelegated.IsHandle = 1
		err := su.gorm.Save(stakerUndelegated).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (su stakerUndelegatedDB) ListStakerUndelegated(address string, page int, pageSize int, order string) ([]StakerUndelegated, uint64) {
	var totalRecord int64
	var stakerUndelegatedList []StakerUndelegated
	queryRoot := su.gorm.Table("staker_undelegated")
	if address != "0x00" {
		err := su.gorm.Table("staker_undelegated").Select("number").Where("staker = ?", address).Count(&totalRecord).Error
		if err != nil {
			log.Error("get staker undelegated count fail")
		}
		queryRoot.Where("staker = ?", address).Offset((page - 1) * pageSize).Limit(pageSize)
	} else {
		err := su.gorm.Table("staker_undelegated").Select("number").Count(&totalRecord).Error
		if err != nil {
			log.Error("get staker undelegated count fail ")
		}
		queryRoot.Offset((page - 1) * pageSize).Limit(pageSize)
	}
	if strings.ToLower(order) == "asc" {
		queryRoot.Order("number asc")
	} else {
		queryRoot.Order("number desc")
	}
	qErr := queryRoot.Find(&stakerUndelegatedList).Error
	if qErr != nil {
		log.Error("list stakerUndelegatedDB fail", "err", qErr)
	}
	return stakerUndelegatedList, uint64(totalRecord)
}

func (db stakerUndelegatedDB) StoreStakerUndelegated(stakerUndelegatedList []StakerUndelegated) error {
	result := db.gorm.Table("staker_undelegated").CreateInBatches(&stakerUndelegatedList, len(stakerUndelegatedList))
	return result.Error
}

func NewStakerUndelegatedDB(db *gorm.DB) StakerUndelegatedDB {
	return &stakerUndelegatedDB{gorm: db}
}
