package contracts

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	common2 "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/google/uuid"

	"github.com/eniac-x-labs/manta-indexer/bindings/sm"
	"github.com/eniac-x-labs/manta-indexer/config"
	"github.com/eniac-x-labs/manta-indexer/database"
	"github.com/eniac-x-labs/manta-indexer/database/event"
	"github.com/eniac-x-labs/manta-indexer/synchronizer/retry"
)

type StrategyManager struct {
	db       *database.DB
	SmAbi    *abi.ABI
	SmFilter *sm.StrategyManagerFilterer
	smCtx    context.Context
}

func NewStrategyManager(db *database.DB) (*StrategyManager, error) {
	strategyAbi, err := sm.StrategyManagerMetaData.GetAbi()
	if err != nil {
		log.Error("get strategy manager abi fail", "err", err)
		return nil, err
	}

	strategyUnpack, err := sm.NewStrategyManagerFilterer(common2.Address{}, nil)
	if err != nil {
		log.Error("new strategy manager fail", "err", err)
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

	deposits := make([]event.StrategyDeposit, 0, len(contractEventList))

	for _, eventItem := range contractEventList {
		rlpLog := eventItem.RLPLog

		header, err := sm.db.Blocks.BlockHeader(eventItem.BlockHash)
		if err != nil {
			log.Error("ProcessStrategyManager db Blocks BlockHeader by BlockHash fail", "err", err)
			return err
		}

		// emit Deposit(staker, weth, strategy, shares);
		if eventItem.EventSignature.String() == sm.SmAbi.Events["Deposit"].ID.String() {
			depositEvent, err := sm.SmFilter.ParseDeposit(*rlpLog)
			if err != nil {
				log.Error("parse deposit event fail", "err", err)
				return err
			}
			log.Info("parse deposit event success",
				"staker", depositEvent.Staker.String(),
				"strategy", depositEvent.Strategy.String(),
				"shares", depositEvent.Shares.String())

			deposit := event.StrategyDeposit{
				GUID:       uuid.New(),
				BlockHash:  eventItem.BlockHash,
				Number:     header.Number,
				TxHash:     eventItem.TransactionHash,
				Staker:     depositEvent.Staker,
				MantaToken: depositEvent.Weth, // 假设 MantaToken 是 Weth
				Strategy:   depositEvent.Strategy,
				Shares:     depositEvent.Shares,
				IsHandle:   0,
				Timestamp:  eventItem.Timestamp,
			}
			deposits = append(deposits, deposit)
		}
	}

	retryStrategy := &retry.ExponentialStrategy{Min: 1000, Max: 20_000, MaxJitter: 250}
	if _, err := retry.Do[interface{}](sm.smCtx, 10, retryStrategy, func() (interface{}, error) {
		if err := sm.db.Transaction(func(tx *database.DB) error {
			if err := tx.StrategyDeposit.StoreStrategyDeposit(deposits); err != nil {
				return err
			}
			return nil
		}); err != nil {
			log.Info("unable to persist batch", err)
			return nil, fmt.Errorf("unable to persist batch: %w", err)
		}
		return nil, nil
	}); err != nil {
		return err
	}

	log.Info("store strategy deposits success", "count", len(deposits))
	return nil
}
