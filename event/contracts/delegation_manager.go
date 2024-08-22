// package contracts

// import (
// 	"context"
// 	"fmt"
// 	"math/big"

// 	"github.com/ethereum/go-ethereum/accounts/abi"
// 	common2 "github.com/ethereum/go-ethereum/common"
// 	"github.com/ethereum/go-ethereum/log"
// 	"github.com/google/uuid"

// 	"github.com/eniac-x-labs/manta-indexer/bindings/dm"
// 	"github.com/eniac-x-labs/manta-indexer/config"
// 	"github.com/eniac-x-labs/manta-indexer/database"
// 	"github.com/eniac-x-labs/manta-indexer/database/event"
// 	"github.com/eniac-x-labs/manta-indexer/synchronizer/retry"
// )

// type DelegationManager struct {
// 	db       *database.DB
// 	DmAbi    *abi.ABI
// 	DmFilter *dm.DelegationManagerFilterer
// 	dmCtx    context.Context
// }

// func NewDelegationManager(db *database.DB) (*DelegationManager, error) {
// 	delegationAbi, err := dm.DelegationManagerMetaData.GetAbi()
// 	if err != nil {
// 		log.Error("get delegate manager abi fail", "err", err)
// 		return nil, err
// 	}

// 	DelegationManagerUnpack, err := dm.NewDelegationManagerFilterer(common2.Address{}, nil)
// 	if err != nil {
// 		log.Error("new delegation manager fail", "err", err)
// 		return nil, err
// 	}

// 	return &DelegationManager{
// 		db:       db,
// 		DmAbi:    delegationAbi,
// 		DmFilter: DelegationManagerUnpack,
// 		dmCtx:    context.Background(),
// 	}, nil
// }

// func (dm *DelegationManager) ProcessDelegationEvent(fromHeight *big.Int, toHeight *big.Int) error {
// 	contractEventFilter := event.ContractEvent{ContractAddress: common2.HexToAddress(config.DelegationManagerAddress)}
// 	contractEventList, err := dm.db.ContractEvent.ContractEventsWithFilter(contractEventFilter, fromHeight, toHeight)
// 	if err != nil {
// 		log.Error("get contracts event list fail", "err", err)
// 		return err
// 	}

// 	var operatorNodeUrlUpdates []event.OperatorNodeUrlUpdate
// 	var operatorRegisters []event.OperatorRegistered
// 	var stakerDelegates []event.StakerDelegate
// 	var operatorSharesIncreases []event.OperatorSharesIncrease
// 	var operatorDetailsModifies []event.OperatorDetailsModify
// 	var operatorSharesDecreases []event.OperatorSharesDecrease
// 	var withdrawalQueues []event.WithdrawalQueue
// 	var withdrawalMigrates []event.WithdrawalMigrate
// 	var minWithdrawalDelayBlocksSets []event.MinWithdrawalDelayBlocksSet
// 	var strategyWithdrawalDelayBlocksSets []event.StrategyWithdrawalDelayBlocksSet

// 	for _, eventItem := range contractEventList {
// 		rlpLog := eventItem.RLPLog

// 		header, err := dm.db.Blocks.BlockHeader(eventItem.BlockHash)
// 		if err != nil {
// 			return err
// 		}

// 		// OperatorNodeUrlUpdated
// 		if eventItem.EventSignature.String() == dm.DmAbi.Events["OperatorNodeUrlUpdated"].ID.String() {
// 			event, err := dm.DmFilter.ParseOperatorNodeUrlUpdated(*rlpLog)
// 			if err != nil {
// 				log.Error("parse operator node updated url fail", "err", err)
// 				return err
// 			}
// 			log.Info("parse operator node updated url success", "operator", event.Operator.String(), "metadataURI", event.MetadataURI)
// 			operatorNodeUrlUpdates = append(operatorNodeUrlUpdates, event.OperatorNodeUrlUpdate{
// 				GUID:        uuid.New(),
// 				BlockHash:   eventItem.BlockHash,
// 				Number:      header.Number,
// 				TxHash:      eventItem.TransactionHash,
// 				Operator:    event.Operator,
// 				MetadataURI: event.MetadataURI,
// 				IsHandle:    0,
// 				Timestamp:   eventItem.Timestamp,
// 			})
// 		}

