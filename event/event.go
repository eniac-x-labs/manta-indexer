package event

import (
	"context"
	"fmt"
	"github.com/eniac-x-labs/manta-indexer/database/event"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/log"

	"github.com/eniac-x-labs/manta-indexer/common/bigint"
	"github.com/eniac-x-labs/manta-indexer/common/tasks"
	"github.com/eniac-x-labs/manta-indexer/database"
	"github.com/eniac-x-labs/manta-indexer/database/common"
	"github.com/eniac-x-labs/manta-indexer/event/contracts"
)

var blocksLimit = 10_000

type EventProcessorConfig struct {
	LoopInterval    time.Duration
	StartHeight     *big.Int
	EventStartBlock uint64
	Epoch           uint64
}

type EventProcessor struct {
	db                *database.DB
	eventBlocksConfig *EventProcessorConfig
	resourceCtx       context.Context
	resourceCancel    context.CancelFunc
	tasks             tasks.Group
	dtManager         *contracts.DelegationManager
	LatestBlockHeader *common.BlockHeader
}

func NewEventProcessor(db *database.DB, eventBlocksConfig *EventProcessorConfig, shutdown context.CancelCauseFunc) (*EventProcessor, error) {
	LatestBlockHeader, err := db.EventBlocks.LatestEventBlockHeader()
	if err != nil {
		return nil, err
	}

	resCtx, resCancel := context.WithCancel(context.Background())

	return &EventProcessor{
		db:                db,
		resourceCtx:       resCtx,
		resourceCancel:    resCancel,
		eventBlocksConfig: eventBlocksConfig,
		tasks: tasks.Group{HandleCrit: func(err error) {
			shutdown(fmt.Errorf("critical error in bridge processor: %w", err))
		}},
		LatestBlockHeader: LatestBlockHeader,
	}, nil
}

func (ep *EventProcessor) Start() error {
	log.Info("starting bridge processor...")
	tickerL1Worker := time.NewTicker(time.Second * 5)
	ep.tasks.Go(func() error {
		for range tickerL1Worker.C {
			err := ep.processEvent()
			if err != nil {
				return err
			}
		}
		return nil
	})
	return nil
}

func (ep *EventProcessor) Close() error {
	ep.resourceCancel()
	return ep.tasks.Wait()
}

func (ep *EventProcessor) processEvent() error {
	lastBlockNumber := big.NewInt(int64(ep.eventBlocksConfig.EventStartBlock))
	if ep.LatestBlockHeader != nil {
		lastBlockNumber = ep.LatestBlockHeader.Number
	}
	log.Info("Process event latest block number", "lastBlockNumber", lastBlockNumber)
	latestHeaderScope := func(db *gorm.DB) *gorm.DB {
		newQuery := db.Session(&gorm.Session{NewDB: true})
		headers := newQuery.Model(common.BlockHeader{}).Where("number > ?", lastBlockNumber)
		return db.Where("number = (?)", newQuery.Table("(?) as block_numbers", headers.Order("number ASC").Limit(blocksLimit)).Select("MAX(number)"))
	}
	if latestHeaderScope == nil {
		return nil
	}
	latestHeader, err := ep.db.Blocks.BlockHeaderWithScope(latestHeaderScope)
	if err != nil {
		return fmt.Errorf("failed to query for latest unfinalized L1 state: %w", err)
	} else if latestHeader == nil {
		log.Debug("no new  state to process event")
		return nil
	}
	fromHeight, toHeight := new(big.Int).Add(lastBlockNumber, bigint.One), latestHeader.Number
	eventBlocks := make([]event.EventBlocks, 0, toHeight.Uint64()-fromHeight.Uint64())
	for index := fromHeight.Uint64(); index < toHeight.Uint64(); index++ {
		blockHeader, err := ep.db.Blocks.BlockHeaderByNumber(big.NewInt(int64(index)))
		if err != nil {
			return err
		}
		evBlock := event.EventBlocks{
			GUID:       uuid.New(),
			Hash:       blockHeader.Hash,
			ParentHash: blockHeader.ParentHash,
			Number:     blockHeader.Number,
			Timestamp:  blockHeader.Timestamp,
		}
		eventBlocks = append(eventBlocks, evBlock)
	}

	dtManager, err := contracts.NewDelegationManager(ep.db, fromHeight, toHeight)
	dtManager.ProcessDelegationRegister()

	if err := ep.db.Transaction(func(tx *database.DB) error {
		log.Info("handle event from to end", "fromHeight", fromHeight, "toHeight", toHeight)
		err := ep.db.EventBlocks.StoreEventBlocks(eventBlocks)
		if err != nil {
			log.Error("store event block fail", "err", err)
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	ep.LatestBlockHeader = latestHeader
	return nil
}
