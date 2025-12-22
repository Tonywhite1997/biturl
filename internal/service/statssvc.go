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

	exists, url := s.PGRepo.ShortCodeExists(statsAccessKey)
	if !exists {
		return nil, errors.New("invalid url")
	}

	totalClicks, err := s.ClickhouseRepo.GetTotalClicks(ctx, url.ShortCode)
	totalCountries, err := s.ClickhouseRepo.GetCountryClicks(ctx, url.ShortCode)
	totalDevices, err := s.ClickhouseRepo.GetDeviceClicks(ctx, url.ShortCode)
	uniqueVisitors, err := s.ClickhouseRepo.GetUniqueVisitors(ctx, url.ShortCode)
	dailyClicks, err := s.ClickhouseRepo.GetDailyClicks(ctx, url.ShortCode)
	totalBrowsers, err := s.ClickhouseRepo.GetBrowserClicks(ctx, url.ShortCode)

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
		OriginalURL:    url.OriginalURL,
		ExpiresAt:      url.ExpiresAt,
	}

	return &results, nil
}