// 		// OperatorRegistered
// 		if eventItem.EventSignature.String() == dm.DmAbi.Events["OperatorRegistered"].ID.String() {
// 			event, err := dm.DmFilter.ParseOperatorRegistered(*rlpLog)
// 			if err != nil {
// 				log.Error("parse operator register fail", "err", err)
// 				return err
// 			}
// 			log.Info("parse operator register success", "operator", event.Operator.String(), "earningsReceiver", event.OperatorDetails.EarningsReceiver.String())
// 			operatorRegisters = append(operatorRegisters, event.OperatorRegistered{
// 				GUID:                     uuid.New(),
// 				BlockHash:                eventItem.BlockHash,
// 				Number:                   header.Number,
// 				TxHash:                   eventItem.TransactionHash,
// 				Operator:                 event.Operator,
// 				EarningsReceiver:         event.OperatorDetails.EarningsReceiver,
// 				DelegationApprover:       event.OperatorDetails.DelegationApprover,
// 				StakerOptOutWindowBlocks: event.OperatorDetails.StakerOptOutWindowBlocks,
// 				MetadataURI:              event.OperatorDetails.MetadataURI,
// 				IsHandle:                 0,
// 				Timestamp:                eventItem.Timestamp,
// 			})
// 		}

// 		// StakerDelegated
// 		if eventItem.EventSignature.String() == dm.DmAbi.Events["StakerDelegated"].ID.String() {
// 			event, err := dm.DmFilter.ParseStakerDelegated(*rlpLog)
// 			if err != nil {
// 				log.Error("parse staker delegate event fail", "err", err)
// 				return err
// 			}
// 			log.Info("parse staker delegate success", "operator", event.Operator.String(), "staker", event.Staker.String())
// 			stakerDelegates = append(stakerDelegates, event.StakerDelegate{
// 				GUID:      uuid.New(),
// 				BlockHash: eventItem.BlockHash,
// 				Number:    header.Number,
// 				TxHash:    eventItem.TransactionHash,
// 				Operator:  event.Operator,
// 				Staker:    event.Staker,
// 				IsHandle:  0,
// 				Timestamp: eventItem.Timestamp,
// 			})
// 		}

// 		// OperatorSharesIncreased
// 		if eventItem.EventSignature.String() == dm.DmAbi.Events["OperatorSharesIncreased"].ID.String() {
// 			event, err := dm.DmFilter.ParseOperatorSharesIncreased(*rlpLog)
// 			if err != nil {
// 				log.Error("parse operator shares increased fail", "err", err)
// 				return err
// 			}
// 			log.Info("parse operator shares increased", "operator", event.Operator.String(), "staker", event.Staker.String())
// 			operatorSharesIncreases = append(operatorSharesIncreases, event.OperatorSharesIncrease{
// 				GUID:      uuid.New(),
// 				BlockHash: eventItem.BlockHash,
// 				Number:    header.Number,
// 				TxHash:    eventItem.TransactionHash,
// 				Operator:  event.Operator,
// 				Staker:    event.Staker,
// 				Shares:    event.Shares,
// 				IsHandle:  0,
// 				Timestamp: eventItem.Timestamp,
// 			})
// 		}

// 		// OperatorDetailsModified
// 		if eventItem.EventSignature.String() == dm.DmAbi.Events["OperatorDetailsModified"].ID.String() {
// 			event, err := dm.DmFilter.ParseOperatorDetailsModified(*rlpLog)
// 			if err != nil {
// 				log.Error("parse operator modified event fail", "err", err)
// 				return err
// 			}
// 			log.Info("parse operator modified success", "operator", event.Operator.String())
// 			operatorDetailsModifies = append(operatorDetailsModifies, event.OperatorDetailsModify{
// 				GUID:                     uuid.New(),
// 				BlockHash:                eventItem.BlockHash,
// 				Number:                   header.Number,
// 				TxHash:                   eventItem.TransactionHash,
// 				Operator:                 event.Operator,
// 				EarningsReceiver:         event.NewOperatorDetails.EarningsReceiver,
// 				DelegationApprover:       event.NewOperatorDetails.DelegationApprover,
// 				StakerOptOutWindowBlocks: event.NewOperatorDetails.StakerOptOutWindowBlocks,
// 				MetadataURI:              event.NewOperatorDetails.MetadataURI,
// 				IsHandle:                 0,
// 				Timestamp:                eventItem.Timestamp,
// 			})
// 		}

