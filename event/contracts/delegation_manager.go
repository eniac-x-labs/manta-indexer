package contracts

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	common2 "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/google/uuid"

	"github.com/eniac-x-labs/manta-indexer/bindings/dm"
	"github.com/eniac-x-labs/manta-indexer/config"
	"github.com/eniac-x-labs/manta-indexer/database"
	"github.com/eniac-x-labs/manta-indexer/database/event"
	"github.com/eniac-x-labs/manta-indexer/synchronizer/retry"
)

type DelegationManager struct {
	db       *database.DB
	DmAbi    *abi.ABI
	DmFilter *dm.DelegationManagerFilterer
	dmCtx    context.Context
}

func NewDelegationManager(db *database.DB) (*DelegationManager, error) {
	delegationAbi, err := dm.DelegationManagerMetaData.GetAbi()
	if err != nil {
		log.Error("get delegate manager abi fail", "err", err)
		return nil, err
	}

	DelegationManagerUnpack, err := dm.NewDelegationManagerFilterer(common2.Address{}, nil)
	if err != nil {
		log.Error("new delegation manager fail", "err", err)
		return nil, err
	}

	return &DelegationManager{
		db:       db,
		DmAbi:    delegationAbi,
		DmFilter: DelegationManagerUnpack,
		dmCtx:    context.Background(),
	}, nil
}

