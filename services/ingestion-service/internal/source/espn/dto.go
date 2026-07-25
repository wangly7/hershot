package espn

// ===========================
// Scoreboard DTO
// ===========================
type ScoreboardResponse struct {
	Events []ScoreboardEventDTO `json:"events"`
}

type ScoreboardEventDTO struct {
	ID           string              `json:"id"`
	Date         string              `json:"date"`
	Competitions []CompetitionDTO    `json:"competitions"`
	Status       ScoreboardStatusDTO `json:"status"`
}

type CompetitionDTO struct {
	ID          string          `json:"id"`
	Competitors []CompetitorDTO `json:"competitors"`
}

type CompetitorDTO struct {
	ID       string  `json:"id"`
	HomeAway string  `json:"homeAway"`
	Team     TeamDTO `json:"team"`
	Score    string  `json:"score"`
}

type TeamDTO struct {
	ID string `json:"id"`
}

type ScoreboardStatusDTO struct {
	Type StatusTypeDTO `json:"type"`
}

type StatusTypeDTO struct {
	State     string `json:"state"`
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

// ===========================
// Play-by-Play DTO
// ===========================
type PlayResponse struct {
	Count     int `json:"count"`
	PageIndex int `json:"pageIndex"`
	PageSize  int `json:"pageSize"`
	PageCount int `json:"pageCount"`

	Items []PlayDTO `json:"items"`
}

type PlayDTO struct {
	ID             string         `json:"id"`
	SequenceNumber FlexiableInt64 `json:"sequenceNumber"`

	Text                 string `json:"text"`
	ShortText            string `json:"shortText"`
	AlternativeText      string `json:"alternativeText"`
	ShortAlternativeText string `json:"shortAlternativeText"`

	Clock  ClockDTO  `json:"clock"`
	Period PeriodDTO `json:"period"`

	Type PlayTypeDTO `json:"type"`

	Valid       bool   `json:"valid"`
	ScoringPlay bool   `json:"scoringPlay"`
	ScoreValue  int    `json:"scoreValue"`
	Modified    string `json:"modified"`

	HomeScore int `json:"homeScore"`
	AwayScore int `json:"awayScore"`

	Team RefDTO `json:"team"`

	Participants []ParticipantDTO `json:"participants"`
}

type ClockDTO struct {
	Value        float32 `json:"value"`
	DisplayValue string  `json:"displayValue"`
}

type PeriodDTO struct {
	Number       int    `json:"number"`
	DisplayValue string `json:"displayValue"`
}

type PlayTypeDTO struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type RefDTO struct {
	Ref string `json:"$ref"`
}

type ParticipantDTO struct {
	Athlete    RefDTO `json:"athlete"`
	Position   RefDTO `json:"position"`
	Statistics RefDTO `json:"statistics"`
	Order      int    `json:"order"`
}