// 		// OperatorSharesDecreased
// 		if eventItem.EventSignature.String() == dm.DmAbi.Events["OperatorSharesDecreased"].ID.String() {
// 			event, err := dm.DmFilter.ParseOperatorSharesDecreased(*rlpLog)
// 			if err != nil {
// 				log.Error("parse operator shares decreased event fail", "err", err)
// 				return err
// 			}
// 			log.Info("parse operator shares decreased success", "operator", event.Operator.String(), "staker", event.Staker.String())
// 			operatorSharesDecreases = append(operatorSharesDecreases, event.OperatorSharesDecrease{
// 				GUID:      uuid.New(),
// 				BlockHash: eventItem.BlockHash,
// 				Number:    header.Number,
// 				TxHash:    eventItem.TransactionHash,
// 				Operator:  event.Operator,
// 				Staker:    event.Staker,
// 				Shares:    event.Shares,
// 				IsHandle:  0,
// 				Timestamp: eventItem.Timestamp,
// 			})
// 		}

// 		// WithdrawalQueued
// 		if eventItem.EventSignature.String() == dm.DmAbi.Events["WithdrawalQueued"].ID.String() {
// 			event, err := dm.DmFilter.ParseWithdrawalQueued(*rlpLog)
// 			if err != nil {
// 				log.Error("parse withdrawal queued event fail", "err", err)
// 				return err
// 			}
// 			log.Info("parse withdrawal queued success", "withdrawalRoot", common2.BytesToHash(event.WithdrawalRoot[:]).String())
// 			withdrawalQueues = append(withdrawalQueues, event.WithdrawalQueue{
// 				GUID:           uuid.New(),
// 				BlockHash:      eventItem.BlockHash,
// 				Number:         header.Number,
// 				TxHash:         eventItem.TransactionHash,
// 				WithdrawalRoot: common2.BytesToHash(event.WithdrawalRoot[:]),
// 				IsHandle:       0,
// 				Timestamp:      eventItem.Timestamp,
// 			})
// 		}

// 		// WithdrawalMigrated
// 		if eventItem.EventSignature.String() == dm.DmAbi.Events["WithdrawalMigrated"].ID.String() {
// 			event, err := dm.DmFilter.ParseWithdrawalMigrated(*rlpLog)
// 			if err != nil {
// 				log.Error("parse withdrawal migrated event fail", "err", err)
// 				return err
// 			}
// 			log.Info("parse withdrawal migrated success",
// 				"oldWithdrawalRoot", common2.BytesToHash(event.OldWithdrawalRoot[:]).String(),
// 				"newWithdrawalRoot", common2.BytesToHash(event.NewWithdrawalRoot[:]).String())
// 			withdrawalMigrates = append(withdrawalMigrates, event.WithdrawalMigrate{
// 				GUID:              uuid.New(),
// 				BlockHash:         eventItem.BlockHash,
// 				Number:            header.Number,
// 				TxHash:            eventItem.TransactionHash,
// 				OldWithdrawalRoot: common2.BytesToHash(event.OldWithdrawalRoot[:]),
// 				NewWithdrawalRoot: common2.BytesToHash(event.NewWithdrawalRoot[:]),
// 				IsHandle:          0,
// 				Timestamp:         eventItem.Timestamp,
// 			})
// 		}

// 		// MinWithdrawalDelayBlocksSet
// 		if eventItem.EventSignature.String() == dm.DmAbi.Events["MinWithdrawalDelayBlocksSet"].ID.String() {
// 			event, err := dm.DmFilter.ParseMinWithdrawalDelayBlocksSet(*rlpLog)
// 			if err != nil {
// 				log.Error("parse min withdrawal delay blocks set event fail", "err", err)
// 				return err
// 			}
// 			log.Info("parse min withdrawal delay blocks set success", "previousValue", event.PreviousValue, "newValue", event.NewValue)
// 			minWithdrawalDelayBlocksSets = append(minWithdrawalDelayBlocksSets, event.MinWithdrawalDelayBlocksSet{
// 				GUID:          uuid.New(),
// 				BlockHash:     eventItem.BlockHash,
// 				Number:        header.Number,
// 				TxHash:        eventItem.TransactionHash,
// 				PreviousValue: event.PreviousValue,
// 				NewValue:      event.NewValue,
// 				IsHandle:      0,
// 				Timestamp:     eventItem.Timestamp,
// 			})
// 		}

