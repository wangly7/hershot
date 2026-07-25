package espn

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	defaultSiteBaseURL = "https://site.api.espn.com"
	defaultCoreBaseURL = "https://sports.core.api.espn.com"

	defaultRequestTimeout = 10 * time.Second
	defaultPlaysLimit     = 1000
)

type Client interface {
	GetScoreboard(ctx context.Context, date time.Time) (ScoreboardResponse, error)
	GetPlays(ctx context.Context, eventID string, competitionID string) (PlayResponse, error)
}

type client struct {
	httpClient *http.Client

	siteBaseURL string
	coreBaseURL string
	playsLimits int
}

type ClientConfig struct {
	SiteBaseURL string
	CoreBaseURL string
	Timeout     time.Duration
	PlaysLimits int
}

func NewClient(config ClientConfig) Client {
	siteBaseURL := strings.TrimRight(
		strings.TrimSpace(config.SiteBaseURL),
		"/",
	)
	if siteBaseURL == "" {
		siteBaseURL = defaultSiteBaseURL
	}

	coreBaseURL := strings.TrimRight(
		strings.TrimSpace(config.CoreBaseURL),
		"/",
	)

	if coreBaseURL == "" {
		coreBaseURL = defaultCoreBaseURL
	}

	timeout := config.Timeout
	if timeout <= 0 {
		timeout = defaultRequestTimeout
	}

	limits := config.PlaysLimits
	if limits <= 0 {
		limits = defaultPlaysLimit
	}

	return &client{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		siteBaseURL: siteBaseURL,
		coreBaseURL: coreBaseURL,
		playsLimits: limits,
	}
}

func (c *client) GetScoreboard(ctx context.Context, date time.Time) (ScoreboardResponse, error) {
	endpoint := fmt.Sprintf(
		"%s/apis/site/v2/sports/basketball/wnba/scoreboard?dates=%s",
		c.siteBaseURL,
		date.Format("20060102"),
	)

	var response ScoreboardResponse

	if err := c.doJSON(ctx, endpoint, &response); err != nil {
		return ScoreboardResponse{}, fmt.Errorf("get ESPN scoreboard: %w", err)
	}

	return response, nil
}

func (c *client) GetPlays(
	ctx context.Context,
	eventID string,
	competitionID string,
) (PlayResponse, error) {
	eventID = strings.TrimSpace(eventID)
	competitionID = strings.TrimSpace(competitionID)

	if eventID == "" {
		return PlayResponse{}, fmt.Errorf("event ID is required")
	}

	if competitionID == "" {
		return PlayResponse{}, fmt.Errorf(
			"competition ID is required",
		)
	}
	endpoint, err := c.buildPlayURL(eventID, competitionID)
	if err != nil {
		return PlayResponse{}, err
	}

	var response PlayResponse
	if err := c.doJSON(ctx, endpoint, &response); err != nil {
		return PlayResponse{}, fmt.Errorf("get ESPN plays for event %q: %w", eventID, err)
	}

	return response, nil
}

func (c *client) buildPlayURL(
	eventID string,
	competitionID string,
) (string, error) {
	base, err := url.Parse(c.coreBaseURL)
	if err != nil {
		return "", fmt.Errorf(
			"parse ESPN core base URL %q: %w",
			c.coreBaseURL,
			err,
		)
	}

	base.Path = fmt.Sprintf(
		"/v2/sports/basketball/leagues/wnba/events/%s/competitions/%s/plays",
		url.PathEscape(eventID),
		url.PathEscape(competitionID),
	)

	query := base.Query()
	query.Set("lang", "en")
	query.Set("region", "us")
	query.Set("limit", strconv.Itoa(c.playsLimits))

	base.RawQuery = query.Encode()

	return base.String(), nil
}

func (c *client) doJSON(
	ctx context.Context,
	endpoint string,
	target any,
) error {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		endpoint,
		nil,
	)
	if err != nil {
		return fmt.Errorf(
			"create request for %q: %w",
			endpoint,
			err,
		)
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", "hershot-ingestion-service/1.0")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf(
			"send request to %q: %w",
			endpoint,
			err,
		)
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK ||
		response.StatusCode >= http.StatusMultipleChoices {
		return c.newHTTPStatusError(endpoint, response)
	}

	decoder := json.NewDecoder(response.Body)

	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("decode response from %q: %w", endpoint, err)
	}
	return nil
}

func (c *client) newHTTPStatusError(endpoint string, response *http.Response) error {
	const maxErrorBodySize = 4 * 1024

	body, readErr := io.ReadAll(
		io.LimitReader(response.Body, maxErrorBodySize),
	)
	if readErr != nil {
		return fmt.Errorf(
			"ESPN request to %q returned status %s; read error body: %w",
			endpoint,
			response.Status,
			readErr,
		)
	}

	message := strings.TrimSpace(string(body))
	return fmt.Errorf(
		"ESPN request to %q returned status %s: %s",
		endpoint,
		response.Status,
		message,
	)
}
