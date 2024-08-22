package config

import (
	"reflect"

	"github.com/ethereum/go-ethereum/common"
)

const (
	DelegationManagerAddress = "0x602FD3f5f327B7e32DCEBa2d170C6460E1544F70"
	RewardManagerAddress     = "0x87FC620E6deADaa909a282F98B36933dAE05D516"
	StrategyManagerAddress   = "0x6076145fd9139180683085901540bd92fA43c52a"
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
