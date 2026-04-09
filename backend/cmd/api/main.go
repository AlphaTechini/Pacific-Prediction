package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"prediction/internal/auth"
	"prediction/internal/balance"
	"prediction/internal/config"
	"prediction/internal/httpapi"
	"prediction/internal/market"
	"prediction/internal/pacifica"
	"prediction/internal/player"
	"prediction/internal/position"
	"prediction/internal/realtime"
	"prediction/internal/settlement"
	"prediction/internal/storage"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := storage.NewDB(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer db.Close()

	migrator := storage.NewMigrator(db.Pool(), cfg.MigrationsDir)
	if err := migrator.Up(ctx); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	app := httpapi.NewApplication(cfg, db)
	playerRepository := storage.NewPlayerPostgresRepository(db.Pool())
	sessionRepository := storage.NewSessionPostgresRepository(db.Pool())
	balanceRepository := storage.NewBalancePostgresRepository(db.Pool())
	marketRepository := storage.NewMarketPostgresRepository(db.Pool())
	positionRepository := storage.NewPositionPostgresRepository(db.Pool())
	realtimeHub := realtime.NewHub()
	pacificaHTTPClient := &http.Client{
		Timeout: cfg.Pacifica.MarketInfoHTTPTimeout,
	}
	pacificaRESTClient := pacifica.NewHTTPRESTClient(
		cfg.Pacifica.RestBaseURL,
		pacificaHTTPClient,
		cfg.Pacifica.MarketInfoCacheTTL,
	)

	authService := auth.NewService(auth.ServiceDeps{
		Config:            cfg,
		TxManager:         app.Dependencies.TxManager,
		SessionRepository: sessionRepository,
	})
	playerService := player.NewService(playerRepository)
	balanceService := balance.NewService(balanceRepository)
	marketValidator := market.NewValidationService(pacificaRESTClient)
	marketService := market.NewServiceWithDeps(market.ServiceDeps{
		MarketRepository:      marketRepository,
		CreateContextProvider: pacificaRESTClient,
		Publisher:             realtimeHub,
		Validator:             marketValidator,
		TxManager:             app.Dependencies.TxManager,
	})
	marketController := market.NewController(marketService)
	balanceController := balance.NewController(balanceService)
	positionValidator := position.NewValidationService(position.ValidationDeps{
		MarketController:  marketController,
		BalanceController: balanceController,
	})
	positionService := position.NewService(position.ServiceDeps{
		PositionRepository: positionRepository,
		TxManager:          app.Dependencies.TxManager,
		Validator:          positionValidator,
		Publisher:          realtimeHub,
	})
	settlementService := settlement.NewService(settlement.ServiceDeps{
		MarketRepository:   marketRepository,
		PacificaClient:     pacificaRESTClient,
		Publisher:          realtimeHub,
		TxManager:          app.Dependencies.TxManager,
		PriceRetryInterval: cfg.Settlement.PriceRetryInterval,
	})
	settlementWorker := settlement.NewWorker(settlement.WorkerDeps{
		Logger:             log.Default(),
		Service:            settlementService,
		ScanInterval:       cfg.Settlement.ScanInterval,
		ScanBatchSize:      cfg.Settlement.ScanBatchSize,
		PriceLookahead:     cfg.Settlement.PriceLookahead,
		PriceRetryInterval: cfg.Settlement.PriceRetryInterval,
	})

	authController := auth.NewController(authService)
	playerController := player.NewController(playerService)
	positionController := position.NewController(positionService)
	realtimeController := realtime.NewController(realtimeHub)
	cookieManager := auth.NewCookieManager(cfg.Auth)
	requireSession := httpapi.NewRequireSessionMiddleware(authController, cookieManager)

	app.WithControllers(httpapi.Controllers{
		Auth:     authController,
		Player:   playerController,
		Balance:  balanceController,
		Market:   marketController,
		Position: positionController,
		Realtime: realtimeController,
	})

	app.RegisterRoute(http.MethodPost, "/api/v1/players/guest", httpapi.NewCreateGuestSessionHandler(authController, playerController, cookieManager))
	app.RegisterRoute(http.MethodGet, "/api/v1/players/me", requireSession(httpapi.NewGetMeHandler(playerController)))
	app.RegisterRoute(http.MethodGet, "/api/v1/players/me/balance", requireSession(httpapi.NewGetBalanceHandler(balanceController)))
	app.RegisterRoute(http.MethodGet, "/api/v1/players/me/positions", requireSession(httpapi.NewListPlayerPositionsHandler(positionController)))
	app.RegisterRoute(http.MethodGet, "/api/v1/stream", httpapi.NewStreamHandler(realtimeController, cfg.Realtime.HeartbeatInterval))
	app.RegisterRoute(http.MethodPost, "/api/v1/markets", requireSession(httpapi.NewCreateMarketHandler(marketController)))
	app.RegisterRoute(http.MethodGet, "/api/v1/markets", httpapi.NewListMarketsHandler(marketController))
	app.RegisterRoute(http.MethodGet, "/api/v1/markets/context", httpapi.NewGetMarketCreateContextHandler(marketController))
	app.RegisterRoute(http.MethodGet, "/api/v1/markets/{market_id}", httpapi.NewGetMarketDetailHandler(marketController))
	app.RegisterRoute(http.MethodPost, "/api/v1/markets/{market_id}/positions", requireSession(httpapi.NewCreatePositionHandler(positionController)))

	server := &http.Server{
		Addr:              cfg.AppAddr,
		Handler:           app.Router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("api listening on %s", cfg.AppAddr)
		errCh <- server.ListenAndServe()
	}()
	go func() {
		if err := settlementWorker.Run(ctx); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("shutdown server: %v", err)
		}
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("serve http: %v", err)
		}
	}
}