// 		// StrategyWithdrawalDelayBlocksSet
// 		if eventItem.EventSignature.String() == dm.DmAbi.Events["StrategyWithdrawalDelayBlocksSet"].ID.String() {
// 			event, err := dm.DmFilter.ParseStrategyWithdrawalDelayBlocksSet(*rlpLog)
// 			if err != nil {
// 				log.Error("parse strategy withdrawal delay blocks set event fail", "err", err)
// 				return err
// 			}
// 			log.Info("parse strategy withdrawal delay blocks set success",
// 				"strategy", event.Strategy.String(),
// 				"previousValue", event.PreviousValue,
// 				"newValue", event.NewValue)
// 			strategyWithdrawalDelayBlocksSets = append(strategyWithdrawalDelayBlocksSets, event.StrategyWithdrawalDelayBlocksSet{
// 				GUID:          uuid.New(),
// 				BlockHash:     eventItem.BlockHash,
// 				Number:        header.Number,
// 				TxHash:        eventItem.TransactionHash,
// 				Strategy:      event.Strategy,
// 				PreviousValue: event.PreviousValue,
// 				NewValue:      event.NewValue,
// 				IsHandle:      0,
// 				Timestamp:     eventItem.Timestamp,
// 			})
// 		}
// 	}

// 	retryStrategy := &retry.ExponentialStrategy{Min: 1000, Max: 20_000, MaxJitter: 250}
// 	if _, err := retry.Do[interface{}](dm.dmCtx, 10, retryStrategy, func() (interface{}, error) {
// 		if err := dm.db.Transaction(func(tx *database.DB) error {
// 			if len(operatorNodeUrlUpdates) > 0 {
// 				if err := tx.OperatorNodeUrlUpdate.StoreOperatorNodeUrlUpdate(operatorNodeUrlUpdates); err != nil {
// 					return err
// 				}
// 			}
// 			if len(operatorRegisters) > 0 {
// 				if err := tx.OperatorRegistered.StoreOperatorRegister(operatorRegisters); err != nil {
// 					return err
// 				}
// 			}
// 			if len(stakerDelegates) > 0 {
// 				if err := tx.StakerDelegate.StoreStakerDelegate(stakerDelegates); err != nil {
// 					return err
// 				}
// 			}
// 			if len(operatorSharesIncreases) > 0 {
// 				if err := tx.OperatorSharesIncrease.StoreOperatorSharesIncrease(operatorSharesIncreases); err != nil {
// 					return err
// 				}
// 			}
// 			if len(operatorDetailsModifies) > 0 {
// 				if err := tx.OperatorDetailsModify.StoreOperatorDetailsModify(operatorDetailsModifies); err != nil {
// 					return err
// 				}
// 			}
// 			if len(operatorSharesDecreases) > 0 {
// 				if err := tx.OperatorSharesDecrease.StoreOperatorSharesDecrease(operatorSharesDecreases); err != nil {
// 					return err
// 				}
// 			}
// 			if len(withdrawalQueues) > 0 {
// 				if err := tx.WithdrawalQueue.StoreWithdrawalQueue(withdrawalQueues); err != nil {
// 					return err
// 				}
// 			}
// 			if len(withdrawalMigrates) > 0 {
// 				if err := tx.WithdrawalMigrate.StoreWithdrawalMigrate(withdrawalMigrates); err != nil {
// 					return err
// 				}
// 			}
// 			if len(minWithdrawalDelayBlocksSets) > 0 {
// 				if err := tx.MinWithdrawalDelayBlocksSet.StoreMinWithdrawalDelayBlocksSet(minWithdrawalDelayBlocksSets); err != nil {
// 					return err
// 				}
// 			}
// 			if len(strategyWithdrawalDelayBlocksSets) > 0 {
// 				if err := tx.StrategyWithdrawalDelayBlocksSet.StoreStrategyWithdrawalDelayBlocksSet(strategyWithdrawalDelayBlocksSets); err != nil {
// 					return err
// 				}
// 			}
// 			return nil
// 		}); err != nil {
// 			log.Info("unable to persist batch", "err", err)
// 			return nil, fmt.Errorf("unable to persist batch: %w", err)
// 		}
// 		return nil, nil
// 	}); err != nil {
// 		return err
// 	}

// 	// Log success messages
// 	log.Info("store delegation events success",
// 		"operatorNodeUrlUpdates", len(operatorNodeUrlUpdates),
// 		"operatorRegisters", len(operatorRegisters),
// 		"stakerDelegates", len(stakerDelegates),
// 		"operatorSharesIncreases", len(operatorSharesIncreases),
// 		"operatorDetailsModifies", len(operatorDetailsModifies),
// 		"operatorSharesDecreases", len(operatorSharesDecreases),
// 		"withdrawalQueues", len(withdrawalQueues),
// 		"withdrawalMigrates", len(withdrawalMigrates),
// 		"minWithdrawalDelayBlocksSets", len(minWithdrawalDelayBlocksSets),
// 		"strategyWithdrawalDelayBlocksSets", len(strategyWithdrawalDelayBlocksSets))

// 	return nil
// }
