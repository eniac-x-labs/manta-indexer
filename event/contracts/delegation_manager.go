package contracts

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	common2 "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"

	"github.com/eniac-x-labs/manta-indexer/bindings/dm"
	"github.com/eniac-x-labs/manta-indexer/config"
	"github.com/eniac-x-labs/manta-indexer/database"
	"github.com/eniac-x-labs/manta-indexer/database/event"
)

type DelegationManager struct {
	db       *database.DB
	DmAbi    *abi.ABI
	DmFilter *dm.DelegationManagerFilterer
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
	}, nil
}

func (dm *DelegationManager) ProcessDelegationEvent(fromHeight *big.Int, toHeight *big.Int) error {
	contractEventFilter := event.ContractEvent{ContractAddress: common2.HexToAddress(config.DelegationManagerAddress)}
	contractEventList, err := dm.db.ContractEvent.ContractEventsWithFilter(contractEventFilter, fromHeight, toHeight)
	if err != nil {
		log.Error("get contracts event list fail", "err", err)
		return err
	}
	for _, eventItem := range contractEventList {
		rlpLog := eventItem.RLPLog
		// OperatorNodeUrlUpdated
		if eventItem.EventSignature.String() == dm.DmAbi.Events["OperatorNodeUrlUpdated"].ID.String() {
			nodeUrlUpdateEvent, err := dm.DmFilter.ParseOperatorNodeUrlUpdated(*rlpLog)
			if err != nil {
				log.Error("parse operator node updated url fail", "err", err)
				return err
			}
			log.Info("parse operator node updated url success", "operator", nodeUrlUpdateEvent.Operator.String(), "metadataURI", nodeUrlUpdateEvent.MetadataURI)

			// operatorNodeUrlUpdates := make([]event.OperatorNodeUrlUpdate, len(nodeUrlUpdateEvent))
			// for i := range nodeUrlUpdateEvent {
			// 	operatorNodeUrlUpdates[0] = event.OperatorNodeUrlUpdate{
			// 		GUID:        uuid.New(),
			// 		BlockHash:   eventItem.BlockHash,
			// 		Number:      new(big.Int).SetUint64(eventItem.BlockNumber),
			// 		TxHash:      eventItem.TxHash,
			// 		Operator:    nodeUrlUpdateEvent[i].Operator,
			// 		MetadataURI: nodeUrlUpdateEvent[i].MetadataURI,
			// 		IsHandle:    0,
			// 		Timestamp:   eventItem.Timestamp,
			// 	}
			// }

			// if len(operatorNodeUrlUpdates) > 0 {
			// 	if err := dm.db.OperatorNodeUrlUpdate.StoreOperatorNodeUrlUpdate(operatorNodeUrlUpdates); err != nil {
			// 		log.Error("store operator node url update fail", "err", err)
			// 		return err
			// 	}
			// }
		}

		// OperatorRegistered
		if eventItem.EventSignature.String() == dm.DmAbi.Events["OperatorRegistered"].ID.String() {
			nodeRegisterEvent, err := dm.DmFilter.ParseOperatorRegistered(*rlpLog)
			if err != nil {
				log.Error("parse operator register fail", "err", err)
				return err
			}
			log.Info("parse operator register success", "operator", nodeRegisterEvent.Operator.String(), "earningsReceiver", nodeRegisterEvent.OperatorDetails.EarningsReceiver)
		}

		// StakerDelegated
		if eventItem.EventSignature.String() == dm.DmAbi.Events["StakerDelegated"].ID.String() {
			stakerDelegatedEvent, err := dm.DmFilter.ParseStakerDelegated(*rlpLog)
			if err != nil {
				log.Error("parse staker delegate event fail", "err", err)
				return err
			}
			log.Info("parse staker delegate success", "operator", stakerDelegatedEvent.Operator.String(), "staker", stakerDelegatedEvent.Staker.String())

		}

		// OperatorSharesIncreased
		if eventItem.EventSignature.String() == dm.DmAbi.Events["OperatorSharesIncreased"].ID.String() {
			operatorSharesIncreasedEvent, err := dm.DmFilter.ParseOperatorSharesIncreased(*rlpLog)
			if err != nil {
				log.Error("parse operator shares increased fail", "err", err)
				return err
			}
			log.Info("parse operator shares increased", "operator", operatorSharesIncreasedEvent.Operator.String(), "staker", operatorSharesIncreasedEvent.Staker.String())
		}

		// StakerDelegated
		if eventItem.EventSignature.String() == dm.DmAbi.Events["StakerDelegated"].ID.String() {
			stakerDelegatedEvent, err := dm.DmFilter.ParseStakerDelegated(*rlpLog)
			if err != nil {
				log.Error("parse staker delegate event fail", "err", err)
				return err
			}
			log.Info("parse staker delegate success", "operator", stakerDelegatedEvent.Operator.String(), "staker", stakerDelegatedEvent.Staker.String())
		}

		// OperatorSharesIncreased
		if eventItem.EventSignature.String() == dm.DmAbi.Events["OperatorSharesIncreased"].ID.String() {
			operatorSharesIncreasedEvent, err := dm.DmFilter.ParseOperatorSharesIncreased(*rlpLog)
			if err != nil {
				log.Error("parse operator shares increased fail", "err", err)
				return err
			}
			log.Info("parse operator shares increased", "operator", operatorSharesIncreasedEvent.Operator.String(), "staker", operatorSharesIncreasedEvent.Staker.String())
		}

		// OperatorRegistered
		if eventItem.EventSignature.String() == dm.DmAbi.Events["OperatorRegistered"].ID.String() {
			operatorRegisteredEvent, err := dm.DmFilter.ParseOperatorRegistered(*rlpLog)
			if err != nil {
				log.Error("parse operator registered event fail", "err", err)
				return err
			}
			log.Info("parse operator registered success", "operator", operatorRegisteredEvent.Operator.String())
		}

		// OperatorDetailsModified
		if eventItem.EventSignature.String() == dm.DmAbi.Events["OperatorDetailsModified"].ID.String() {
			operatorModifiedEvent, err := dm.DmFilter.ParseOperatorDetailsModified(*rlpLog)
			if err != nil {
				log.Error("parse operator modified event fail", "err", err)
				return err
			}
			log.Info("parse operator modified success", "operator", operatorModifiedEvent.Operator.String())
		}

		// OperatorNodeUrlUpdated
		if eventItem.EventSignature.String() == dm.DmAbi.Events["OperatorNodeUrlUpdated"].ID.String() {
			operatorNodeUrlUpdateEvent, err := dm.DmFilter.ParseOperatorNodeUrlUpdated(*rlpLog)
			if err != nil {
				log.Error("parse operator node url update event fail", "err", err)
				return err
			}
			log.Info("parse operator node url update success", "operator", operatorNodeUrlUpdateEvent.Operator.String())
		}

		// OperatorSharesDecreased
		if eventItem.EventSignature.String() == dm.DmAbi.Events["OperatorSharesDecreased"].ID.String() {
			operatorSharesDecreasedEvent, err := dm.DmFilter.ParseOperatorSharesDecreased(*rlpLog)
			if err != nil {
				log.Error("parse operator shares decreased event fail", "err", err)
				return err
			}
			log.Info("parse operator shares decreased success", "operator", operatorSharesDecreasedEvent.Operator.String(), "staker", operatorSharesDecreasedEvent.Staker.String())
		}

		// WithdrawalQueued
		if eventItem.EventSignature.String() == dm.DmAbi.Events["WithdrawalQueued"].ID.String() {
			withdrawalQueuedEvent, err := dm.DmFilter.ParseWithdrawalQueued(*rlpLog)
			if err != nil {
				log.Error("parse withdrawal queued event fail", "err", err)
				return err
			}
			log.Info("parse withdrawal queued success", "withdrawalRoot", common2.BytesToHash(withdrawalQueuedEvent.WithdrawalRoot[:]).String())
		}

		// WithdrawalMigrated
		if eventItem.EventSignature.String() == dm.DmAbi.Events["WithdrawalMigrated"].ID.String() {
			withdrawalMigratedEvent, err := dm.DmFilter.ParseWithdrawalMigrated(*rlpLog)
			if err != nil {
				log.Error("parse withdrawal migrated event fail", "err", err)
				return err
			}
			log.Info("parse withdrawal migrated success", "oldWithdrawalRoot", common2.BytesToHash(withdrawalMigratedEvent.OldWithdrawalRoot[:]).String())
			log.Info("parse withdrawal migrated success", "newWithdrawalRoot", common2.BytesToHash(withdrawalMigratedEvent.NewWithdrawalRoot[:]).String())
		}

		// MinWithdrawalDelayBlocksSet
		if eventItem.EventSignature.String() == dm.DmAbi.Events["MinWithdrawalDelayBlocksSet"].ID.String() {
			minWithdrawalDelayBlocksSetEvent, err := dm.DmFilter.ParseMinWithdrawalDelayBlocksSet(*rlpLog)
			if err != nil {
				log.Error("parse min withdrawal delay blocks set event fail", "err", err)
				return err
			}
			log.Info("parse min withdrawal delay blocks set success", "previousValue", minWithdrawalDelayBlocksSetEvent.PreviousValue, "newValue", minWithdrawalDelayBlocksSetEvent.NewValue)
		}

		// StrategyWithdrawalDelayBlocksSet
		if eventItem.EventSignature.String() == dm.DmAbi.Events["StrategyWithdrawalDelayBlocksSet"].ID.String() {
			strategyWithdrawalDelayBlocksSetEvent, err := dm.DmFilter.ParseStrategyWithdrawalDelayBlocksSet(*rlpLog)
			if err != nil {
				log.Error("parse strategy withdrawal delay blocks set event fail", "err", err)
				return err
			}
			log.Info("parse strategy withdrawal delay blocks set success", "strategy", strategyWithdrawalDelayBlocksSetEvent.Strategy.String(), "previousValue", strategyWithdrawalDelayBlocksSetEvent.PreviousValue, "newValue", strategyWithdrawalDelayBlocksSetEvent.NewValue)
		}
	}
	return nil
}
