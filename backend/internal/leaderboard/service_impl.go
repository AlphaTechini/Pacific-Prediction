package leaderboard

import (
	"context"
	"fmt"
	"time"

	"prediction/internal/domain"
	"prediction/internal/storage"

	"golang.org/x/sync/errgroup"
)

type service struct {
	repository storage.LeaderboardRepository
}

func NewService(repository storage.LeaderboardRepository) Service {
	return &service{repository: repository}
}

func (s *service) GetSnapshot(ctx context.Context, limit int) (Snapshot, error) {
	if limit < 1 || limit > 25 {
		return Snapshot{}, domain.NewValidationError("limit", "limit must be between 1 and 25", limit)
	}

	var (
		overview      storage.LeaderboardOverview
		topPredictors []storage.LeaderboardPredictor
		topCreators   []storage.LeaderboardCreator
		bestStreaks   []storage.LeaderboardStreak
		mostActive    []storage.LeaderboardActivity
	)

	group, groupCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		record, err := s.repository.GetOverview(groupCtx)
		if err != nil {
			return fmt.Errorf("get leaderboard overview: %w", err)
		}

		overview = record
		return nil
	})

	group.Go(func() error {
		records, err := s.repository.ListTopPredictors(groupCtx, limit)
		if err != nil {
			return fmt.Errorf("list top predictors: %w", err)
		}

		topPredictors = records
		return nil
	})

	group.Go(func() error {
		records, err := s.repository.ListTopCreators(groupCtx, limit)
		if err != nil {
			return fmt.Errorf("list top creators: %w", err)
		}

		topCreators = records
		return nil
	})

	group.Go(func() error {
		records, err := s.repository.ListBestStreaks(groupCtx, limit)
		if err != nil {
			return fmt.Errorf("list best streaks: %w", err)
		}

		bestStreaks = records
		return nil
	})

	group.Go(func() error {
		records, err := s.repository.ListMostActive(groupCtx, limit)
		if err != nil {
			return fmt.Errorf("list most active players: %w", err)
		}

		mostActive = records
		return nil
	})

	if err := group.Wait(); err != nil {
		return Snapshot{}, err
	}

	return Snapshot{
		GeneratedAt:   time.Now().UTC(),
		Overview:      toOverview(overview),
		TopPredictors: toPredictorEntries(topPredictors),
		TopCreators:   toCreatorEntries(topCreators),
		BestStreaks:   toStreakEntries(bestStreaks),
		MostActive:    toActivityEntries(mostActive),
	}, nil
}

func toOverview(record storage.LeaderboardOverview) Overview {
	return Overview{
		TotalPredictions:    record.TotalPredictions,
		ResolvedPredictions: record.ResolvedPredictions,
		ActivePredictors:    record.ActivePredictors,
		ActiveCreators:      record.ActiveCreators,
		AverageWinRate:      record.AverageWinRate,
	}
}

func toPredictorEntries(records []storage.LeaderboardPredictor) []PredictorEntry {
	items := make([]PredictorEntry, 0, len(records))
	for index, record := range records {
		items = append(items, PredictorEntry{
			Rank:              index + 1,
			PlayerID:          string(record.PlayerID),
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

func toCreatorEntries(records []storage.LeaderboardCreator) []CreatorEntry {
	items := make([]CreatorEntry, 0, len(records))
	for index, record := range records {
		items = append(items, CreatorEntry{
			Rank:                 index + 1,
			PlayerID:             string(record.PlayerID),
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

func toStreakEntries(records []storage.LeaderboardStreak) []StreakEntry {
	items := make([]StreakEntry, 0, len(records))
	for index, record := range records {
		items = append(items, StreakEntry{
			Rank:              index + 1,
			PlayerID:          string(record.PlayerID),
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

func toActivityEntries(records []storage.LeaderboardActivity) []ActivityEntry {
	items := make([]ActivityEntry, 0, len(records))
	for index, record := range records {
		items = append(items, ActivityEntry{
			Rank:              index + 1,
			PlayerID:          string(record.PlayerID),
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
