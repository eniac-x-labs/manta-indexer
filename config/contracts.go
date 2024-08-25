package config

import (
	"reflect"

	"github.com/ethereum/go-ethereum/common"
)

const (
	DelegationManagerAddress = "0xCef760344c1e50FCF9cd931fe9fa6e3d84EacF94"
	RewardManagerAddress     = "0xE2425ba98b51Cbf07076DD5308D55954b43eb3a9"
	StrategyManagerAddress   = "0x231Ac717723Bd14C721A2a551E8d81ADb5831250"
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
