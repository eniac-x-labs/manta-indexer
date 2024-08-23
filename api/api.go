package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/eniac-x-labs/manta-indexer/api/common/httputil"
	"github.com/eniac-x-labs/manta-indexer/api/routes"
	"github.com/eniac-x-labs/manta-indexer/api/service"
	"github.com/eniac-x-labs/manta-indexer/config"
	"github.com/eniac-x-labs/manta-indexer/database"
)

const ethereumAddressRegex = `^0x[a-fA-F0-9]{40}$`

const (
	HealthPath                 = "/healthz"
	RegisterOperatorV1Path     = "/api/v1/operator/register"
	RegisterOperatorListV1Path = "/api/v1/operator/register/list"

	OperatorNodeUrlUpdateGetV1Path  = "/api/v1/operatorNodeUrlUpdate/get"
	OperatorNodeUrlUpdateListV1Path = "/api/v1/operatorNodeUrlUpdate/list"

	OperatorGetV1Path  = "/api/v1/operator/get"
	OperatorListV1Path = "/api/v1/operator/list"

	StakerGetV1Path  = "/api/v1/staker/get"
	StakerListV1Path = "/api/v1/staker/list"

	StrategyDepositGetV1Path  = "/api/v1/strategyDeposit/get"
	StrategyDepositListV1Path = "/api/v1/strategyDeposit/list"

	WithdrawalQueuedListV1Path = "/api/v1/withdrawalQueued/list"
	WithdrawalQueuedGetV1Path  = "/api/v1/withdrawalQueued/get"

	WithdrawalCompletedGetV1Path  = "/api/v1/withdrawalCompleted/get"
	WithdrawalCompletedListV1Path = "/api/v1/withdrawalCompleted/list"

	StakerDelegatedGetV1Path  = "/api/v1/stakerDelegated/get"
	StakerDelegatedListV1Path = "/api/v1/stakerDelegated/list"

	StakerUndelegatedGetV1Path  = "/api/v1/stakerUndelegated/get"
	StakerUndelegatedListV1Path = "/api/v1/stakerUndelegated/list"

	StakeHolderClaimRewardGetV1Path  = "/api/v1/stakeHolderClaimReward/get"
	StakeHolderClaimRewardListV1Path = "/api/v1/stakeHolderClaimReward/list"

	OperatorSharesDecreasedGetV1Path  = "/api/v1/operatorSharesDecreased/get"
	OperatorSharesDecreasedListV1Path = "/api/v1/operatorSharesDecreased/list"

	OperatorSharesIncreasedGetV1Path  = "/api/v1/operatorSharesIncreased/get"
	OperatorSharesIncreasedListV1Path = "/api/v1/operatorSharesIncreased/list"

	OperatorAndStakeRewardGetV1Path  = "/api/v1/operatorAndStakeReward/get"
	OperatorAndStakeRewardListV1Path = "/api/v1/operatorAndStakeReward/list"

	OperatorClaimRewardGetV1Path  = "/api/v1/operatorClaimReward/get"
	OperatorClaimRewardListV1Path = "/api/v1/operatorClaimReward/list"
)

type APIConfig struct {
	HTTPServer    config.ServerConfig
	MetricsServer config.ServerConfig
}

type API struct {
	router    *chi.Mux
	apiServer *httputil.HTTPServer
	db        *database.DB
	stopped   atomic.Bool
}

func NewApi(ctx context.Context, cfg *config.Config) (*API, error) {
	out := &API{}
	if err := out.initFromConfig(ctx, cfg); err != nil {
		return nil, errors.Join(err, out.Stop(ctx))
	}
	return out, nil
}

func (a *API) initFromConfig(ctx context.Context, cfg *config.Config) error {
	if err := a.initDB(ctx, cfg); err != nil {
		return fmt.Errorf("failed to init DB: %w", err)
	}
	a.initRouter(cfg.HTTPServer, cfg)
	if err := a.startServer(cfg.HTTPServer); err != nil {
		return fmt.Errorf("failed to start API server: %w", err)
	}
	return nil
}

