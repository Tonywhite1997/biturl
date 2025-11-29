package service

import (
	"biturl/internal/repository"
	"context"
	"errors"
	"fmt"
)

type StatsSVC struct {
	ClickhouseRepo repository.ClkHouseRepo
	PGRepo         repository.PGrepo
}

func (s StatsSVC) GetStatsByShortCode(ctx context.Context, statsAccessKey string) ([]repository.Stats, error) {
	if len(statsAccessKey) == 0 {
		return nil, errors.New("invalid shortcode")
	}

	exists, shortCode := s.PGRepo.ShortCodeExists(statsAccessKey)
	if !exists {
		return nil, errors.New("invalid url")
	}

	stats, err := s.ClickhouseRepo.GetBySHortID(ctx, *shortCode)

	if err != nil {
		fmt.Println("could not get stats", err)
		return nil, errors.New("could not get stats")
	}

	return stats, nil
}
