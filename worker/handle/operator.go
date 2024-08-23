package handle

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"

	"github.com/ethereum/go-ethereum/log"

	"github.com/eniac-x-labs/manta-indexer/common/tasks"
	"github.com/eniac-x-labs/manta-indexer/database"
	"github.com/eniac-x-labs/manta-indexer/database/worker"
	"github.com/ethereum/go-ethereum/common"
)

type OperatorHandle struct {
	db             *database.DB
	resourceCtx    context.Context
	resourceCancel context.CancelFunc
	tasks          tasks.Group
}

func NewOperatorHandle(db *database.DB, shutdown context.CancelCauseFunc) (*OperatorHandle, error) {
	resCtx, resCancel := context.WithCancel(context.Background())
	return &OperatorHandle{
		db:             db,
		resourceCtx:    resCtx,
		resourceCancel: resCancel,
		tasks: tasks.Group{HandleCrit: func(err error) {
			shutdown(fmt.Errorf("critical error in bridge processor: %w", err))
		}},
	}, nil
}

func (oh *OperatorHandle) Start() error {
	log.Info("=======================================================")
	log.Info("============starting operator worker task==============")
	log.Info("=======================================================")
	tickerOperator := time.NewTicker(time.Second * 5)
	oh.tasks.Go(func() error {
		for range tickerOperator.C {
			err := oh.processOperator()
			if err != nil {
				return err
			}

			err = oh.processOperatorNodeUrlUpdate()
			if err != nil {
				return err
			}

			err = oh.processOperatorSharesDecreased()
			if err != nil {
				return err
			}

			err = oh.processOperatorSharesIncreased()
			if err != nil {
				return err
			}

			err = oh.processOperatorRewardIncreased()
			if err != nil {
				return err
			}

			err = oh.processOperatorRewardDecreased()
			if err != nil {
				return err
			}
		}
		return nil
	})
	return nil
}

func (oh *OperatorHandle) Close() error {
	oh.resourceCancel()
	return oh.tasks.Wait()
}

func (oh *OperatorHandle) processOperator() error {
	unHandleRegisteredList, err := oh.db.OperatorRegistered.QueryUnHandleOperatorRegistered()
	if err != nil {
		log.Error("QueryUnHandleOperatorRegistered fail", "err", err)
		return err
	}
	operators := make([]worker.Operators, 0, len(unHandleRegisteredList))
	for _, unHandleRegistered := range unHandleRegisteredList {
		operator := worker.Operators{
			GUID:                     uuid.New(),
			BlockHash:                unHandleRegistered.BlockHash,
			Number:                   unHandleRegistered.Number,
			TxHash:                   unHandleRegistered.TxHash,
			Operator:                 unHandleRegistered.Operator,
			Socket:                   "unknown",
			EarningsReceiver:         unHandleRegistered.EarningsReceiver,
			DelegationApprover:       unHandleRegistered.DelegationApprover,
			StakerOptoutWindowBlocks: unHandleRegistered.StakerOptoutWindowBlocks,
			TotalMantaStake:          big.NewInt(0),
			TotalStakeReward:         big.NewInt(0),
			RateReturn:               "0",
			Status:                   0,
			Timestamp:                unHandleRegistered.Timestamp,
		}
		operators = append(operators, operator)
	}
	if len(operators) > 0 {
		if err := oh.db.Operators.StoreOperators(operators); err != nil {
			log.Error("StoreOperators fail", "err", err)
			return err
		}
		if err := oh.db.OperatorRegistered.MarkedOperatorRegisteredHandled(unHandleRegisteredList); err != nil {
			log.Error("MarkedOperatorRegisteredHandled fail", "err", err)
			return err
		}
	}
	return nil
}

func (oh *OperatorHandle) processOperatorNodeUrlUpdate() error {
	unHandleUrlUpdateList, err := oh.db.OperatorNodeUrlUpdate.QueryUnHandleOperatorNodeUrlUpdate()
	if err != nil {
		log.Error("QueryUnHandleOperatorNodeUrlUpdate fail", "err", err)
		return err
	}
	for _, unHandleUrlUpdate := range unHandleUrlUpdateList {
		opType := worker.OperatorsType{
			Socket:                   unHandleUrlUpdate.MetadataUri,
			EarningsReceiver:         common.Address{},
			DelegationApprover:       common.Address{},
			StakerOptoutWindowBlocks: nil,
			TotalMantaStake:          nil,
			TotalStakeReward:         nil,
			RateReturn:               "",
			Status:                   0,
		}
		err := oh.db.Operators.QueryAndUpdateOperator(unHandleUrlUpdate.Operator, opType)
		if err != nil {
			log.Error("QueryAndUpdateOperator fail", "err", err)
			return err
		}
	}
	if len(unHandleUrlUpdateList) > 0 {
		if err := oh.db.OperatorNodeUrlUpdate.MarkedOperatorNodeUrlUpdateHandled(unHandleUrlUpdateList); err != nil {
			log.Error("MarkedOperatorNodeUrlUpdateHandled fail", "err", err)
			return err
		}
	}
	return nil
}

