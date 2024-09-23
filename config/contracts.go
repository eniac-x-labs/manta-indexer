package config

import (
	"reflect"

	"github.com/ethereum/go-ethereum/common"
)

const (
	DelegationManagerAddress   = "0x8F847f86563873577f237191C44246D5408F1c55"
	RewardManagerAddress       = "0x36FB46E66283F42784d1Ce6aF636fbcC8f0cEE87"
	StrategyManagerAddress     = "0x99113e1989C3845e1Ad9DF768442Db24C1e1cC6e"
	MantaServiceManagerAddress = "0xA7918D253764E42d60C3ce2010a34d5a1e7C1398"
)

type MantaLayerContracts struct {
	DelegationManager   common.Address
	RewardManager       common.Address
	StrategyManager     common.Address
	MantaServiceManager common.Address
}

func ContractsFromConst() MantaLayerContracts {
	return MantaLayerContracts{
		DelegationManager:   common.HexToAddress(DelegationManagerAddress),
		RewardManager:       common.HexToAddress(RewardManagerAddress),
		StrategyManager:     common.HexToAddress(StrategyManagerAddress),
		MantaServiceManager: common.HexToAddress(MantaServiceManagerAddress),
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
