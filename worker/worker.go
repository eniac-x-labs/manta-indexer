package worker

import (
	"context"
	"github.com/ethereum/go-ethereum/log"

	"github.com/eniac-x-labs/manta-indexer/database"
	"github.com/eniac-x-labs/manta-indexer/worker/handle"
)

type Worker struct {
	operatorHandle    *handle.OperatorHandle
	stakeHolderHandle *handle.StakeHolderHandle
}

func NewWorker(db *database.DB, shutdown context.CancelCauseFunc) (*Worker, error) {
	operatorHandle, err := handle.NewOperatorHandle(db, shutdown)
	if err != nil {
		log.Error("New operator fail", "err", err)
		return nil, err
	}

	stakeHolderHandle, err := handle.NewStakeHolderHandle(db, shutdown)
	if err != nil {
		log.Error("New stake holder handle fail", "err", err)
		return nil, err
	}

	return &Worker{
		operatorHandle:    operatorHandle,
		stakeHolderHandle: stakeHolderHandle,
	}, nil
}

func (ep *Worker) Start() error {
	log.Info("...starting worker...")
	err := ep.operatorHandle.Start()
	if err != nil {
		log.Error("start operator handler fail", "err", err)
		return err
	}
	err = ep.stakeHolderHandle.Start()
	if err != nil {
		log.Error("start stake holder handler fail", "err", err)
		return err
	}
	return nil
}

func (ep *Worker) Close() error {
	err := ep.operatorHandle.Close()
	if err != nil {
		return err
	}
	err = ep.stakeHolderHandle.Close()
	if err != nil {
		return err
	}
	return nil
}
