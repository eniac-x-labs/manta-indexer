package worker

import (
	"errors"
	"math/big"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

type Operators struct {
	GUID                     uuid.UUID      `gorm:"primaryKey"`
	BlockHash                common.Hash    `gorm:"serializer:bytes"`
	Number                   *big.Int       `gorm:"serializer:u256"`
	TxHash                   common.Hash    `gorm:"serializer:bytes"`
	Operator                 common.Address `gorm:"serializer:bytes"`
	Socket                   string         `json:"socket"`
	EarningsReceiver         common.Address `gorm:"serializer:bytes"`
	DelegationApprover       common.Address `gorm:"serializer:bytes"`
	StakerOptoutWindowBlocks *big.Int       `gorm:"serializer:u256"`
	TotalMantaStake          *big.Int       `gorm:"serializer:u256"`
	TotalStakeReward         *big.Int       `gorm:"serializer:u256"`
	RateReturn               string         `json:"rate_return"`
	Status                   uint8
	Timestamp                uint64
}

func (Operators) TableName() string {
	return "operators"
}

type OperatorsView interface {
	QueryAndUpdateOperator(operator common.Address, opType OperatorsType) error
	GetOperator(string) (*Operators, error)
	ListOperator(page int, pageSize int, order string) ([]Operators, uint64)
}

type OperatorsDB interface {
	OperatorsView
	StoreOperators([]Operators) error
}

type operatorsDB struct {
	gorm *gorm.DB
}

func (op *operatorsDB) QueryAndUpdateOperator(operator common.Address, opType OperatorsType) error {
	var operatorEntity Operators
	result := op.gorm.Where(&Operators{Operator: operator}).Take(&operatorEntity)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil
		}
		return result.Error
	}
	var zeroAddress common.Address
	if opType.Socket != "" {
		operatorEntity.Socket = opType.Socket
	}
	if opType.EarningsReceiver != zeroAddress {
		operatorEntity.EarningsReceiver = opType.EarningsReceiver
	}
	if opType.DelegationApprover != zeroAddress {
		operatorEntity.DelegationApprover = opType.DelegationApprover
	}
	if opType.StakerOptoutWindowBlocks != nil {
		operatorEntity.StakerOptoutWindowBlocks = opType.StakerOptoutWindowBlocks
	}
	if opType.TotalMantaStake != nil {
		totalStake := new(big.Int).Add(operatorEntity.TotalMantaStake, opType.TotalMantaStake)
		operatorEntity.TotalMantaStake = totalStake
		operatorEntity.RateReturn = new(big.Int).Div(operatorEntity.TotalStakeReward, totalStake).String()
	}
	if opType.TotalStakeReward != nil {
		totalStakeReward := new(big.Int).Add(operatorEntity.TotalStakeReward, opType.TotalStakeReward)
		operatorEntity.TotalStakeReward = totalStakeReward
		operatorEntity.RateReturn = new(big.Int).Div(totalStakeReward, operatorEntity.TotalMantaStake).String()
	}
	if operatorEntity.Status != 0 {
		operatorEntity.Status = opType.Status
	}
	err := op.gorm.Save(operatorEntity).Error
	if err != nil {
		log.Error("Update node url fail", "err", err)
		return err
	}
	return nil
}

func (ov operatorsDB) GetOperator(operator string) (*Operators, error) {
	var operatorDetail Operators
	result := ov.gorm.Where(&Operators{Operator: common.HexToAddress(operator)}).Take(&operatorDetail)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &operatorDetail, nil
}

func (ov operatorsDB) ListOperator(page int, pageSize int, order string) ([]Operators, uint64) {
	var totalRecord int64
	var operatorList []Operators
	queryRoot := ov.gorm.Table("operators")
	err := queryRoot.Select("guid").Count(&totalRecord).Error
	if err != nil {
		log.Error("list operatorsDB count fail", "err", err)
	}

	queryRoot.Offset((page - 1) * pageSize).Limit(pageSize)
	if strings.ToLower(order) == "asc" {
		queryRoot.Order("guid asc")
	} else {
		queryRoot.Order("guid desc")
	}
	qErr := queryRoot.Find(&operatorList).Error
	if qErr != nil {
		log.Error("list operatorsDB fail", "err", qErr)
	}
	return operatorList, uint64(totalRecord)
}

func NewOperatorsDB(db *gorm.DB) OperatorsDB {
	return &operatorsDB{gorm: db}
}

func (op *operatorsDB) StoreOperators(events []Operators) error {
	result := op.gorm.CreateInBatches(&events, len(events))
	return result.Error
}
