package handle

import (
	"context"
	"fmt"
	"strings"

	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/log"

	"github.com/eniac-x-labs/manta-indexer/common/tasks"
	"github.com/eniac-x-labs/manta-indexer/database"
	"github.com/eniac-x-labs/manta-indexer/database/event/staker"
	"github.com/eniac-x-labs/manta-indexer/database/event/strategies"
	"github.com/eniac-x-labs/manta-indexer/database/worker"
	"github.com/eniac-x-labs/manta-indexer/synchronizer/retry"
)

type StakeHolderHandle struct {
	db             *database.DB
	resourceCtx    context.Context
	resourceCancel context.CancelFunc
	tasks          tasks.Group
}

func NewStakeHolderHandle(db *database.DB, shutdown context.CancelCauseFunc) (*StakeHolderHandle, error) {
	resCtx, resCancel := context.WithCancel(context.Background())
	return &StakeHolderHandle{
		db:             db,
		resourceCtx:    resCtx,
		resourceCancel: resCancel,
		tasks: tasks.Group{HandleCrit: func(err error) {
			shutdown(fmt.Errorf("critical error in bridge processor: %w", err))
		}},
	}, nil
}

func (sh *StakeHolderHandle) Close() error {
	sh.resourceCancel()
	return sh.tasks.Wait()
}

func (sh *StakeHolderHandle) Start() error {
	log.Info("=======================================================")
	log.Info("===========start stake holder worker task===========")
	log.Info("=======================================================")
	tickerOperator := time.NewTicker(time.Second * 5)
	sh.tasks.Go(func() error {
		for range tickerOperator.C {
			err := sh.processStrategyDeposit()
			if err != nil {
				log.Error("Process strategy deposit fail", "err", err)
				return err
			}
			err = sh.processWithdrawalCompleted()
			if err != nil {
				log.Error("Process withdraw completed fail", "err", err)
				return err
			}
		}
		return nil
	})
	return nil
}

