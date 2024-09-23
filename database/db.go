package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum/log"

	"github.com/eniac-x-labs/manta-indexer/config"
	"github.com/eniac-x-labs/manta-indexer/database/common"
	"github.com/eniac-x-labs/manta-indexer/database/event"
	"github.com/eniac-x-labs/manta-indexer/database/event/finality"
	"github.com/eniac-x-labs/manta-indexer/database/event/operator"
	"github.com/eniac-x-labs/manta-indexer/database/event/staker"
	"github.com/eniac-x-labs/manta-indexer/database/event/strategies"
	_ "github.com/eniac-x-labs/manta-indexer/database/utils/serializers"
	"github.com/eniac-x-labs/manta-indexer/database/worker"
	"github.com/eniac-x-labs/manta-indexer/synchronizer/retry"
)

type DB struct {
	gorm                             *gorm.DB
	Blocks                           common.BlocksDB
	ContractEvent                    event.ContractEventDB
	EventBlocks                      event.EventBlocksDB
	MinWithdrawalDelayBlocksSet      strategies.MinWithdrawalDelayBlocksSetDB
	OperatorAndStakeReward           operator.OperatorAndStakeRewardDB
	OperatorClaimReward              operator.OperatorClaimRewardDB
	OperatorModified                 operator.OperatorModifiedDB
	OperatorNodeUrlUpdate            operator.OperatorNodeUrlUpdateDB
	OperatorRegistered               operator.OperatorRegisteredDB
	OperatorSharesDecreased          operator.OperatorSharesDecreasedDB
	OperatorSharesIncreased          operator.OperatorSharesIncreasedDB
	StakeHolderClaimReward           staker.StakeHolderClaimRewardDB
	StakerDelegated                  staker.StakerDelegatedDB
	StakerUndelegated                staker.StakerUndelegatedDB
	Strategies                       strategies.StrategiesDB
	StrategyDeposit                  staker.StrategyDepositDB
	StrategyWithdrawalDelayBlocksSet strategies.StrategyWithdrawalDelayBlocksSetDB
	WithdrawalMigrated               staker.WithdrawalMigratedDB
	WithdrawalQueued                 staker.WithdrawalQueuedDB
	WithdrawalCompleted              staker.WithdrawalCompletedDB
	Operators                        worker.OperatorsDB
	OperatorPublicKeys               worker.OperatorPublicKeysDB
	StakerOperator                   worker.StakerOperatorDB
	StakeStrategy                    worker.StakeStrategyDB
	TotalOperator                    worker.TotalOperatorDB
	FinalityVerified                 finality.FinalityVerifiedDB
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
	log.Info("database NewDB dsn ", "info", dsn)
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
		MinWithdrawalDelayBlocksSet:      strategies.NewMinWithdrawalDelayBlocksSetDB(gorm),
		OperatorAndStakeReward:           operator.NewOperatorAndStakeRewardDB(gorm),
		OperatorClaimReward:              operator.NewOperatorClaimRewardDB(gorm),
		OperatorModified:                 operator.NewOperatorModifiedDB(gorm),
		OperatorNodeUrlUpdate:            operator.NewOperatorNodeUrlUpdateDB(gorm),
		OperatorRegistered:               operator.NewOperatorRegisteredDB(gorm),
		OperatorSharesDecreased:          operator.NewOperatorSharesDecreasedDB(gorm),
		OperatorSharesIncreased:          operator.NewOperatorSharesIncreasedDB(gorm),
		StakeHolderClaimReward:           staker.NewStakeHolderClaimRewardDB(gorm),
		StakerDelegated:                  staker.NewStakerDelegatedDB(gorm),
		StakerUndelegated:                staker.NewStakerUndelegatedDB(gorm),
		Strategies:                       strategies.NewStrategiesDB(gorm),
		StrategyDeposit:                  staker.NewStrategyDepositDB(gorm),
		StrategyWithdrawalDelayBlocksSet: strategies.NewStrategyWithdrawalDelayBlocksSetDB(gorm),
		WithdrawalMigrated:               staker.NewWithdrawalMigratedDB(gorm),
		WithdrawalQueued:                 staker.NewWithdrawalQueuedDB(gorm),
		WithdrawalCompleted:              staker.NewWithdrawalCompletedDB(gorm),
		Operators:                        worker.NewOperatorsDB(gorm),
		OperatorPublicKeys:               worker.NewOperatorPublicKeysDB(gorm),
		StakerOperator:                   worker.NewStakerOperatorDB(gorm),
		StakeStrategy:                    worker.NewStakeStrategyDB(gorm),
		TotalOperator:                    worker.NewTotalOperatorDB(gorm),
		FinalityVerified:                 finality.NewFinalityVerifiedDB(gorm),
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
			MinWithdrawalDelayBlocksSet:      strategies.NewMinWithdrawalDelayBlocksSetDB(tx),
			OperatorAndStakeReward:           operator.NewOperatorAndStakeRewardDB(tx),
			OperatorClaimReward:              operator.NewOperatorClaimRewardDB(tx),
			OperatorModified:                 operator.NewOperatorModifiedDB(tx),
			OperatorNodeUrlUpdate:            operator.NewOperatorNodeUrlUpdateDB(tx),
			OperatorRegistered:               operator.NewOperatorRegisteredDB(tx),
			OperatorSharesDecreased:          operator.NewOperatorSharesDecreasedDB(tx),
			OperatorSharesIncreased:          operator.NewOperatorSharesIncreasedDB(tx),
			StakeHolderClaimReward:           staker.NewStakeHolderClaimRewardDB(tx),
			StakerDelegated:                  staker.NewStakerDelegatedDB(tx),
			StakerUndelegated:                staker.NewStakerUndelegatedDB(tx),
			Strategies:                       strategies.NewStrategiesDB(tx),
			StrategyDeposit:                  staker.NewStrategyDepositDB(tx),
			StrategyWithdrawalDelayBlocksSet: strategies.NewStrategyWithdrawalDelayBlocksSetDB(tx),
			WithdrawalMigrated:               staker.NewWithdrawalMigratedDB(tx),
			WithdrawalQueued:                 staker.NewWithdrawalQueuedDB(tx),
			WithdrawalCompleted:              staker.NewWithdrawalCompletedDB(tx),
			Operators:                        worker.NewOperatorsDB(tx),
			OperatorPublicKeys:               worker.NewOperatorPublicKeysDB(tx),
			StakerOperator:                   worker.NewStakerOperatorDB(tx),
			StakeStrategy:                    worker.NewStakeStrategyDB(tx),
			TotalOperator:                    worker.NewTotalOperatorDB(tx),
			FinalityVerified:                 finality.NewFinalityVerifiedDB(tx),
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
