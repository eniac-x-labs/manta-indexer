package contracts

import (
	"fmt"
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
	db                *database.DB
	DmAbi             *abi.ABI
	ContractEventList []event.ContractEvent
	DmFilter          *dm.DelegationManagerFilterer
}

func NewDelegationManager(db *database.DB, fromHeight *big.Int, toHeight *big.Int) (*DelegationManager, error) {
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
	log.Info("====================================================================================")
	log.Info("new delegation manager")
	log.Info("====================================================================================")
	contractEventFilter := event.ContractEvent{ContractAddress: common2.HexToAddress(config.DelegationManagerAddress)}
	eventList, err := db.ContractEvent.ContractEventsWithFilter(contractEventFilter, fromHeight, toHeight)
	if err != nil {
		log.Error("get contracts event list fail", "err", err)
		return nil, err
	}
	return &DelegationManager{
		db:                db,
		DmAbi:             delegationAbi,
		ContractEventList: eventList,
		DmFilter:          DelegationManagerUnpack,
	}, nil
}

func (dm *DelegationManager) ProcessDelegationRegister() {
	for _, eventItem := range dm.ContractEventList {
		/*
			0xfebe5cd24b2cbc7b065b9d0fdeb904461e4afcff57dd57acda1e7832031ba7ac
			 0xc3ee9f2e5fda98e8066a1f745b2df9285f416fe98cf2559cd21484b3d8743304
			 0x8e8485583a2310d41f7c82b9427d0bd49bad74bb9cff9d3402a29d8f9b28a0e2
			 0x826d13513a58153c5878cd93af2008d3f8dfc32049e748b380b1b385645e280b
		*/
		log.Info("====================================================================================")
		log.Info("eventItem.EventSignature.String()", eventItem.EventSignature.String())
		log.Info("dm.DmAbi.Events[StakerDelegated].ID.String()", dm.DmAbi.Events["StakerDelegated"].ID.String())
		log.Info("dm.DmAbi.Events[OperatorRegistered].ID.String()", dm.DmAbi.Events["OperatorRegistered"].ID.String())
		log.Info("dm.DmAbi.Events[OperatorNodeUrlUpdated].ID.String()", dm.DmAbi.Events["OperatorNodeUrlUpdated"].ID.String())
		log.Info("dm.DmAbi.Events[OperatorSharesIncreased].ID.String()", dm.DmAbi.Events["OperatorSharesIncreased"].ID.String())
		log.Info("====================================================================================")
		// OperatorRegistered(msg.sender, registeringOperatorDetails);
		if eventItem.EventSignature.String() == dm.DmAbi.Events["OperatorNodeUrlUpdated"].ID.String() {
			rlpLog := eventItem.RLPLog
			unPackEvent, err := dm.DmFilter.ParseOperatorNodeUrlUpdated(*rlpLog)
			if err != nil {
				log.Error("parse operator register fail", "err", err)
				return
			}
			fmt.Println("unPackEvent.Operator.String()===", unPackEvent.Operator.String())
			fmt.Println("unPackEvent.OperatorDetails.DelegationApprover.String()===", unPackEvent.MetadataURI)
		}

		// OperatorNodeUrlUpdated(msg.sender, nodeUrl);
		if eventItem.EventSignature.String() == dm.DmAbi.Events["OperatorNodeUrlUpdated"].ID.String() {

		}
	}
}
