package team

type CreateTeamRequest struct {
	Name         string `json:"name"`
	City         string `json:"city"`
	Abbreviation string `json:"abbreviation"`
}

type UpdateTeamRequest struct {
	Name         string `json:"name"`
	City         string `json:"city"`
	Abbreviation string `json:"abbreviation"`
}

type TeamResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	City         string `json:"city"`
	Abbreviation string `json:"abbreviation"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func toResponse(team Team) TeamResponse {
	return TeamResponse{
		ID:           team.ID,
		Name:         team.Name,
		City:         team.City,
		Abbreviation: team.Abbreviation,
	}
}

func toResponseList(teams []Team) []TeamResponse {
	responses := make([]TeamResponse, 0, len(teams))

	for _, t := range teams {
		responses = append(responses, toResponse(t))
	}

	return responses
}
