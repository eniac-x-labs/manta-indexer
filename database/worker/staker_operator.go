package worker

import (
	"errors"
	"strings"

	"gorm.io/gorm"
	"math/big"

	"github.com/google/uuid"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

type StakerOperator struct {
	GUID          uuid.UUID      `gorm:"primaryKey" json:"guid"`
	Staker        common.Address `gorm:"serializer:bytes" json:"staker"`
	Operator      common.Address `gorm:"serializer:bytes" json:"operator"`
	TotalStake    *big.Int       `gorm:"serializer:u256" json:"total_stake"`
	TotalReward   *big.Int       `gorm:"serializer:u256" json:"total_reward"`
	ClaimedAmount *big.Int       `gorm:"serializer:u256" json:"claimed_amount"`
	Timestamp     uint64
}

func (StakerOperator) TableName() string {
	return "staker_operator"
}

type StakerOperatorView interface {
	GetStakerOperator(string) (*StakerOperator, error)
	ListStakerOperator(operator string, staker string, page int, pageSize int, order string) ([]StakerOperator, uint64)
}

type StakerOperatorDB interface {
	StakerOperatorView
	QueryAndUpdateStakerOperator(stakeAddress string, operatorAddress string, mantaStake *big.Int) error
	StoreStakerOperator([]StakerOperator) error
}

type stakerOperatorDB struct {
	gorm *gorm.DB
}

func (sh *stakerOperatorDB) QueryAndUpdateStakerOperator(stakeAddress string, operatorAddress string, mantaStake *big.Int) error {
	var stakerOperator StakerOperator
	result := sh.gorm.Table("staker_operator").Where("staker = ? and operator = ?", strings.ToLower(stakeAddress), strings.ToLower(operatorAddress)).Take(&stakerOperator)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	stakerOperator.Operator = common.Address{}
	stakerOperator.TotalStake = big.NewInt(0)
	if mantaStake.Cmp(big.NewInt(0)) == 0 {
		stakerOperator.TotalStake = big.NewInt(0)
	} else {
		stakerOperator.TotalStake = new(big.Int).Add(stakerOperator.TotalStake, mantaStake)
	}
	err := sh.gorm.Save(stakerOperator).Error
	if err != nil {
		log.Error("update stake operator fail", "err", err)
		return err
	}
	return nil
}

func (shv stakerOperatorDB) GetStakerOperator(staker string) (*StakerOperator, error) {
	var stakerOperator StakerOperator
	result := shv.gorm.Table("staker_operator").Where("staker = ?", strings.ToLower(staker)).Take(&stakerOperator)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &stakerOperator, nil
}

func (shv stakerOperatorDB) ListStakerOperator(operator string, staker string, page int, pageSize int, order string) ([]StakerOperator, uint64) {
	var totalRecord int64
	var stakerOperatorList []StakerOperator
	queryRoot := shv.gorm.Table("staker_operator")
	if operator != "0x00" && staker != "0x00" {
		err := shv.gorm.Table("staker_operator").Select("total_stake").Where("operator = ? and staker = ?", operator, staker).Count(&totalRecord).Error
		if err != nil {
			log.Error("get stake operator count fail")
		}
		queryRoot.Where("operator = ? and staker = ?", operator, staker).Offset((page - 1) * pageSize).Limit(pageSize)
	} else {
		err := shv.gorm.Table("staker_operator").Select("total_stake").Count(&totalRecord).Error
		if err != nil {
			log.Error("get stake operator count fail ")
		}
		queryRoot.Offset((page - 1) * pageSize).Limit(pageSize)
	}
	if strings.ToLower(order) == "asc" {
		queryRoot.Order("timestamp asc")
	} else {
		queryRoot.Order("timestamp desc")
	}
	qErr := queryRoot.Find(&stakerOperatorList).Error
	if qErr != nil {
		log.Error("list stake operator db fail", "err", qErr)
	}
	return stakerOperatorList, uint64(totalRecord)
}

func NewStakerOperatorDB(db *gorm.DB) StakerOperatorDB {
	return &stakerOperatorDB{gorm: db}
}

func (sh *stakerOperatorDB) StoreStakerOperator(stakerOperators []StakerOperator) error {
	result := sh.gorm.Table("staker_operator").CreateInBatches(&stakerOperators, len(stakerOperators))
	return result.Error
}
