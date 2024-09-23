package contracts

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	common2 "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/log"

	"github.com/eniac-x-labs/manta-indexer/bindings/msm"
	"github.com/eniac-x-labs/manta-indexer/config"
	"github.com/eniac-x-labs/manta-indexer/database"
	"github.com/eniac-x-labs/manta-indexer/database/event"
	"github.com/eniac-x-labs/manta-indexer/database/event/finality"
	"github.com/eniac-x-labs/manta-indexer/synchronizer/retry"
)

type MantaServiceManager struct {
	MsmAbi    *abi.ABI
	MsmFilter *msm.MantaServiceManagerFilterer
	MsmCtx    context.Context
}

func NewMantaServiceManager() (*MantaServiceManager, error) {
	mantaServiceManagerAbi, err := msm.MantaServiceManagerMetaData.GetAbi()
	if err != nil {
		log.Error("get manta service manager abi fail", "err", err)
		return nil, err
	}

	mantaServiceManagerUnpack, err := msm.NewMantaServiceManagerFilterer(common2.Address{}, nil)
	if err != nil {
		log.Error("new delegation manager fail", "err", err)
		return nil, err
	}

	return &MantaServiceManager{
		MsmAbi:    mantaServiceManagerAbi,
		MsmFilter: mantaServiceManagerUnpack,
		MsmCtx:    context.Background(),
	}, nil
}

func (msm *MantaServiceManager) ProcessMantaServiceManager(db *database.DB, fromHeight *big.Int, toHeight *big.Int) error {
	contractEventFilter := event.ContractEvent{ContractAddress: common2.HexToAddress(config.MantaServiceManagerAddress)}
	contractEventList, err := db.ContractEvent.ContractEventsWithFilter(contractEventFilter, fromHeight, toHeight)
	if err != nil {
		log.Error("get contracts event list fail", "err", err)
		return err
	}

	var finalityVerifiedList []finality.FinalityVerified

	for _, eventItem := range contractEventList {
		rlpLog := eventItem.RLPLog

		header, err := db.Blocks.BlockHeader(eventItem.BlockHash)
		if err != nil {
			log.Error("ProcessMantaServiceManager db Blocks BlockHeader by BlockHash fail", "err", err)
			return err
		}

		if eventItem.EventSignature.String() == msm.MsmAbi.Events["FinalityVerified"].ID.String() {
			finalityVerifiedEvent, err := msm.MsmFilter.ParseFinalityVerified(*rlpLog)
			if err != nil {
				log.Error("parse finality verified event fail", "err", err)
				return err
			}
			log.Info("parse finality verified success",
				"proposer", finalityVerifiedEvent.Proposer.String(),
				"l1BlockNumber", finalityVerifiedEvent.L1BlockNumber.String(),
				"l2BlockNumber", finalityVerifiedEvent.L2BlockNumber.String(),
				"outputRoot", hexutil.Encode(finalityVerifiedEvent.OutputRoot[:]))

			temp := finality.FinalityVerified{
				GUID:          uuid.New(),
				Proposer:      finalityVerifiedEvent.Proposer,
				TxBlockNumber: header.Number,
				L1BlockNumber: finalityVerifiedEvent.L1BlockNumber,
				L2BlockNumber: finalityVerifiedEvent.L2BlockNumber,
				L1BlockHash:   hexutil.Encode(finalityVerifiedEvent.L1BlockHash[:]),
				OutputRoot:    hexutil.Encode(finalityVerifiedEvent.OutputRoot[:]),
				Timestamp:     eventItem.Timestamp,
			}
			finalityVerifiedList = append(finalityVerifiedList, temp)
		}

	}

	retryStrategy := &retry.ExponentialStrategy{Min: 1000, Max: 20_000, MaxJitter: 250}
	if _, err := retry.Do[interface{}](msm.MsmCtx, 10, retryStrategy, func() (interface{}, error) {
		if err := db.Transaction(func(tx *database.DB) error {
			if err := tx.FinalityVerified.StoreFinalityVerified(finalityVerifiedList); err != nil {
				return err
			}

			// Log success messages
			log.Info("store finality verified events success",
				"finalityVerifiedList", len(finalityVerifiedList),
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

	return nil
}
