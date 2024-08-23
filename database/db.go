package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/eniac-x-labs/manta-indexer/config"
	"github.com/eniac-x-labs/manta-indexer/database/common"
	"github.com/eniac-x-labs/manta-indexer/database/event"
	_ "github.com/eniac-x-labs/manta-indexer/database/utils/serializers"
	"github.com/eniac-x-labs/manta-indexer/database/worker"
	"github.com/eniac-x-labs/manta-indexer/synchronizer/retry"
)

type DB struct {
	gorm                             *gorm.DB
	Blocks                           common.BlocksDB
	ContractEvent                    event.ContractEventDB
	EventBlocks                      event.EventBlocksDB
	MinWithdrawalDelayBlocksSet      event.MinWithdrawalDelayBlocksSetDB
	OperatorAndStakeReward           event.OperatorAndStakeRewardDB
	OperatorClaimReward              event.OperatorClaimRewardDB
	OperatorModified                 event.OperatorModifiedDB
	OperatorNodeUrlUpdate            event.OperatorNodeUrlUpdateDB
	OperatorRegistered               event.OperatorRegisteredDB
	OperatorSharesDecreased          event.OperatorSharesDecreasedDB
	OperatorSharesIncreased          event.OperatorSharesIncreasedDB
	StakeHolderClaimReward           event.StakeHolderClaimRewardDB
	StakerDelegated                  event.StakerDelegatedDB
	StakerUndelegated                event.StakerUndelegatedDB
	StrategyDeposit                  event.StrategyDepositDB
	StrategyWithdrawalDelayBlocksSet event.StrategyWithdrawalDelayBlocksSetDB
	WithdrawalMigrated               event.WithdrawalMigratedDB
	WithdrawalQueued                 event.WithdrawalQueuedDB
	WithdrawalCompleted              event.WithdrawalCompletedDB
	Operators                        worker.OperatorsDB
	OperatorPublicKeys               worker.OperatorPublicKeysDB
	StakeHolder                      worker.StakeHolderDB
	TotalOperator                    worker.TotalOperatorDB
}

func NewDB(ctx context.Context, dbConfig config.DBConfig) (*DB, error) {
	dsn := fmt.Sprintf("host=%s dbname=%s sslmode=disable", dbConfig.Host, dbConfig.Name)
	if dbConfig.Port != 0 {
		dsn += fmt.Sprintf(" port=%d", dbConfig.Port)
	}
	if dbConfig.User != "" {
		dsn += fmt.Sprintf(" user=%s", dbConfig.User)
	}
	if dbConfig.Password != "" {
		dsn += fmt.Sprintf(" password=%s", dbConfig.Password)
	}

	gormConfig := gorm.Config{
		SkipDefaultTransaction: true,
		CreateBatchSize:        3_000,
	}

	retryStrategy := &retry.ExponentialStrategy{Min: 1000, Max: 20_000, MaxJitter: 250}
	gorm, err := retry.Do[*gorm.DB](context.Background(), 10, retryStrategy, func() (*gorm.DB, error) {
		gorm, err := gorm.Open(postgres.Open(dsn), &gormConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}
		return gorm, nil
	})

	if err != nil {
		return nil, err
	}

	db := &DB{
		gorm:                             gorm,
		Blocks:                           common.NewBlocksDB(gorm),
		ContractEvent:                    event.NewContractEventsDB(gorm),
		EventBlocks:                      event.NewEventBlocksDB(gorm),
		MinWithdrawalDelayBlocksSet:      event.NewMinWithdrawalDelayBlocksSetDB(gorm),
		OperatorAndStakeReward:           event.NewOperatorAndStakeRewardDB(gorm),
		OperatorClaimReward:              event.NewOperatorClaimRewardDB(gorm),
		OperatorModified:                 event.NewOperatorModifiedDB(gorm),
		OperatorNodeUrlUpdate:            event.NewOperatorNodeUrlUpdateDB(gorm),
		OperatorRegistered:               event.NewOperatorRegisteredDB(gorm),
		OperatorSharesDecreased:          event.NewOperatorSharesDecreasedDB(gorm),
		OperatorSharesIncreased:          event.NewOperatorSharesIncreasedDB(gorm),
		StakeHolderClaimReward:           event.NewStakeHolderClaimRewardDB(gorm),
		StakerDelegated:                  event.NewStakerDelegatedDB(gorm),
		StakerUndelegated:                event.NewStakerUndelegatedDB(gorm),
		StrategyDeposit:                  event.NewStrategyDepositDB(gorm),
		StrategyWithdrawalDelayBlocksSet: event.NewStrategyWithdrawalDelayBlocksSetDB(gorm),
		WithdrawalMigrated:               event.NewWithdrawalMigratedDB(gorm),
		WithdrawalQueued:                 event.NewWithdrawalQueuedDB(gorm),
		WithdrawalCompleted:              event.NewWithdrawalCompletedDB(gorm),
		Operators:                        worker.NewOperatorsDB(gorm),
		OperatorPublicKeys:               worker.NewOperatorPublicKeysDB(gorm),
		StakeHolder:                      worker.NewStakeHolderDB(gorm),
		TotalOperator:                    worker.NewTotalOperatorDB(gorm),
	}
	return db, nil
}

