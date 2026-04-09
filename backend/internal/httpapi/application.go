package httpapi

import (
	"net/http"

	"prediction/internal/auth"
	"prediction/internal/balance"
	"prediction/internal/config"
	"prediction/internal/leaderboard"
	"prediction/internal/market"
	"prediction/internal/player"
	"prediction/internal/position"
	"prediction/internal/realtime"
	"prediction/internal/storage"
)

type Dependencies struct {
	Config    config.Config
	DB        *storage.DB
	TxManager *storage.TxManager
}

type Controllers struct {
	Auth        auth.Controller
	Player      player.Controller
	Balance     balance.Controller
	Leaderboard leaderboard.Controller
	Market      market.Controller
	Position    position.Controller
	Realtime    realtime.Controller
}

type Application struct {
	Dependencies Dependencies
	Controllers  Controllers
	Router       *Router
}

func NewApplication(cfg config.Config, db *storage.DB) *Application {
	return &Application{
		Dependencies: Dependencies{
			Config:    cfg,
			DB:        db,
			TxManager: storage.NewTxManager(db.Pool()),
		},
		Router: NewRouter(),
	}
}

func (a *Application) WithControllers(controllers Controllers) {
	a.Controllers = controllers
}

func (a *Application) RegisterRoute(method, pattern string, handler http.Handler) {
	a.Router.Handle(Route{
		Method:  method,
		Pattern: pattern,
		Handler: handler,
	})
}
