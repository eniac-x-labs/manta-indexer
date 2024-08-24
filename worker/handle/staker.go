package handle

import (
	"context"
	"fmt"
	"github.com/eniac-x-labs/manta-indexer/database/event/strategies"

	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/log"

	"github.com/eniac-x-labs/manta-indexer/common/tasks"
	"github.com/eniac-x-labs/manta-indexer/database"
	"github.com/eniac-x-labs/manta-indexer/database/worker"
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
	var strategyList []strategies.Strategies
	for _, unHandleDeposit := range unHandleDepositList {
		stkType := worker.StakeHolderType{
			MantaStake:    unHandleDeposit.Shares,
			Reward:        big.NewInt(0),
			ClaimedAmount: big.NewInt(0),
			Timestamp:     unHandleDeposit.Timestamp,
		}
		strategy := strategies.Strategies{
			Strategy: unHandleDeposit.Strategy,
			Tvl:      unHandleDeposit.Shares,
		}
		strategyList = append(strategyList, strategy)
		err := sh.db.StakeHolder.QueryAndUpdateStakeHolder(unHandleDeposit.Staker.String(), unHandleDeposit.Strategy.String(), stkType)
		if err != nil {
			log.Error("processStrategyDeposit query and update operator fail", "err", err)
			return err
		}
	}
	if len(unHandleDepositList) > 0 {
		if err := sh.db.StrategyDeposit.MarkedStrategyDepositHandled(unHandleDepositList); err != nil {
			log.Error("MarkedStrategyDepositHandled fail", "err", err)
			return err
		}
	}
	if len(strategyList) > 0 {
		if err := sh.db.Strategies.UpdateStrategyTvlHandled(strategyList); err != nil {
			log.Error("UpdateStrategyTvlHandled fail", "err", err)
			return err
		}
	}
	return nil
}

func (sh *StakeHolderHandle) processWithdrawalCompleted() error {
	unHandleCompletedList, err := sh.db.WithdrawalCompleted.QueryUnHandleWithdrawalCompleted()
	if err != nil {
		return err
	}
	for _, unHandleCompleted := range unHandleCompletedList {
		stkType := worker.StakeHolderType{
			MantaStake:    new(big.Int).Neg(unHandleCompleted.Shares),
			Reward:        big.NewInt(0),
			ClaimedAmount: big.NewInt(0),
			Timestamp:     unHandleCompleted.Timestamp,
		}
		err := sh.db.StakeHolder.QueryAndUpdateStakeHolder(unHandleCompleted.Staker.String(), unHandleCompleted.Strategy.String(), stkType)
		if err != nil {
			log.Error("processWithdrawalCompleted query and update staker fail", "err", err)
			return err
		}
	}
	if len(unHandleCompletedList) > 0 {
		if err := sh.db.WithdrawalCompleted.MarkedWithdrawalCompleted(unHandleCompletedList); err != nil {
			log.Error("MarkedWithdrawalCompleted fail", "err", err)
			return err
		}
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
			log.Error("processWithdrawalCompleted query and update staker fail", "err", err)
			return err
		}
	}
	if len(unHandleStakeHolderClaimRewardList) > 0 {
		if err := sh.db.StakeHolderClaimReward.MarkedStakeHolderClaimRewardHandled(unHandleStakeHolderClaimRewardList); err != nil {
			log.Error("MarkedStakeHolderClaimRewardHandled fail", "err", err)
			return err
		}
	}
	return nil
}