func (sh *StakeHolderHandle) processStrategyDeposit() error {
	unHandleDepositList, err := sh.db.StrategyDeposit.QueryUnHandleStrategyDeposit()
	if err != nil {
		log.Error("Query unhandled strategy deposit fail", "err", err)
		return err
	}
	var strategyList []strategies.StrategyType
	for _, unHandleDeposit := range unHandleDepositList {
		stkType := worker.StakeHolderType{
			MantaStake:    unHandleDeposit.Shares,
			Reward:        big.NewInt(0),
			ClaimedAmount: big.NewInt(0),
			Timestamp:     unHandleDeposit.Timestamp,
		}
		strategy := strategies.StrategyType{
			Strategy: unHandleDeposit.Strategy.String(),
			Tvl:      unHandleDeposit.Shares,
		}
		strategyList = append(strategyList, strategy)
		log.Info("process strategy deposit query and update stake holder", "Staker", unHandleDeposit.Staker.String(), "Strategy", unHandleDeposit.Strategy.String())
		err := sh.db.StakeHolder.QueryAndUpdateStakeHolder(unHandleDeposit.Staker.String(), unHandleDeposit.Strategy.String(), stkType)
		if err != nil {
			log.Error("process strategy deposit query and update operator fail", "err", err)
			return err
		}
	}
	log.Info("process strategy deposit", "unHandleDepositList", len(unHandleDepositList), "strategyList", len(strategyList))

	retryStrategy := &retry.ExponentialStrategy{Min: 1000, Max: 20_000, MaxJitter: 250}
	if _, err := retry.Do[interface{}](sh.resourceCtx, 10, retryStrategy, func() (interface{}, error) {
		if err := sh.db.Transaction(func(tx *database.DB) error {
			if len(strategyList) > 0 {
				if err := tx.Strategies.UpdateStrategyTvlHandled(strategyList); err != nil {
					log.Error("Update strategy tvl handled fail", "err", err)
					return err
				}
			}

			if len(unHandleDepositList) > 0 {
				if err := tx.StrategyDeposit.MarkedStrategyDepositHandled(unHandleDepositList); err != nil {
					log.Error("Marked strategy deposit handled fail", "err", err)
					return err
				}
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

func (sh *StakeHolderHandle) processWithdrawalCompleted() error {
	unHandleCompletedList, err := sh.db.WithdrawalCompleted.QueryUnHandleWithdrawalCompleted()
	if err != nil {
		return err
	}

	var withdrawalQueuedList []staker.WithdrawalQueuedType
	var strategyList []strategies.StrategyType
	for _, unHandleCompleted := range unHandleCompletedList {
		stkType := worker.StakeHolderType{
			MantaStake:    new(big.Int).Neg(unHandleCompleted.Shares),
			Reward:        big.NewInt(0),
			ClaimedAmount: big.NewInt(0),
			Timestamp:     unHandleCompleted.Timestamp,
		}
		err := sh.db.StakeHolder.QueryAndUpdateStakeHolder(unHandleCompleted.Staker.String(), unHandleCompleted.Strategy.String(), stkType)
		if err != nil {
			log.Error("process withdrawal completed query and update staker fail", "err", err)
			return err
		}
		withdrawalQueued := staker.WithdrawalQueuedType{
			Staker:     strings.ToLower(unHandleCompleted.Staker.String()),
			Strategies: strings.ToLower(unHandleCompleted.Strategy.String()),
		}
		withdrawalQueuedList = append(withdrawalQueuedList, withdrawalQueued)

		strategy := strategies.StrategyType{
			Strategy: unHandleCompleted.Strategy.String(),
			Tvl:      new(big.Int).Neg(unHandleCompleted.Shares),
		}
		strategyList = append(strategyList, strategy)
	}
	log.Info("process withdrawal completed", "unHandleCompletedList number", len(unHandleCompletedList), "withdrawalQueuedList number", len(withdrawalQueuedList))

	retryStrategy := &retry.ExponentialStrategy{Min: 1000, Max: 20_000, MaxJitter: 250}
	if _, err := retry.Do[interface{}](sh.resourceCtx, 10, retryStrategy, func() (interface{}, error) {
		if err := sh.db.Transaction(func(tx *database.DB) error {
			if len(strategyList) > 0 {
				if err := tx.Strategies.UpdateStrategyTvlHandled(strategyList); err != nil {
					log.Error("update strategy tvl handled fail", "err", err)
					return err
				}
			}

			if len(withdrawalQueuedList) > 0 {
				if err := tx.WithdrawalQueued.MarkedWithdrawalQueuedHandled(withdrawalQueuedList); err != nil {
					log.Error("marked withdrawal queued handled fail", "err", err)
					return err
				}
			}

			if len(unHandleCompletedList) > 0 {
				if err := tx.WithdrawalCompleted.MarkedWithdrawalCompleted(unHandleCompletedList); err != nil {
					log.Error("marked withdrawal completed fail", "err", err)
					return err
				}
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

func (sh *StakeHolderHandle) processStakeHolderClaimReward() error {
	unHandleStakeHolderClaimRewardList, err := sh.db.StakeHolderClaimReward.QueryUnHandleStakeHolderClaimReward()
	if err != nil {
		return err
	}
	for _, unHandleStakeHolderClaimReward := range unHandleStakeHolderClaimRewardList {
		stkType := worker.StakeHolderType{
			MantaStake:    big.NewInt(0),
			Reward:        big.NewInt(0),
			ClaimedAmount: new(big.Int).Neg(unHandleStakeHolderClaimReward.Amount),
			Timestamp:     unHandleStakeHolderClaimReward.Timestamp,
		}
		err := sh.db.StakeHolder.QueryAndUpdateStakeHolder(unHandleStakeHolderClaimReward.StakeHolder.String(), unHandleStakeHolderClaimReward.Strategy.String(), stkType)
		if err != nil {
			log.Error("process withdrawal completed query and update staker fail", "err", err)
			return err
		}
	}

	retryStrategy := &retry.ExponentialStrategy{Min: 1000, Max: 20_000, MaxJitter: 250}
	if _, err := retry.Do[interface{}](sh.resourceCtx, 10, retryStrategy, func() (interface{}, error) {
		if err := sh.db.Transaction(func(tx *database.DB) error {
			if len(unHandleStakeHolderClaimRewardList) > 0 {
				if err := tx.StakeHolderClaimReward.MarkedStakeHolderClaimRewardHandled(unHandleStakeHolderClaimRewardList); err != nil {
					log.Error("marked stake holder claim reward handled fail", "err", err)
					return err
				}
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
