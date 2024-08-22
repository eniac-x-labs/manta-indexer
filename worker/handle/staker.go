package handle

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"

	"github.com/eniac-x-labs/manta-indexer/common/tasks"
	"github.com/eniac-x-labs/manta-indexer/database"
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
	log.Info("...starting staker holder...")
	tickerOperator := time.NewTicker(time.Second * 5)
	sh.tasks.Go(func() error {
		for range tickerOperator.C {

		}
		return nil
	})
	return nil
}
