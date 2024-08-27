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
	"github.com/eniac-x-labs/manta-indexer/database/event/operator"
	"github.com/eniac-x-labs/manta-indexer/database/event/staker"
	"github.com/eniac-x-labs/manta-indexer/synchronizer/retry"
)

type RewardManager struct {
	RmAbi    *abi.ABI
	RmFilter *rm.RewardManagerFilterer
	rmCtx    context.Context
}

func NewRewardManager() (*RewardManager, error) {
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
		RmAbi:    rewardManagerAbi,
		RmFilter: rewardManagerUnpack,
		rmCtx:    context.Background(),
	}, nil
}

func (rm *RewardManager) ProcessRewardManager(db *database.DB, fromHeight *big.Int, toHeight *big.Int) error {
	contractEventFilter := event.ContractEvent{ContractAddress: common2.HexToAddress(config.RewardManagerAddress)}
	contractEventList, err := db.ContractEvent.ContractEventsWithFilter(contractEventFilter, fromHeight, toHeight)
	if err != nil {
		log.Error("get contracts event list fail", "err", err)
		return err
	}

	var operatorAndStakeRewardList []operator.OperatorAndStakeReward
	var operatorClaimRewardList []operator.OperatorClaimReward
	var stakeHolderClaimRewardList []staker.StakeHolderClaimReward
	for _, eventItem := range contractEventList {
		rlpLog := eventItem.RLPLog

		header, err := db.Blocks.BlockHeader(eventItem.BlockHash)
		if err != nil {
			log.Error("ProcessRewardManager db Blocks BlockHeader by BlockHash fail", "err", err)
			return err
		}

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

			temp := operator.OperatorAndStakeReward{
				GUID:             uuid.New(),
				BlockHash:        eventItem.BlockHash,
				Number:           header.Number,
				TxHash:           eventItem.TransactionHash,
				Strategy:         operatorAndStakeRewardEvent.Strategy,
				Operator:         operatorAndStakeRewardEvent.Operator,
				StakerFee:        operatorAndStakeRewardEvent.StakerFee,
				OperatorFee:      operatorAndStakeRewardEvent.OperatorFee,
				IsOperatorHandle: 0,
				IsStakerHandle:   0,
				Timestamp:        eventItem.Timestamp,
			}
			operatorAndStakeRewardList = append(operatorAndStakeRewardList, temp)
		}

		if eventItem.EventSignature.String() == rm.RmAbi.Events["OperatorClaimReward"].ID.String() {
			operatorClaimRewardEvent, err := rm.RmFilter.ParseOperatorClaimReward(*rlpLog)
			if err != nil {
				log.Error("parse operator claim reward event fail", "err", err)
				return err
			}
			log.Info("parse operator claim reward success",
				"operator", operatorClaimRewardEvent.Operator.String(),
				"amount", operatorClaimRewardEvent.Amount.String())

			temp := operator.OperatorClaimReward{
				GUID:      uuid.New(),
				BlockHash: eventItem.BlockHash,
				Number:    header.Number,
				TxHash:    eventItem.TransactionHash,
				Operator:  operatorClaimRewardEvent.Operator,
				Amount:    operatorClaimRewardEvent.Amount,
				IsHandle:  0,
				Timestamp: eventItem.Timestamp,
			}
			operatorClaimRewardList = append(operatorClaimRewardList, temp)
		}

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

			temp := staker.StakeHolderClaimReward{
				GUID:        uuid.New(),
				BlockHash:   eventItem.BlockHash,
				Number:      header.Number,
				TxHash:      eventItem.TransactionHash,
				StakeHolder: stakeHolderClaimRewardEvent.StakeHolder,
				Strategy:    stakeHolderClaimRewardEvent.Strategy,
				Amount:      stakeHolderClaimRewardEvent.Amount,
				IsHandle:    0,
				Timestamp:   eventItem.Timestamp,
			}
			stakeHolderClaimRewardList = append(stakeHolderClaimRewardList, temp)
		}
	}

	retryStrategy := &retry.ExponentialStrategy{Min: 1000, Max: 20_000, MaxJitter: 250}
	if _, err := retry.Do[interface{}](rm.rmCtx, 10, retryStrategy, func() (interface{}, error) {
		if err := db.Transaction(func(tx *database.DB) error {
			if err := tx.OperatorAndStakeReward.StoreOperatorAndStakeReward(operatorAndStakeRewardList); err != nil {
				return err
			}
			if err := tx.OperatorClaimReward.StoreOperatorClaimReward(operatorClaimRewardList); err != nil {
				return err
			}
			if err := tx.StakeHolderClaimReward.StoreStakeHolderClaimReward(stakeHolderClaimRewardList); err != nil {
				return err
			}
			// Log success messages
			log.Info("store reward manager events success",
				"operatorAndStakeRewardList", len(operatorAndStakeRewardList),
				"operatorClaimRewardList", len(operatorClaimRewardList),
				"stakeHolderClaimRewardList", len(stakeHolderClaimRewardList),
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
