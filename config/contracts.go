package config

import (
	"reflect"

	"github.com/ethereum/go-ethereum/common"
)

const (
	DelegationManagerAddress = "0x319d48C4BBA7AB1a7C413aDa8DBe5a82cE09FD58"
	RewardManagerAddress     = "0x77fb31ebbB7312f5C3CaC0290aF0bda24287f3A8"
	StrategyManagerAddress   = "0x87b9a05800b0EDf9A4AC482C54A348bdd2070D28"
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
