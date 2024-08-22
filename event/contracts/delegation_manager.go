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

		if eventItem.EventSignature.String() == dm.DmAbi.Events["OperatorSharesDecreased"].ID.String() {

		}

		if eventItem.EventSignature.String() == dm.DmAbi.Events["WithdrawalQueued"].ID.String() {

		}

		if eventItem.EventSignature.String() == dm.DmAbi.Events["WithdrawalMigrated"].ID.String() {

		}

		if eventItem.EventSignature.String() == dm.DmAbi.Events["WithdrawalCompleted"].ID.String() {

		}
	}
	return nil
}