func (oh *OperatorHandle) processOperatorSharesDecreased() error {
	decreasedList, err := oh.db.OperatorSharesDecreased.QueryUnHandlerOperatorSharesDecreased()
	if err != nil {
		log.Error("QueryUnHandlerOperatorSharesDecreased fail", "err", err)
		return err
	}
	for _, decreased := range decreasedList {
		opType := worker.OperatorsType{
			Socket:                   "",
			EarningsReceiver:         common.Address{},
			DelegationApprover:       common.Address{},
			StakerOptoutWindowBlocks: nil,
			TotalMantaStake:          new(big.Int).Neg(decreased.Shares),
			TotalStakeReward:         nil,
			RateReturn:               "",
			Status:                   0,
		}
		err := oh.db.Operators.QueryAndUpdateOperator(decreased.Operator, opType)
		if err != nil {
			log.Error("QueryAndUpdateOperator fail", "err", err)
			return err
		}
	}
	if len(decreasedList) > 0 {
		if err := oh.db.OperatorSharesDecreased.MarkedOperatorSharesDecreasedHandled(decreasedList); err != nil {
			log.Error("MarkedOperatorSharesDecreasedHandled fail", "err", err)
			return err
		}
	}
	return nil
}

func (oh *OperatorHandle) processOperatorSharesIncreased() error {
	increasedList, err := oh.db.OperatorSharesIncreased.QueryUnHandleOperatorSharesIncreased()
	if err != nil {
		log.Error("QueryUnHandleOperatorSharesIncreased fail", "err", err)
		return err
	}
	for _, increased := range increasedList {
		opType := worker.OperatorsType{
			Socket:                   "",
			EarningsReceiver:         common.Address{},
			DelegationApprover:       common.Address{},
			StakerOptoutWindowBlocks: nil,
			TotalMantaStake:          increased.Shares,
			TotalStakeReward:         nil,
			RateReturn:               "",
			Status:                   0,
		}
		err := oh.db.Operators.QueryAndUpdateOperator(increased.Operator, opType)
		if err != nil {
			log.Error("QueryAndUpdateOperator fail", "err", err)
			return err
		}
	}
	if len(increasedList) > 0 {
		if err := oh.db.OperatorSharesIncreased.MarkedOperatorSharesIncreasedHandled(increasedList); err != nil {
			log.Error("MarkedOperatorSharesIncreasedHandled fail", "err", err)
			return err
		}
	}
	return nil
}

func (oh *OperatorHandle) processOperatorRewardIncreased() error {
	operatorAndStakeRewardList, err := oh.db.OperatorAndStakeReward.QueryUnHandleOperatorAndStakeReward(true)
	if err != nil {
		log.Error("QueryUnHandleOperatorAndStakeReward fail", "err", err)
		return err
	}
	for _, operatorAndStakeReward := range operatorAndStakeRewardList {
		opType := worker.OperatorsType{
			Socket:                   "",
			EarningsReceiver:         common.Address{},
			DelegationApprover:       common.Address{},
			StakerOptoutWindowBlocks: nil,
			TotalMantaStake:          operatorAndStakeReward.OperatorFee,
			TotalStakeReward:         nil,
			RateReturn:               "",
			Status:                   0,
		}
		err := oh.db.Operators.QueryAndUpdateOperator(operatorAndStakeReward.Operator, opType)
		if err != nil {
			log.Error("QueryAndUpdateOperator fail", "err", err)
			return err
		}
	}
	if len(operatorAndStakeRewardList) > 0 {
		if err := oh.db.OperatorAndStakeReward.MarkedOperatorAndStakeRewardHandled(operatorAndStakeRewardList, true); err != nil {
			log.Error("MarkedOperatorAndStakeRewardHandled fail", "err", err)
			return err
		}
	}
	return nil
}

func (oh *OperatorHandle) processOperatorRewardDecreased() error {
	operatorClaimRewardList, err := oh.db.OperatorClaimReward.QueryUnHandleOperatorClaimReward()
	if err != nil {
		log.Error("QueryUnHandleOperatorClaimReward fail", "err", err)
		return err
	}
	for _, operatorClaimReward := range operatorClaimRewardList {
		opType := worker.OperatorsType{
			Socket:                   "",
			EarningsReceiver:         common.Address{},
			DelegationApprover:       common.Address{},
			StakerOptoutWindowBlocks: nil,
			TotalMantaStake:          new(big.Int).Neg(operatorClaimReward.Amount),
			TotalStakeReward:         nil,
			RateReturn:               "",
			Status:                   0,
		}
		err := oh.db.Operators.QueryAndUpdateOperator(operatorClaimReward.Operator, opType)
		if err != nil {
			log.Error("QueryAndUpdateOperator fail", "err", err)
			return err
		}
	}
	if len(operatorClaimRewardList) > 0 {
		if err := oh.db.OperatorClaimReward.MarkedOperatorClaimRewardHandled(operatorClaimRewardList); err != nil {
			log.Error("MarkedOperatorClaimRewardHandled fail", "err", err)
			return err
		}
	}
	return nil
}
