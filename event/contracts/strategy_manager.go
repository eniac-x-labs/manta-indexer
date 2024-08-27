package contracts

import (
	"context"
	"fmt"

	"math/big"

	"github.com/google/uuid"

	"github.com/ethereum/go-ethereum/accounts/abi"
	common2 "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"

	"github.com/eniac-x-labs/manta-indexer/bindings/sm"
	"github.com/eniac-x-labs/manta-indexer/config"
	"github.com/eniac-x-labs/manta-indexer/database"
	"github.com/eniac-x-labs/manta-indexer/database/event"
	"github.com/eniac-x-labs/manta-indexer/database/event/staker"
	"github.com/eniac-x-labs/manta-indexer/database/event/strategies"
	"github.com/eniac-x-labs/manta-indexer/synchronizer/retry"
)

var MantaTokenAddress = "0xFEE297254eC9B60d06f6e5af4E154962f9dCcE88"

type StrategyManager struct {
	SmAbi    *abi.ABI
	SmFilter *sm.StrategyManagerFilterer
	smCtx    context.Context
}

func NewStrategyManager() (*StrategyManager, error) {
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
		SmAbi:    strategyAbi,
		SmFilter: strategyUnpack,
		smCtx:    context.Background(),
	}, nil
}

func (sm *StrategyManager) ProcessStrategyManager(db *database.DB, fromHeight *big.Int, toHeight *big.Int) error {
	contractEventFilter := event.ContractEvent{ContractAddress: common2.HexToAddress(config.StrategyManagerAddress)}
	contractEventList, err := db.ContractEvent.ContractEventsWithFilter(contractEventFilter, fromHeight, toHeight)
	if err != nil {
		log.Error("get contracts event list fail", "err", err)
		return err
	}

	var deposits []staker.StrategyDeposit
	var strategiesAdd []strategies.Strategies
	var strategiesRemove []strategies.Strategies
	for _, eventItem := range contractEventList {
		rlpLog := eventItem.RLPLog

		header, err := db.Blocks.BlockHeader(eventItem.BlockHash)
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

			deposit := staker.StrategyDeposit{
				GUID:       uuid.New(),
				BlockHash:  eventItem.BlockHash,
				Number:     header.Number,
				TxHash:     eventItem.TransactionHash,
				Staker:     depositEvent.Staker,
				MantaToken: depositEvent.MantaToken,
				Strategy:   depositEvent.Strategy,
				Shares:     depositEvent.Shares,
				IsHandle:   0,
				Timestamp:  eventItem.Timestamp,
			}
			deposits = append(deposits, deposit)
		}

		if eventItem.EventSignature.String() == sm.SmAbi.Events["StrategyAddedToDepositWhitelist"].ID.String() {
			strategyAdded, err := sm.SmFilter.ParseStrategyAddedToDepositWhitelist(*rlpLog)
			if err != nil {
				log.Error("parse deposit event fail", "err", err)
				return err
			}
			log.Info("parse strategy added to Deposit whitelist success", "Strategy", strategyAdded.Strategy.String())
			strategy := strategies.Strategies{
				GUID:       uuid.New(),
				BlockHash:  eventItem.BlockHash,
				Number:     header.Number,
				TxHash:     eventItem.TransactionHash,
				Strategy:   strategyAdded.Strategy,
				Tvl:        big.NewInt(0),
				MantaToken: common2.HexToAddress(MantaTokenAddress),
				IsHandle:   0,
				Timestamp:  eventItem.Timestamp,
			}
			strategiesAdd = append(strategiesAdd, strategy)
		}

		if eventItem.EventSignature.String() == sm.SmAbi.Events["StrategyRemovedFromDepositWhitelist"].ID.String() {
			strategyRemoved, err := sm.SmFilter.ParseStrategyRemovedFromDepositWhitelist(*rlpLog)
			if err != nil {
				log.Error("parse strategy removed from deposit whitelist fail", "err", err)
				return err
			}
			log.Info("parse strategy removed from deposit whitelist success", "Strategy", strategyRemoved.Strategy.String())
			strategy := strategies.Strategies{
				GUID:       uuid.New(),
				BlockHash:  eventItem.BlockHash,
				Number:     header.Number,
				TxHash:     eventItem.TransactionHash,
				Strategy:   strategyRemoved.Strategy,
				Tvl:        big.NewInt(0),
				MantaToken: common2.HexToAddress(MantaTokenAddress),
				IsHandle:   0,
				Timestamp:  eventItem.Timestamp,
			}
			strategiesRemove = append(strategiesRemove, strategy)
		}
	}

	retryStrategy := &retry.ExponentialStrategy{Min: 1000, Max: 20_000, MaxJitter: 250}
	if _, err := retry.Do[interface{}](sm.smCtx, 10, retryStrategy, func() (interface{}, error) {
		if err := db.Transaction(func(tx *database.DB) error {
			if len(deposits) > 0 {
				if err := tx.StrategyDeposit.StoreStrategyDeposit(deposits); err != nil {
					return err
				}
			}
			if len(strategiesAdd) > 0 {
				if err := tx.Strategies.StoreStrategies(strategiesAdd); err != nil {
					return err
				}
			}
			if len(strategiesRemove) > 0 {
				if err := tx.Strategies.RemoveStoreStrategies(strategiesRemove); err != nil {
					return err
				}
			}
			// Log success messages
			log.Info("store strategy manager events success",
				"deposits", len(deposits),
				"strategiesAdd", len(strategiesAdd),
				"strategiesRemove", len(strategiesRemove),
			)
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
