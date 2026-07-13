package team

import (
	"context"
	"fmt"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) ListTeams(ctx context.Context) ([]TeamResponse, error) {
	teams, err := s.repository.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list teams: %w", err)
	}

	return toResponseList(teams), nil
}
