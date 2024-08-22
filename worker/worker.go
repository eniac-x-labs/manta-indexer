package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"

	"github.com/eniac-x-labs/manta-indexer/common/tasks"
	"github.com/eniac-x-labs/manta-indexer/database"
)

type WorkerConfig struct {
	LoopInterval time.Duration
}

type Worker struct {
	db             *database.DB
	workerConfig   *WorkerConfig
	resourceCtx    context.Context
	resourceCancel context.CancelFunc
	tasks          tasks.Group
}

func NewWorker(db *database.DB, workerConfig *WorkerConfig, shutdown context.CancelCauseFunc) (*WorkerHandler, error) {
	resCtx, resCancel := context.WithCancel(context.Background())
	return &Worker{
		db:             db,
		resourceCtx:    resCtx,
		resourceCancel: resCancel,
		workerConfig:   workerConfig,
		tasks: tasks.Group{HandleCrit: func(err error) {
			shutdown(fmt.Errorf("critical error in bridge processor: %w", err))
		}},
	}, nil
}

func (ep *Worker) Start() error {
	log.Info("...starting worker...")
	tickerOperator := time.NewTicker(time.Second * 5)
	ep.tasks.Go(func() error {
		for range tickerOperator.C {
			err := ep.processOperator()
			if err != nil {
				return err
			}
		}
		return nil
	})

	tickerStakeHolder := time.NewTicker(time.Second * 5)
	ep.tasks.Go(func() error {
		for range tickerStakeHolder.C {
			err := ep.processStakeHolder()
			if err != nil {
				return err
			}
		}
		return nil
	})
	return nil
}

func (ep *Worker) Close() error {
	ep.resourceCancel()
	return ep.tasks.Wait()
}

func (ep *Worker) processOperator() error {
	return nil
}

func (ep *Worker) processStakeHolder() error {
	return nil
}
