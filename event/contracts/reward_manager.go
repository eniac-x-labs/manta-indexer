package contracts

import (
	"context"
	"fmt"
	"math/big"

	"github.com/google/uuid"

	"github.com/ethereum/go-ethereum/accounts/abi"
	common2 "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"

	"github.com/eniac-x-labs/manta-indexer/bindings/rm"
	"github.com/eniac-x-labs/manta-indexer/config"
	"github.com/eniac-x-labs/manta-indexer/database"
	"github.com/eniac-x-labs/manta-indexer/database/event"
	"github.com/eniac-x-labs/manta-indexer/synchronizer/retry"
)

type RewardManager struct {
	db       *database.DB
	RmAbi    *abi.ABI
	RmFilter *rm.RewardManagerFilterer
	rmCtx    context.Context
}

func NewRewardManager(db *database.DB) (*RewardManager, error) {
	rewardManagerAbi, err := rm.RewardManagerMetaData.GetAbi()
	if err != nil {
		log.Error("get delegate manager abi fail", "err", err)
		return nil, err
	}

	rewardManagerUnpack, err := rm.NewRewardManagerFilterer(common2.Address{}, nil)
	if err != nil {
		log.Error("new delegation manager fail", "err", err)
		return nil, err
	}

	return &RewardManager{
		db:       db,
		RmAbi:    rewardManagerAbi,
		RmFilter: rewardManagerUnpack,
	}, nil
}

func (rm *RewardManager) ProcessRewardManager(fromHeight *big.Int, toHeight *big.Int) error {
	contractEventFilter := event.ContractEvent{ContractAddress: common2.HexToAddress(config.RewardManagerAddress)}
	contractEventList, err := rm.db.ContractEvent.ContractEventsWithFilter(contractEventFilter, fromHeight, toHeight)
	if err != nil {
		log.Error("get contracts event list fail", "err", err)
		return err
	}

	operatorAndStakeRewardList := make([]event.OperatorAndStakeReward, 0, len(contractEventList))
	for _, eventItem := range contractEventList {
		rlpLog := eventItem.RLPLog
		if eventItem.EventSignature.String() == rm.RmAbi.Events["OperatorAndStakeReward"].ID.String() {
			operatorAndStakeRewardEvent, err := rm.RmFilter.ParseOperatorAndStakeReward(*rlpLog)
			if err != nil {
				log.Error("parse operator and stake reward event fail", "err", err)
				return err
			}
			log.Info("parse operator and stake reward success",
				"strategy", operatorAndStakeRewardEvent.Strategy.String(),
				"operator", operatorAndStakeRewardEvent.Operator.String(),
				"stakerFee", operatorAndStakeRewardEvent.StakerFee.String(),
				"operatorFee", operatorAndStakeRewardEvent.OperatorFee.String())

			header, err := rm.db.Blocks.BlockHeader(eventItem.BlockHash)
			if err != nil {
				log.Error("db Blocks BlockHeader by BlockHash fail", "err", err)
				return err
			}

			temp := event.OperatorAndStakeReward{
				GUID:        uuid.New(),
				BlockHash:   eventItem.BlockHash,
				Number:      header.Number,
				TxHash:      eventItem.TransactionHash,
				Strategy:    operatorAndStakeRewardEvent.Strategy,
				Operator:    operatorAndStakeRewardEvent.Operator,
				StakerFee:   operatorAndStakeRewardEvent.StakerFee,
				OperatorFee: operatorAndStakeRewardEvent.OperatorFee,
				IsHandle:    0,
				Timestamp:   eventItem.Timestamp,
			}
			operatorAndStakeRewardList = append(operatorAndStakeRewardList, temp)
		}

		// OperatorClaimReward
		if eventItem.EventSignature.String() == rm.RmAbi.Events["OperatorClaimReward"].ID.String() {
			operatorClaimRewardEvent, err := rm.RmFilter.ParseOperatorClaimReward(*rlpLog)
			if err != nil {
				log.Error("parse operator claim reward event fail", "err", err)
				return err
			}
			log.Info("parse operator claim reward success",
				"operator", operatorClaimRewardEvent.Operator.String(),
				"amount", operatorClaimRewardEvent.Amount.String())
		}

		// StakeHolderClaimReward
		if eventItem.EventSignature.String() == rm.RmAbi.Events["StakeHolderClaimReward"].ID.String() {
			stakeHolderClaimRewardEvent, err := rm.RmFilter.ParseStakeHolderClaimReward(*rlpLog)
			if err != nil {
				log.Error("parse stake holder claim reward event fail", "err", err)
				return err
			}
			log.Info("parse stake holder claim reward success",
				"stakeHolder", stakeHolderClaimRewardEvent.StakeHolder.String(),
				"strategy", stakeHolderClaimRewardEvent.Strategy.String(),
				"amount", stakeHolderClaimRewardEvent.Amount.String())
		}
	}

	retryStrategy := &retry.ExponentialStrategy{Min: 1000, Max: 20_000, MaxJitter: 250}
	if _, err := retry.Do[interface{}](rm.rmCtx, 10, retryStrategy, func() (interface{}, error) {
		if err := rm.db.Transaction(func(tx *database.DB) error {
			if err := tx.OperatorAndStakeReward.StoreOperatorAndStakeReward(operatorAndStakeRewardList); err != nil {
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

	return nil
}
