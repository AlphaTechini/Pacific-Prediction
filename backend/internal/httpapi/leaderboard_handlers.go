package httpapi

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"prediction/internal/leaderboard"
)

const defaultLeaderboardLimit = 8

type leaderboardResponse struct {
	GeneratedAt   string                         `json:"generated_at"`
	Overview      leaderboardOverviewResponse    `json:"overview"`
	TopPredictors []predictorLeaderboardResponse `json:"top_predictors"`
	TopCreators   []creatorLeaderboardResponse   `json:"top_creators"`
	BestStreaks   []streakLeaderboardResponse    `json:"best_streaks"`
	MostActive    []activityLeaderboardResponse  `json:"most_active"`
}

type leaderboardOverviewResponse struct {
	TotalPredictions    int    `json:"total_predictions"`
	ResolvedPredictions int    `json:"resolved_predictions"`
	ActivePredictors    int    `json:"active_predictors"`
	ActiveCreators      int    `json:"active_creators"`
	AverageWinRate      string `json:"average_win_rate"`
}

type predictorLeaderboardResponse struct {
	Rank              int    `json:"rank"`
	PlayerID          string `json:"player_id"`
	DisplayName       string `json:"display_name"`
	ResolvedPositions int    `json:"resolved_positions"`
	WonPositions      int    `json:"won_positions"`
	LostPositions     int    `json:"lost_positions"`
	WinRate           string `json:"win_rate"`
	NetProfit         string `json:"net_profit"`
	TotalStaked       string `json:"total_staked"`
}

type creatorLeaderboardResponse struct {
	Rank                 int    `json:"rank"`
	PlayerID             string `json:"player_id"`
	DisplayName          string `json:"display_name"`
	CreatedMarkets       int    `json:"created_markets"`
	ResolvedMarkets      int    `json:"resolved_markets"`
	TotalPositions       int    `json:"total_positions"`
	UniqueParticipants   int    `json:"unique_participants"`
	TotalStakedOnMarkets string `json:"total_staked_on_markets"`
}

type streakLeaderboardResponse struct {
	Rank              int    `json:"rank"`
	PlayerID          string `json:"player_id"`
	DisplayName       string `json:"display_name"`
	CurrentWinStreak  int    `json:"current_win_streak"`
	LongestWinStreak  int    `json:"longest_win_streak"`
	ResolvedPositions int    `json:"resolved_positions"`
	WinRate           string `json:"win_rate"`
	NetProfit         string `json:"net_profit"`
}

type activityLeaderboardResponse struct {
	Rank              int    `json:"rank"`
	PlayerID          string `json:"player_id"`
	DisplayName       string `json:"display_name"`
	TotalPositions    int    `json:"total_positions"`
	OpenPositions     int    `json:"open_positions"`
	ResolvedPositions int    `json:"resolved_positions"`
	CreatedMarkets    int    `json:"created_markets"`
	TotalStaked       string `json:"total_staked"`
}

func NewGetLeaderboardHandler(controller leaderboard.Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limit, err := parseLeaderboardLimit(r)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
			return
		}

		snapshot, err := controller.GetSnapshot(r.Context(), limit)
		if err != nil {
			writeError(w, r, err)
			return
		}

		w.Header().Set("Cache-Control", "public, max-age=15, stale-while-revalidate=60")
		writeJSON(w, http.StatusOK, toLeaderboardResponse(snapshot))
	})
}

func parseLeaderboardLimit(r *http.Request) (int, error) {
	value := r.URL.Query().Get("limit")
	if value == "" {
		return defaultLeaderboardLimit, nil
	}

	limit, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid_limit")
	}

	if limit < 1 || limit > 25 {
		return 0, fmt.Errorf("invalid_limit")
	}

	return limit, nil
}

func toLeaderboardResponse(snapshot leaderboard.Snapshot) leaderboardResponse {
	return leaderboardResponse{
		GeneratedAt:   snapshot.GeneratedAt.Format(time.RFC3339),
		Overview:      toLeaderboardOverviewResponse(snapshot.Overview),
		TopPredictors: toPredictorLeaderboardResponses(snapshot.TopPredictors),
		TopCreators:   toCreatorLeaderboardResponses(snapshot.TopCreators),
		BestStreaks:   toStreakLeaderboardResponses(snapshot.BestStreaks),
		MostActive:    toActivityLeaderboardResponses(snapshot.MostActive),
	}
}

func toLeaderboardOverviewResponse(overview leaderboard.Overview) leaderboardOverviewResponse {
	return leaderboardOverviewResponse{
		TotalPredictions:    overview.TotalPredictions,
		ResolvedPredictions: overview.ResolvedPredictions,
		ActivePredictors:    overview.ActivePredictors,
		ActiveCreators:      overview.ActiveCreators,
		AverageWinRate:      overview.AverageWinRate,
	}
}

func toPredictorLeaderboardResponses(records []leaderboard.PredictorEntry) []predictorLeaderboardResponse {
	items := make([]predictorLeaderboardResponse, 0, len(records))
	for _, record := range records {
		items = append(items, predictorLeaderboardResponse{
			Rank:              record.Rank,
			PlayerID:          record.PlayerID,
			DisplayName:       record.DisplayName,
			ResolvedPositions: record.ResolvedPositions,
			WonPositions:      record.WonPositions,
			LostPositions:     record.LostPositions,
			WinRate:           record.WinRate,
			NetProfit:         record.NetProfit,
			TotalStaked:       record.TotalStaked,
		})
	}

	return items
}

func toCreatorLeaderboardResponses(records []leaderboard.CreatorEntry) []creatorLeaderboardResponse {
	items := make([]creatorLeaderboardResponse, 0, len(records))
	for _, record := range records {
		items = append(items, creatorLeaderboardResponse{
			Rank:                 record.Rank,
			PlayerID:             record.PlayerID,
			DisplayName:          record.DisplayName,
			CreatedMarkets:       record.CreatedMarkets,
			ResolvedMarkets:      record.ResolvedMarkets,
			TotalPositions:       record.TotalPositions,
			UniqueParticipants:   record.UniqueParticipants,
			TotalStakedOnMarkets: record.TotalStakedOnMarkets,
		})
	}

	return items
}

func toStreakLeaderboardResponses(records []leaderboard.StreakEntry) []streakLeaderboardResponse {
	items := make([]streakLeaderboardResponse, 0, len(records))
	for _, record := range records {
		items = append(items, streakLeaderboardResponse{
			Rank:              record.Rank,
			PlayerID:          record.PlayerID,
			DisplayName:       record.DisplayName,
			CurrentWinStreak:  record.CurrentWinStreak,
			LongestWinStreak:  record.LongestWinStreak,
			ResolvedPositions: record.ResolvedPositions,
			WinRate:           record.WinRate,
			NetProfit:         record.NetProfit,
		})
	}

	return items
}

func toActivityLeaderboardResponses(records []leaderboard.ActivityEntry) []activityLeaderboardResponse {
	items := make([]activityLeaderboardResponse, 0, len(records))
	for _, record := range records {
		items = append(items, activityLeaderboardResponse{
			Rank:              record.Rank,
			PlayerID:          record.PlayerID,
			DisplayName:       record.DisplayName,
			TotalPositions:    record.TotalPositions,
			OpenPositions:     record.OpenPositions,
			ResolvedPositions: record.ResolvedPositions,
			CreatedMarkets:    record.CreatedMarkets,
			TotalStaked:       record.TotalStaked,
		})
	}

	return items
}
