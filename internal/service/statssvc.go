package service

import (
	"biturl/internal/dto"
	"biturl/internal/repository"
	"context"
	"errors"
	"fmt"
)

type StatsSVC struct {
	ClickhouseRepo repository.ClkHouseRepo
	PGRepo         repository.PGrepo
}

func (s StatsSVC) GetStats(ctx context.Context, statsAccessKey string) (*dto.StatsResponse, error) {
	if len(statsAccessKey) == 0 {
		return nil, errors.New("invalid shortcode")
	}

	exists, shortCode, originalURL := s.PGRepo.ShortCodeExists(statsAccessKey)
	if !exists {
		return nil, errors.New("invalid url")
	}

	totalClicks, err := s.ClickhouseRepo.GetTotalClicks(ctx, *shortCode)
	totalCountries, err := s.ClickhouseRepo.GetCountryClicks(ctx, *shortCode)
	totalDevices, err := s.ClickhouseRepo.GetDeviceClicks(ctx, *shortCode)
	uniqueVisitors, err := s.ClickhouseRepo.GetUniqueVisitors(ctx, *shortCode)
	dailyClicks, err := s.ClickhouseRepo.GetDailyClicks(ctx, *shortCode)
	totalBrowsers, err := s.ClickhouseRepo.GetBrowserClicks(ctx, *shortCode)

	if err != nil {
		fmt.Println("could not get stats", err)
		return nil, errors.New("could not get stats")
	}

	results := dto.StatsResponse{
		TotalClicks:    totalClicks,
		UniqueVisitors: uniqueVisitors,
		Browsers:       totalBrowsers,
		Devices:        totalDevices,
		Countries:      totalCountries,
		DailyStats:     dailyClicks,
		OriginalURL:    *originalURL,
	}

	return &results, nil
}
