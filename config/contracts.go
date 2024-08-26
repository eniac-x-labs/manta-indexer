package config

import (
	"reflect"

	"github.com/ethereum/go-ethereum/common"
)

const (
	DelegationManagerAddress = "0x89471CCB7AC32C8526132A49F9abA6c4cA0595D0"
	RewardManagerAddress     = "0xFBE94848630F5C3C84eF767be585e30f478b282E"
	StrategyManagerAddress   = "0x44654787F728B01E00B5284632974a928431eAE0"
)

type MantaLayerContracts struct {
	DelegationManager common.Address
	RewardManager     common.Address
	StrategyManager   common.Address
}

func ContractsFromConst() MantaLayerContracts {
	return MantaLayerContracts{
		DelegationManager: common.HexToAddress(DelegationManagerAddress),
		RewardManager:     common.HexToAddress(RewardManagerAddress),
		StrategyManager:   common.HexToAddress(StrategyManagerAddress),
	}
}

func (mc MantaLayerContracts) ForEach(cb func(string, common.Address) error) error {
	contracts := reflect.ValueOf(mc)
	fields := reflect.VisibleFields(reflect.TypeOf(mc))
	for _, field := range fields {
		addr := (contracts.FieldByName(field.Name).Interface()).(common.Address)
		if err := cb(field.Name, addr); err != nil {
			return err
		}
	}
	return nil
}
