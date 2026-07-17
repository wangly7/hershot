package team

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
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

func (s *Service) GetTeam(ctx context.Context, id string) (*TeamResponse, error) {
	team, err := s.repository.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get team by id: %w", err)
	}

	resp := toResponse(*team)

	return &resp, nil
}

func (s *Service) CreateTeam(ctx context.Context, request CreateTeamRequest) (*TeamResponse, error) {
	name, city, abbreviation, err := normalizeAndValidateTeamInput(
		request.Name,
		request.City,
		request.Abbreviation,
	)
	if err != nil {
		return nil, err
	}

	newTeam := &Team{
		ID:           uuid.NewString(),
		Name:         name,
		City:         city,
		Abbreviation: abbreviation,
	}

	if err := s.repository.Create(ctx, newTeam); err != nil {
		return nil, fmt.Errorf("create team: %w", err)
	}

	response := toResponse(*newTeam)

	return &response, nil
}

func (s *Service) UpdateTeam(
	ctx context.Context,
	id string,
	request UpdateTeamRequest,
) (*TeamResponse, error) {
	if err := validateTeamID(id); err != nil {
		return nil, err
	}

	name, city, abbreviation, err := normalizeAndValidateTeamInput(
		request.Name,
		request.City,
		request.Abbreviation,
	)
	if err != nil {
		return nil, err
	}

	updatedTeam := &Team{
		ID:           id,
		Name:         name,
		City:         city,
		Abbreviation: abbreviation,
	}

	if err := s.repository.Update(ctx, updatedTeam); err != nil {
		return nil, fmt.Errorf("update team: %w", err)
	}

	response := toResponse(*updatedTeam)

	return &response, nil
}

func (s *Service) DeleteTeam(
	ctx context.Context,
	id string,
) error {
	if err := validateTeamID(id); err != nil {
		return err
	}

	if err := s.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete team: %w", err)
	}
	return nil
}

func validateTeamID(id string) error {
	if strings.TrimSpace(id) == "" {
		return ErrInvalidID
	}
	if _, err := uuid.Parse(id); err != nil {
		return ErrInvalidID
	}
	return nil
}

func normalizeAndValidateTeamInput(
	name string,
	city string,
	abbreviation string,
) (string, string, string, error) {
	name = strings.TrimSpace(name)
	city = strings.TrimSpace(city)
	abbreviation = strings.TrimSpace(abbreviation)

	if name == "" {
		return "", "", "", ErrInvalidTeamName
	}

	if city == "" {
		return "", "", "", ErrInvalidTeamCity
	}

	if abbreviation == "" {
		return "", "", "", ErrInvalidTeamAbbreviation
	}

	if len(abbreviation) < 2 || len(abbreviation) > 4 {
		return "", "", "", ErrInvalidTeamAbbreviation
	}

	return name, city, abbreviation, nil
}