func (dm *DelegationManager) ProcessDelegationEvent(fromHeight *big.Int, toHeight *big.Int) error {
	contractEventFilter := event.ContractEvent{ContractAddress: common2.HexToAddress(config.DelegationManagerAddress)}
	contractEventList, err := dm.db.ContractEvent.ContractEventsWithFilter(contractEventFilter, fromHeight, toHeight)
	if err != nil {
		log.Error("get contracts event list fail", "err", err)
		return err
	}

	var operatorNodeUrlUpdates []event.OperatorNodeUrlUpdate
	var operatorRegisters []event.OperatorRegistered
	var stakerDelegates []event.StakerDelegated
	var operatorSharesIncreases []event.OperatorSharesIncreased
	var operatorDetailsModifies []event.OperatorModified
	var operatorSharesDecreases []event.OperatorSharesDecreased
	var withdrawalQueues []event.WithdrawalQueued
	var withdrawalMigrates []event.WithdrawalMigrated
	var minWithdrawalDelayBlocksSets []event.MinWithdrawalDelayBlocksSet
	var strategyWithdrawalDelayBlocksSets []event.StrategyWithdrawalDelayBlocksSet

	for _, eventItem := range contractEventList {
		rlpLog := eventItem.RLPLog

		header, err := dm.db.Blocks.BlockHeader(eventItem.BlockHash)
		if err != nil {
			log.Error("ProcessDelegationEvent db Blocks BlockHeader by BlockHash fail", "err", err)
			return err
		}

		// OperatorNodeUrlUpdated
		if eventItem.EventSignature.String() == dm.DmAbi.Events["OperatorNodeUrlUpdated"].ID.String() {
			operatorNodeUrlUpdatedEvent, err := dm.DmFilter.ParseOperatorNodeUrlUpdated(*rlpLog)
			if err != nil {
				log.Error("parse operator node updated url fail", "err", err)
				return err
			}
			log.Info("parse operator node updated url success",
				"operator", operatorNodeUrlUpdatedEvent.Operator.String(),
				"metadataURI", operatorNodeUrlUpdatedEvent.MetadataURI)

			temp := event.OperatorNodeUrlUpdate{
				GUID:        uuid.New(),
				BlockHash:   eventItem.BlockHash,
				Number:      header.Number,
				TxHash:      eventItem.TransactionHash,
				Operator:    operatorNodeUrlUpdatedEvent.Operator,
				MetadataUri: operatorNodeUrlUpdatedEvent.MetadataURI,
				IsHandle:    0,
				Timestamp:   eventItem.Timestamp,
			}
			operatorNodeUrlUpdates = append(operatorNodeUrlUpdates, temp)
		}

		// OperatorRegistered
		if eventItem.EventSignature.String() == dm.DmAbi.Events["OperatorRegistered"].ID.String() {
			operatorRegisteredEvent, err := dm.DmFilter.ParseOperatorRegistered(*rlpLog)
			if err != nil {
				log.Error("parse operator register fail", "err", err)
				return err
			}
			log.Info("parse operator register success",
				"operator", operatorRegisteredEvent.Operator.String(),
				"earningsReceiver", operatorRegisteredEvent.OperatorDetails.EarningsReceiver.String())

			tempStakerOptOutWindowBlocks := new(big.Int)
			tempStakerOptOutWindowBlocks.SetUint64(uint64(operatorRegisteredEvent.OperatorDetails.StakerOptOutWindowBlocks))

			temp := event.OperatorRegistered{
				GUID:                     uuid.New(),
				BlockHash:                eventItem.BlockHash,
				Number:                   header.Number,
				TxHash:                   eventItem.TransactionHash,
				Operator:                 operatorRegisteredEvent.Operator,
				EarningsReceiver:         operatorRegisteredEvent.OperatorDetails.EarningsReceiver,
				DelegationApprover:       operatorRegisteredEvent.OperatorDetails.DelegationApprover,
				StakerOptoutWindowBlocks: tempStakerOptOutWindowBlocks,
				IsHandle:                 0,
				Timestamp:                eventItem.Timestamp,
			}
			operatorRegisters = append(operatorRegisters, temp)
		}

		// StakerDelegated
		if eventItem.EventSignature.String() == dm.DmAbi.Events["StakerDelegated"].ID.String() {
			stakerDelegatedEvnet, err := dm.DmFilter.ParseStakerDelegated(*rlpLog)
			if err != nil {
				log.Error("parse staker delegate event fail", "err", err)
				return err
			}
			log.Info("parse staker delegate success",
				"operator", stakerDelegatedEvnet.Operator.String(),
				"staker", stakerDelegatedEvnet.Staker.String())

			temp := event.StakerDelegated{
				GUID:      uuid.New(),
				BlockHash: eventItem.BlockHash,
				Number:    header.Number,
				TxHash:    eventItem.TransactionHash,
				Operator:  stakerDelegatedEvnet.Operator,
				Staker:    stakerDelegatedEvnet.Staker,
				IsHandle:  0,
				Timestamp: eventItem.Timestamp,
			}
			stakerDelegates = append(stakerDelegates, temp)
		}

		// OperatorSharesIncreased
		if eventItem.EventSignature.String() == dm.DmAbi.Events["OperatorSharesIncreased"].ID.String() {
			operatorSharesIncreasedEvent, err := dm.DmFilter.ParseOperatorSharesIncreased(*rlpLog)
			if err != nil {
				log.Error("parse operator shares increased fail", "err", err)
				return err
			}
			log.Info("parse operator shares increased",
				"operator", operatorSharesIncreasedEvent.Operator.String(),
				"staker", operatorSharesIncreasedEvent.Staker.String())

			temp := event.OperatorSharesIncreased{
				GUID:      uuid.New(),
				BlockHash: eventItem.BlockHash,
				Number:    header.Number,
				TxHash:    eventItem.TransactionHash,
				Operator:  operatorSharesIncreasedEvent.Operator,
				Staker:    operatorSharesIncreasedEvent.Staker,
				Strategy:  operatorSharesIncreasedEvent.Strategy,
				Shares:    operatorSharesIncreasedEvent.Shares,
				IsHandle:  0,
				Timestamp: eventItem.Timestamp,
			}
			operatorSharesIncreases = append(operatorSharesIncreases, temp)
		}

		// OperatorDetailsModified
		if eventItem.EventSignature.String() == dm.DmAbi.Events["OperatorDetailsModified"].ID.String() {
			operatorDetailsModifiedEvent, err := dm.DmFilter.ParseOperatorDetailsModified(*rlpLog)
			if err != nil {
				log.Error("parse operator modified event fail", "err", err)
				return err
			}
			log.Info("parse operator modified success", "operator", operatorDetailsModifiedEvent.Operator.String())

			tempStakerOptOutWindowBlocks := new(big.Int)
			tempStakerOptOutWindowBlocks.SetUint64(uint64(operatorDetailsModifiedEvent.NewOperatorDetails.StakerOptOutWindowBlocks))

			temp := event.OperatorModified{
				GUID:                     uuid.New(),
				BlockHash:                eventItem.BlockHash,
				Number:                   header.Number,
				TxHash:                   eventItem.TransactionHash,
				Operator:                 operatorDetailsModifiedEvent.Operator,
				EarningsReceiver:         operatorDetailsModifiedEvent.NewOperatorDetails.EarningsReceiver,
				DelegationApprover:       operatorDetailsModifiedEvent.NewOperatorDetails.DelegationApprover,
				StakerOptoutWindowBlocks: tempStakerOptOutWindowBlocks,
				IsHandle:                 0,
				Timestamp:                eventItem.Timestamp,
			}

			operatorDetailsModifies = append(operatorDetailsModifies, temp)
		}

		// OperatorSharesDecreased
		if eventItem.EventSignature.String() == dm.DmAbi.Events["OperatorSharesDecreased"].ID.String() {
			operatorSharesDecreasedEvent, err := dm.DmFilter.ParseOperatorSharesDecreased(*rlpLog)
			if err != nil {
				log.Error("parse operator shares decreased event fail", "err", err)
				return err
			}
			log.Info("parse operator shares decreased success",
				"operator", operatorSharesDecreasedEvent.Operator.String(),
				"staker", operatorSharesDecreasedEvent.Staker.String())

			temp := event.OperatorSharesDecreased{
				GUID:      uuid.New(),
				BlockHash: eventItem.BlockHash,
				Number:    header.Number,
				TxHash:    eventItem.TransactionHash,
				Operator:  operatorSharesDecreasedEvent.Operator,
				Staker:    operatorSharesDecreasedEvent.Staker,
				Strategy:  operatorSharesDecreasedEvent.Strategy,
				Shares:    operatorSharesDecreasedEvent.Shares,
				IsHandle:  0,
				Timestamp: eventItem.Timestamp,
			}

			operatorSharesDecreases = append(operatorSharesDecreases, temp)
		}

		// WithdrawalQueued
		if eventItem.EventSignature.String() == dm.DmAbi.Events["WithdrawalQueued"].ID.String() {
			withdrawalQueuedEvent, err := dm.DmFilter.ParseWithdrawalQueued(*rlpLog)
			if err != nil {
				log.Error("parse withdrawal queued event fail", "err", err)
				return err
			}
			log.Info("parse withdrawal queued success", "withdrawalRoot", common2.BytesToHash(withdrawalQueuedEvent.WithdrawalRoot[:]).String())

			startBlockBigInt := new(big.Int)
			startBlockBigInt.SetUint64(uint64(withdrawalQueuedEvent.Withdrawal.StartBlock))

			temp := event.WithdrawalQueued{
				GUID:           uuid.New(),
				BlockHash:      eventItem.BlockHash,
				Number:         header.Number,
				TxHash:         eventItem.TransactionHash,
				WithdrawalRoot: common2.BytesToHash(withdrawalQueuedEvent.WithdrawalRoot[:]),

				Staker:      withdrawalQueuedEvent.Withdrawal.Staker,
				DelegatedTo: withdrawalQueuedEvent.Withdrawal.DelegatedTo,
				Withdrawer:  withdrawalQueuedEvent.Withdrawal.Withdrawer,
				Nonce:       withdrawalQueuedEvent.Withdrawal.Nonce,
				StartBlock:  startBlockBigInt,
				Strategies:  addressListToString(withdrawalQueuedEvent.Withdrawal.Strategies),
				Shares:      bigIntListToString(withdrawalQueuedEvent.Withdrawal.Shares),
				IsHandle:    0,
				Timestamp:   eventItem.Timestamp,
			}

			withdrawalQueues = append(withdrawalQueues, temp)
		}

		// WithdrawalMigrated
		if eventItem.EventSignature.String() == dm.DmAbi.Events["WithdrawalMigrated"].ID.String() {
			withdrawalMigratedEvent, err := dm.DmFilter.ParseWithdrawalMigrated(*rlpLog)
			if err != nil {
				log.Error("parse withdrawal migrated event fail", "err", err)
				return err
			}
			log.Info("parse withdrawal migrated success",
				"oldWithdrawalRoot", common2.BytesToHash(withdrawalMigratedEvent.OldWithdrawalRoot[:]).String(),
				"newWithdrawalRoot", common2.BytesToHash(withdrawalMigratedEvent.NewWithdrawalRoot[:]).String())

			temp := event.WithdrawalMigrated{
				GUID:              uuid.New(),
				BlockHash:         eventItem.BlockHash,
				Number:            header.Number,
				TxHash:            eventItem.TransactionHash,
				OldWithdrawalRoot: common2.BytesToHash(withdrawalMigratedEvent.OldWithdrawalRoot[:]),
				NewWithdrawalRoot: common2.BytesToHash(withdrawalMigratedEvent.NewWithdrawalRoot[:]),
				IsHandle:          0,
				Timestamp:         eventItem.Timestamp,
			}

			withdrawalMigrates = append(withdrawalMigrates, temp)
		}

		// MinWithdrawalDelayBlocksSet
		if eventItem.EventSignature.String() == dm.DmAbi.Events["MinWithdrawalDelayBlocksSet"].ID.String() {
			minWithdrawalDelayBlocksSetEvent, err := dm.DmFilter.ParseMinWithdrawalDelayBlocksSet(*rlpLog)
			if err != nil {
				log.Error("parse min withdrawal delay blocks set event fail", "err", err)
				return err
			}
			log.Info("parse min withdrawal delay blocks set success",
				"previousValue", minWithdrawalDelayBlocksSetEvent.PreviousValue,
				"newValue", minWithdrawalDelayBlocksSetEvent.NewValue)
			temp := event.MinWithdrawalDelayBlocksSet{
				GUID:          uuid.New(),
				BlockHash:     eventItem.BlockHash,
				Number:        header.Number,
				TxHash:        eventItem.TransactionHash,
				PreviousValue: minWithdrawalDelayBlocksSetEvent.PreviousValue,
				NewValue:      minWithdrawalDelayBlocksSetEvent.NewValue,
				IsHandle:      0,
				Timestamp:     eventItem.Timestamp,
			}

			minWithdrawalDelayBlocksSets = append(minWithdrawalDelayBlocksSets, temp)
		}

		// StrategyWithdrawalDelayBlocksSet
		if eventItem.EventSignature.String() == dm.DmAbi.Events["StrategyWithdrawalDelayBlocksSet"].ID.String() {
			strategyWithdrawalDelayBlocksSetEvent, err := dm.DmFilter.ParseStrategyWithdrawalDelayBlocksSet(*rlpLog)
			if err != nil {
				log.Error("parse strategy withdrawal delay blocks set event fail", "err", err)
				return err
			}
			log.Info("parse strategy withdrawal delay blocks set success",
				"strategy", strategyWithdrawalDelayBlocksSetEvent.Strategy.String(),
				"previousValue", strategyWithdrawalDelayBlocksSetEvent.PreviousValue,
				"newValue", strategyWithdrawalDelayBlocksSetEvent.NewValue)

			temp := event.StrategyWithdrawalDelayBlocksSet{
				GUID:          uuid.New(),
				BlockHash:     eventItem.BlockHash,
				Number:        header.Number,
				TxHash:        eventItem.TransactionHash,
				Strategy:      strategyWithdrawalDelayBlocksSetEvent.Strategy,
				PreviousValue: strategyWithdrawalDelayBlocksSetEvent.PreviousValue,
				NewValue:      strategyWithdrawalDelayBlocksSetEvent.NewValue,
				IsHandle:      0,
				Timestamp:     eventItem.Timestamp,
			}

			strategyWithdrawalDelayBlocksSets = append(strategyWithdrawalDelayBlocksSets, temp)
		}
	}

	retryStrategy := &retry.ExponentialStrategy{Min: 1000, Max: 20_000, MaxJitter: 250}
	if _, err := retry.Do[interface{}](dm.dmCtx, 10, retryStrategy, func() (interface{}, error) {
		if err := dm.db.Transaction(func(tx *database.DB) error {
			if len(operatorNodeUrlUpdates) > 0 {
				if err := tx.OperatorNodeUrlUpdate.StoreOperatorNodeUrlUpdate(operatorNodeUrlUpdates); err != nil {
					return err
				}
			}
			if len(operatorRegisters) > 0 {
				if err := tx.OperatorRegistered.StoreOperatorRegistered(operatorRegisters); err != nil {
					return err
				}
			}
			if len(stakerDelegates) > 0 {
				if err := tx.StakerDelegated.StoreStakerDelegated(stakerDelegates); err != nil {
					return err
				}
			}
			if len(operatorSharesIncreases) > 0 {
				if err := tx.OperatorSharesIncreased.StoreOperatorSharesIncreased(operatorSharesIncreases); err != nil {
					return err
				}
			}
			if len(operatorDetailsModifies) > 0 {
				if err := tx.OperatorModified.StoreOperatorModified(operatorDetailsModifies); err != nil {
					return err
				}
			}
			if len(operatorSharesDecreases) > 0 {
				if err := tx.OperatorSharesDecreased.StoreOperatorSharesDecreased(operatorSharesDecreases); err != nil {
					return err
				}
			}
			if len(withdrawalQueues) > 0 {
				if err := tx.WithdrawalQueued.StoreWithdrawalQueued(withdrawalQueues); err != nil {
					return err
				}
			}
			if len(withdrawalMigrates) > 0 {
				if err := tx.WithdrawalMigrated.StoreWithdrawalMigrated(withdrawalMigrates); err != nil {
					return err
				}
			}
			if len(minWithdrawalDelayBlocksSets) > 0 {
				if err := tx.MinWithdrawalDelayBlocksSet.StoreMinWithdrawalDelayBlocksSet(minWithdrawalDelayBlocksSets); err != nil {
					return err
				}
			}
			if len(strategyWithdrawalDelayBlocksSets) > 0 {
				if err := tx.StrategyWithdrawalDelayBlocksSet.StoreStrategyWithdrawalDelayBlocksSet(strategyWithdrawalDelayBlocksSets); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			log.Info("unable to persist batch", "err", err)
			return nil, fmt.Errorf("unable to persist batch: %w", err)
		}
		return nil, nil
	}); err != nil {
		return err
	}

	// Log success messages
	log.Info("store delegation events success",
		"operatorNodeUrlUpdates", len(operatorNodeUrlUpdates),
		"operatorRegisters", len(operatorRegisters),
		"stakerDelegates", len(stakerDelegates),
		"operatorSharesIncreases", len(operatorSharesIncreases),
		"operatorDetailsModifies", len(operatorDetailsModifies),
		"operatorSharesDecreases", len(operatorSharesDecreases),
		"withdrawalQueues", len(withdrawalQueues),
		"withdrawalMigrates", len(withdrawalMigrates),
		"minWithdrawalDelayBlocksSets", len(minWithdrawalDelayBlocksSets),
		"strategyWithdrawalDelayBlocksSets", len(strategyWithdrawalDelayBlocksSets))

	return nil
}

func addressListToString(addressList []common2.Address) string {
	if len(addressList) == 0 {
		return "addressList is empty"
	}
	var result []string
	for _, addr := range addressList {
		result = append(result, addr.Hex())
	}
	return strings.Join(result, ",")
}

func bigIntListToString(bigIntList []*big.Int) string {
	if len(bigIntList) == 0 {
		return "bigIntList is empty"
	}
	var result []string
	for _, bi := range bigIntList {
		result = append(result, bi.String())
	}
	return strings.Join(result, ",")
}
