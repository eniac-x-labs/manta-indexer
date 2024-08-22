package manta_indexer

import (
	"context"
	"math/big"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/log"

	"github.com/eniac-x-labs/manta-indexer/config"
	"github.com/eniac-x-labs/manta-indexer/database"
	"github.com/eniac-x-labs/manta-indexer/event"
	"github.com/eniac-x-labs/manta-indexer/synchronizer"
	"github.com/eniac-x-labs/manta-indexer/synchronizer/node"
)

type MantaIndexer struct {
	synchronizer   *synchronizer.Synchronizer
	eventProcessor *event.EventProcessor

	shutdown context.CancelCauseFunc
	stopped  atomic.Bool
}

func NewMantaIndexer(ctx context.Context, cfg *config.Config, shutdown context.CancelCauseFunc) (*MantaIndexer, error) {
	ethClient, err := node.DialEthClient(ctx, cfg.Chain.ChainRpcUrl)
	if err != nil {
		return nil, err
	}

	db, err := database.NewDB(ctx, cfg.MasterDB)
	if err != nil {
		log.Error("init database fail", err)
		return nil, err
	}

	syncer, err := synchronizer.NewSynchronizer(cfg, db, ethClient, shutdown)
	if err != nil {
		log.Error("new synchronizer fail", "err", err)
		return nil, err
	}

	eventConfigm := &event.EventProcessorConfig{
		LoopInterval:    cfg.Chain.LoopInterval,
		StartHeight:     big.NewInt(int64(cfg.Chain.StartingHeight)),
		EventStartBlock: cfg.Chain.StartingHeight,
		Epoch:           500,
	}

	eventProcessor, err := event.NewEventProcessor(db, eventConfigm, shutdown)
	if err != nil {
		log.Error("new event processor fail", "err", err)
		return nil, err
	}

	out := &MantaIndexer{
		synchronizer:   syncer,
		eventProcessor: eventProcessor,
		shutdown:       shutdown,
	}
	return out, nil
}

func (ew *MantaIndexer) Start(ctx context.Context) error {
	err := ew.synchronizer.Start()
	if err != nil {
		return err
	}
	err = ew.eventProcessor.Start()
	if err != nil {
		return err
	}
	return nil
}

func (ew *MantaIndexer) Stop(ctx context.Context) error {
	err := ew.synchronizer.Close()
	if err != nil {
		return err
	}
	err = ew.eventProcessor.Close()
	if err != nil {
		return err
	}
	return nil
}

func (ew *MantaIndexer) Stopped() bool {
	return ew.stopped.Load()
}
