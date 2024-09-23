package api

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ethereum/go-ethereum/log"

	"github.com/eniac-x-labs/manta-indexer/api/common/httputil"
	"github.com/eniac-x-labs/manta-indexer/api/routes"
	"github.com/eniac-x-labs/manta-indexer/api/service"
	"github.com/eniac-x-labs/manta-indexer/config"
	"github.com/eniac-x-labs/manta-indexer/database"
)

const ethereumAddressRegex = `^0x[a-fA-F0-9]{40}$`

const (
	HealthPath = "/healthz"

	StrategyV1Path     = "/api/v1/strategies/strategy"
	StrategyListV1Path = "/api/v1/strategies/strategy/list"

	OperatorGetV1Path                 = "/api/v1/operator/get"
	OperatorListV1Path                = "/api/v1/operator/list"
	OperatorRegisterV1Path            = "/api/v1/operator/register"
	OperatorRegisterListV1Path        = "/api/v1/operator/register/list"
	OperatorNodeUrlUpdateListV1Path   = "/api/v1/operator/node/url/update/list"
	OperatorReceiveStakerDelegateList = "/api/v1/operator/receiver/staker/delegate/list"
	OperatorSharesIncreasedListV1Path = "/api/v1/operator/shares/increased/list"
	OperatorSharesDecreasedListV1Path = "/api/v1/operator/shares/decreased/list"
	OperatorAndStakeRewardListV1Path  = "/api/v1/operator/stake/reward/list"
	OperatorClaimRewardListV1Path     = "/api/v1/operator/claim/reward/list"

	StakerGetV1Path                          = "/api/v1/staker/get"
	StakerListV1Path                         = "/api/v1/staker/list"
	StakerOperatorListV1Path                 = "/api/v1/staker/operator/list"
	StakerDepositStrategyListV1Path          = "/api/v1/staker/deposit/strategy/list"
	StakerDelegatedListV1Path                = "/api/v1/staker/delegated/list"
	StakerUndelegatedListV1Path              = "/api/v1/staker/undelegated/list"
	StakeHolderClaimRewardListV1Path         = "/api/v1/staker/claim/reward/list"
	StakeHolderWithdrawalQueuedListV1Path    = "/api/v1/staker/withdrawal/queued/list"
	StakeHolderWithdrawalCompletedListV1Path = "/api/v1/staker/withdrawal/completed/list"

	FinalityVerifiedV1Path = "/api/v1/finality/verified"
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

	svc := service.New(v, a.db.OperatorRegistered, a.db.OperatorNodeUrlUpdate, a.db.Operators, a.db.StakeStrategy, a.db.StakerOperator, a.db.StrategyDeposit,
		a.db.WithdrawalQueued, a.db.WithdrawalCompleted, a.db.StakerDelegated, a.db.StakerUndelegated, a.db.StakeHolderClaimReward,
		a.db.OperatorSharesDecreased, a.db.OperatorSharesIncreased, a.db.OperatorAndStakeReward, a.db.OperatorClaimReward, a.db.Strategies, a.db.FinalityVerified)
	apiRouter := chi.NewRouter()
	h := routes.NewRoutes(apiRouter, svc)

	apiRouter.Use(middleware.Timeout(time.Second * 12))
	apiRouter.Use(middleware.Recoverer)

	apiRouter.Use(middleware.Heartbeat(HealthPath))

	/*
	* ============== Strategy ==============
	 */
	apiRouter.Get(fmt.Sprintf(StrategyV1Path), h.StrategyHandler)
	apiRouter.Get(fmt.Sprintf(StrategyListV1Path), h.StrategyListHandler)

	/*
	* ============== Operator ==============
	 */
	apiRouter.Get(fmt.Sprintf(OperatorGetV1Path), h.GetOperatorHandler)
	apiRouter.Get(fmt.Sprintf(OperatorListV1Path), h.ListOperatorHandler)
	apiRouter.Get(fmt.Sprintf(OperatorRegisterV1Path), h.RegisterOperatorHandler)
	apiRouter.Get(fmt.Sprintf(OperatorRegisterListV1Path), h.RegisterOperatorListHandler)
	apiRouter.Get(fmt.Sprintf(OperatorNodeUrlUpdateListV1Path), h.ListOperatorNodeUrlUpdateHandler)
	apiRouter.Get(fmt.Sprintf(OperatorReceiveStakerDelegateList), h.ListOperatorReceiveStakerDelegateHandler)
	apiRouter.Get(fmt.Sprintf(OperatorSharesDecreasedListV1Path), h.ListOperatorSharesDecreasedHandler)
	apiRouter.Get(fmt.Sprintf(OperatorSharesIncreasedListV1Path), h.ListOperatorSharesIncreasedHandler)
	apiRouter.Get(fmt.Sprintf(OperatorAndStakeRewardListV1Path), h.ListOperatorAndStakeRewardHandler)
	apiRouter.Get(fmt.Sprintf(OperatorClaimRewardListV1Path), h.ListOperatorClaimRewardHandler)

	/*
	* ============== Stakeholder ==============
	 */
	apiRouter.Get(fmt.Sprintf(StakerGetV1Path), h.GetStakeHolderHandler)
	apiRouter.Get(fmt.Sprintf(StakerListV1Path), h.ListStakeHolderHandler)
	apiRouter.Get(fmt.Sprintf(StakerOperatorListV1Path), h.ListStakeOperatorHandler)
	apiRouter.Get(fmt.Sprintf(StakerDepositStrategyListV1Path), h.ListStakerDepositStrategyHandler)
	apiRouter.Get(fmt.Sprintf(StakerDelegatedListV1Path), h.ListStakerDelegatedHandler)
	apiRouter.Get(fmt.Sprintf(StakerUndelegatedListV1Path), h.ListStakerUndelegatedHandler)
	apiRouter.Get(fmt.Sprintf(StakeHolderClaimRewardListV1Path), h.ListStakeHolderClaimRewardHandler)
	apiRouter.Get(fmt.Sprintf(StakeHolderWithdrawalQueuedListV1Path), h.ListStakeHolderWithdrawalQueuedHandler)
	apiRouter.Get(fmt.Sprintf(StakeHolderWithdrawalCompletedListV1Path), h.ListStakeHolderWithdrawalCompletedHandler)

	apiRouter.Get(fmt.Sprintf(FinalityVerifiedV1Path), h.GetFinalityVerifiedHandler)

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
