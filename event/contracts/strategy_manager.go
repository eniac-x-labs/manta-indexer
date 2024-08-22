package contracts

import (
	"math/big"

	common2 "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"

	"github.com/eniac-x-labs/manta-indexer/bindings/sm"
	"github.com/eniac-x-labs/manta-indexer/config"
	"github.com/eniac-x-labs/manta-indexer/database"
	"github.com/eniac-x-labs/manta-indexer/database/event"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

type StrategyManager struct {
	db       *database.DB
	SmAbi    *abi.ABI
	SmFilter *sm.StrategyManagerFilterer
}

func NewStrategyManager(db *database.DB) (*StrategyManager, error) {
	strategyAbi, err := sm.StrategyManagerMetaData.GetAbi()
	if err != nil {
		log.Error("get delegate manager abi fail", "err", err)
		return nil, err
	}

	strategyUnpack, err := sm.NewStrategyManagerFilterer(common2.Address{}, nil)
	if err != nil {
		log.Error("new delegation manager fail", "err", err)
		return nil, err
	}

	return &StrategyManager{
		db:       db,
		SmAbi:    strategyAbi,
		SmFilter: strategyUnpack,
	}, nil
}

func (sm *StrategyManager) ProcessStrategyManager(fromHeight *big.Int, toHeight *big.Int) error {
	contractEventFilter := event.ContractEvent{ContractAddress: common2.HexToAddress(config.StrategyManagerAddress)}
	contractEventList, err := sm.db.ContractEvent.ContractEventsWithFilter(contractEventFilter, fromHeight, toHeight)
	if err != nil {
		log.Error("get contracts event list fail", "err", err)
		return err
	}
	for _, eventItem := range contractEventList {
		// emit Deposit(staker, weth, strategy, shares);
		rlpLog := eventItem.RLPLog

		// Deposit
		if eventItem.EventSignature.String() == sm.SmAbi.Events["Deposit"].ID.String() {
			depositEvent, err := sm.SmFilter.ParseDeposit(*rlpLog)
			if err != nil {
				log.Error("parse deposit event fail", "err", err)
				return err
			}

			log.Info("parse deposit event success",
				"staker", depositEvent.Staker.String(),
				// "mantaToken", depositEvent.MantaToken.String(),
				"strategy", depositEvent.Strategy.String(),
				"shares", depositEvent.Shares.String())
		}
	}

	return nil
}