func (db *DB) Transaction(fn func(db *DB) error) error {
	return db.gorm.Transaction(func(tx *gorm.DB) error {
		txDB := &DB{
			gorm:                             tx,
			Blocks:                           common.NewBlocksDB(tx),
			ContractEvent:                    event.NewContractEventsDB(tx),
			EventBlocks:                      event.NewEventBlocksDB(tx),
			MinWithdrawalDelayBlocksSet:      event.NewMinWithdrawalDelayBlocksSetDB(tx),
			OperatorAndStakeReward:           event.NewOperatorAndStakeRewardDB(tx),
			OperatorClaimReward:              event.NewOperatorClaimRewardDB(tx),
			OperatorModified:                 event.NewOperatorModifiedDB(tx),
			OperatorNodeUrlUpdate:            event.NewOperatorNodeUrlUpdateDB(tx),
			OperatorRegistered:               event.NewOperatorRegisteredDB(tx),
			OperatorSharesDecreased:          event.NewOperatorSharesDecreasedDB(tx),
			OperatorSharesIncreased:          event.NewOperatorSharesIncreasedDB(tx),
			StakeHolderClaimReward:           event.NewStakeHolderClaimRewardDB(tx),
			StakerDelegated:                  event.NewStakerDelegatedDB(tx),
			StakerUndelegated:                event.NewStakerUndelegatedDB(tx),
			StrategyDeposit:                  event.NewStrategyDepositDB(tx),
			StrategyWithdrawalDelayBlocksSet: event.NewStrategyWithdrawalDelayBlocksSetDB(tx),
			WithdrawalMigrated:               event.NewWithdrawalMigratedDB(tx),
			WithdrawalQueued:                 event.NewWithdrawalQueuedDB(tx),
			WithdrawalCompleted:              event.NewWithdrawalCompletedDB(tx),
			Operators:                        worker.NewOperatorsDB(tx),
			OperatorPublicKeys:               worker.NewOperatorPublicKeysDB(tx),
			StakeHolder:                      worker.NewStakeHolderDB(tx),
			TotalOperator:                    worker.NewTotalOperatorDB(tx),
		}
		return fn(txDB)
	})
}

func (db *DB) Close() error {
	sql, err := db.gorm.DB()
	if err != nil {
		return err
	}
	return sql.Close()
}

func (db *DB) ExecuteSQLMigration(migrationsFolder string) error {
	err := filepath.Walk(migrationsFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to process migration file: %s", path))
		}
		if info.IsDir() {
			return nil
		}
		fileContent, readErr := os.ReadFile(path)
		if readErr != nil {
			return errors.Wrap(readErr, fmt.Sprintf("Error reading SQL file: %s", path))
		}

		execErr := db.gorm.Exec(string(fileContent)).Error
		if execErr != nil {
			return errors.Wrap(execErr, fmt.Sprintf("Error executing SQL script: %s", path))
		}
		return nil
	})
	return err
}
