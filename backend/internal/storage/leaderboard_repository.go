package storage

import (
	"context"

	"prediction/internal/domain"
)

type LeaderboardOverview struct {
	TotalPredictions    int
	ResolvedPredictions int
	ActivePredictors    int
	ActiveCreators      int
	AverageWinRate      string
}

type LeaderboardPredictor struct {
	PlayerID          domain.PlayerID
	DisplayName       string
	ResolvedPositions int
	WonPositions      int
	LostPositions     int
	WinRate           string
	NetProfit         string
	TotalStaked       string
}

type LeaderboardCreator struct {
	PlayerID             domain.PlayerID
	DisplayName          string
	CreatedMarkets       int
	ResolvedMarkets      int
	TotalPositions       int
	UniqueParticipants   int
	TotalStakedOnMarkets string
}

type LeaderboardStreak struct {
	PlayerID          domain.PlayerID
	DisplayName       string
	CurrentWinStreak  int
	LongestWinStreak  int
	ResolvedPositions int
	WinRate           string
	NetProfit         string
}

type LeaderboardActivity struct {
	PlayerID          domain.PlayerID
	DisplayName       string
	TotalPositions    int
	OpenPositions     int
	ResolvedPositions int
	CreatedMarkets    int
	TotalStaked       string
}

type LeaderboardRepository interface {
	GetOverview(ctx context.Context) (LeaderboardOverview, error)
	ListTopPredictors(ctx context.Context, limit int) ([]LeaderboardPredictor, error)
	ListTopCreators(ctx context.Context, limit int) ([]LeaderboardCreator, error)
	ListBestStreaks(ctx context.Context, limit int) ([]LeaderboardStreak, error)
	ListMostActive(ctx context.Context, limit int) ([]LeaderboardActivity, error)
}
