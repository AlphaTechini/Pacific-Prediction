package leaderboard

import "time"

type Snapshot struct {
	GeneratedAt   time.Time
	Overview      Overview
	TopPredictors []PredictorEntry
	TopCreators   []CreatorEntry
	BestStreaks   []StreakEntry
	MostActive    []ActivityEntry
}

type Overview struct {
	TotalPredictions    int
	ResolvedPredictions int
	ActivePredictors    int
	ActiveCreators      int
	AverageWinRate      string
}

type PredictorEntry struct {
	Rank              int
	PlayerID          string
	DisplayName       string
	ResolvedPositions int
	WonPositions      int
	LostPositions     int
	WinRate           string
	NetProfit         string
	TotalStaked       string
}

type CreatorEntry struct {
	Rank                 int
	PlayerID             string
	DisplayName          string
	CreatedMarkets       int
	ResolvedMarkets      int
	TotalPositions       int
	UniqueParticipants   int
	TotalStakedOnMarkets string
}

type StreakEntry struct {
	Rank              int
	PlayerID          string
	DisplayName       string
	CurrentWinStreak  int
	LongestWinStreak  int
	ResolvedPositions int
	WinRate           string
	NetProfit         string
}

type ActivityEntry struct {
	Rank              int
	PlayerID          string
	DisplayName       string
	TotalPositions    int
	OpenPositions     int
	ResolvedPositions int
	CreatedMarkets    int
	TotalStaked       string
}
