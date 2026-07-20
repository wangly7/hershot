package leagueclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var ErrTeamNotFound = errors.New("team not found")

type Team struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	City         string `json:"city"`
	Abbreviation string `json:"abbreviation"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Client) ListTeams(ctx context.Context) ([]Team, error) {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.baseURL+"/teams",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("create list teams request: %w", err)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("request teams from league service: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("league service returned status %d", response.StatusCode)
	}

	var teams []Team

	if err := json.NewDecoder(response.Body).Decode(&teams); err != nil {
		return nil, fmt.Errorf("decode teams response: %w", err)
	}

	return teams, nil
}

func (c *Client) GetTeam(ctx context.Context, id string) (*Team, error) {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.baseURL+"/teams/"+id,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("create list team request: %w", err)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("request team from league service: %w", err)
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, ErrTeamNotFound
	default:
		return nil, fmt.Errorf(
			"league service returned status %d",
			response.StatusCode,
		)
	}

	var team Team
	if err := json.NewDecoder(response.Body).Decode(&team); err != nil {
		return nil, fmt.Errorf("decode team response: %w", err)
	}

	return &team, nil
}