func (a *API) initRouter(conf config.ServerConfig, cfg *config.Config) {
	v := new(service.Validator)

	svc := service.New(v, a.db.OperatorRegistered, a.db.OperatorNodeUrlUpdate, a.db.Operators, a.db.StakeHolder, a.db.StrategyDeposit,
		a.db.WithdrawalQueued, a.db.WithdrawalCompleted, a.db.StakerDelegated, a.db.StakerUndelegated, a.db.StakeHolderClaimReward,
		a.db.OperatorSharesDecreased, a.db.OperatorSharesIncreased, a.db.OperatorAndStakeReward, a.db.OperatorClaimReward)
	apiRouter := chi.NewRouter()
	h := routes.NewRoutes(apiRouter, svc)

	apiRouter.Use(middleware.Timeout(time.Second * 12))
	apiRouter.Use(middleware.Recoverer)

	apiRouter.Use(middleware.Heartbeat(HealthPath))

	apiRouter.Get(fmt.Sprintf(RegisterOperatorV1Path), h.RegisterOperatorHandler)
	apiRouter.Get(fmt.Sprintf(RegisterOperatorListV1Path), h.RegisterOperatorListHandler)
	apiRouter.Get(fmt.Sprintf(OperatorNodeUrlUpdateGetV1Path), h.GetOperatorNodeUrlUpdate)
	apiRouter.Get(fmt.Sprintf(OperatorNodeUrlUpdateListV1Path), h.ListOperatorNodeUrlUpdate)
	apiRouter.Get(fmt.Sprintf(StrategyDepositGetV1Path), h.GetStrategyDeposit)
	apiRouter.Get(fmt.Sprintf(StrategyDepositListV1Path), h.ListStrategyDepositHandler)
	apiRouter.Get(fmt.Sprintf(StakerGetV1Path), h.GetStakeHolder)
	apiRouter.Get(fmt.Sprintf(StakerListV1Path), h.ListStakeHolderHandler)
	apiRouter.Get(fmt.Sprintf(OperatorGetV1Path), h.GetOperator)
	apiRouter.Get(fmt.Sprintf(OperatorListV1Path), h.ListOperatorHandler)

	apiRouter.Get(fmt.Sprintf(WithdrawalQueuedGetV1Path), h.GetWithdrawalQueued)
	apiRouter.Get(fmt.Sprintf(WithdrawalQueuedListV1Path), h.ListWithdrawalQueuedHandler)

	apiRouter.Get(fmt.Sprintf(WithdrawalCompletedGetV1Path), h.GetWithdrawalCompleted)
	apiRouter.Get(fmt.Sprintf(WithdrawalCompletedListV1Path), h.ListWithdrawalCompletedHandler)

	apiRouter.Get(fmt.Sprintf(StakerDelegatedGetV1Path), h.GetStakerDelegated)
	apiRouter.Get(fmt.Sprintf(StakerDelegatedListV1Path), h.ListStakerDelegatedHandler)

	apiRouter.Get(fmt.Sprintf(StakerUndelegatedGetV1Path), h.GetStakerUndelegated)
	apiRouter.Get(fmt.Sprintf(StakerUndelegatedListV1Path), h.ListStakerUndelegatedHandler)

	apiRouter.Get(fmt.Sprintf(StakeHolderClaimRewardGetV1Path), h.GetStakeHolderClaimReward)
	apiRouter.Get(fmt.Sprintf(StakeHolderClaimRewardListV1Path), h.ListStakeHolderClaimRewardHandler)

	apiRouter.Get(fmt.Sprintf(OperatorSharesDecreasedGetV1Path), h.GetOperatorSharesDecreased)
	apiRouter.Get(fmt.Sprintf(OperatorSharesDecreasedListV1Path), h.ListOperatorSharesDecreasedHandler)

	apiRouter.Get(fmt.Sprintf(OperatorSharesIncreasedGetV1Path), h.GetOperatorSharesIncreased)
	apiRouter.Get(fmt.Sprintf(OperatorSharesIncreasedListV1Path), h.ListOperatorSharesIncreasedHandler)

	apiRouter.Get(fmt.Sprintf(OperatorAndStakeRewardGetV1Path), h.GetOperatorAndStakeReward)
	apiRouter.Get(fmt.Sprintf(OperatorAndStakeRewardListV1Path), h.ListOperatorAndStakeRewardHandler)

	apiRouter.Get(fmt.Sprintf(OperatorClaimRewardGetV1Path), h.GetOperatorClaimReward)
	apiRouter.Get(fmt.Sprintf(OperatorClaimRewardListV1Path), h.ListOperatorClaimRewardHandler)

	a.router = apiRouter
}

func (a *API) initDB(ctx context.Context, cfg *config.Config) error {
	var initDb *database.DB
	var err error
	if !cfg.SlaveDbEnable {
		initDb, err = database.NewDB(ctx, cfg.MasterDB)
		if err != nil {
			log.Error("failed to connect to master database", "err", err)
			return err
		}
	} else {
		initDb, err = database.NewDB(ctx, cfg.SlaveDB)
		if err != nil {
			log.Error("failed to connect to slave database", "err", err)
			return err
		}
	}
	a.db = initDb
	return nil
}

func (a *API) Start(ctx context.Context) error {
	return nil
}

func (a *API) Stop(ctx context.Context) error {
	var result error
	if a.apiServer != nil {
		if err := a.apiServer.Stop(ctx); err != nil {
			result = errors.Join(result, fmt.Errorf("failed to stop API server: %w", err))
		}
	}
	if a.db != nil {
		if err := a.db.Close(); err != nil {
			result = errors.Join(result, fmt.Errorf("failed to close DB: %w", err))
		}
	}
	a.stopped.Store(true)
	log.Info("API service shutdown complete")
	return result
}

func (a *API) startServer(serverConfig config.ServerConfig) error {
	log.Debug("API server listening...", "port", serverConfig.Port)
	addr := net.JoinHostPort(serverConfig.Host, strconv.Itoa(serverConfig.Port))
	srv, err := httputil.StartHTTPServer(addr, a.router)
	if err != nil {
		return fmt.Errorf("failed to start API server: %w", err)
	}
	log.Info("API server started", "addr", srv.Addr().String())
	a.apiServer = srv
	return nil
}

func (a *API) Stopped() bool {
	return a.stopped.Load()
}
