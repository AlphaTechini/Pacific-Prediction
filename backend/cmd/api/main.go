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
	"prediction/internal/player"
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

	authService := auth.NewService(auth.ServiceDeps{
		Config:            cfg,
		TxManager:         app.Dependencies.TxManager,
		SessionRepository: sessionRepository,
	})
	playerService := player.NewService(playerRepository)
	balanceService := balance.NewService(balanceRepository)

	authController := auth.NewController(authService)
	playerController := player.NewController(playerService)
	balanceController := balance.NewController(balanceService)
	cookieManager := auth.NewCookieManager(cfg.Auth)
	requireSession := httpapi.NewRequireSessionMiddleware(authController, cookieManager)

	app.WithControllers(httpapi.Controllers{
		Auth:    authController,
		Player:  playerController,
		Balance: balanceController,
	})

	app.RegisterRoute(http.MethodPost, "/api/v1/players/guest", httpapi.NewCreateGuestSessionHandler(authController, playerController, cookieManager))
	app.RegisterRoute(http.MethodGet, "/api/v1/players/me", requireSession(httpapi.NewGetMeHandler(playerController)))
	app.RegisterRoute(http.MethodGet, "/api/v1/players/me/balance", requireSession(httpapi.NewGetBalanceHandler(balanceController)))

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
