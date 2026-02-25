package domain

type Tournament struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	TournamentType string `json:"tournament_type"` // e.g., "single_elimination", "round_robin", "league"
	MaxPlayers     int    `json:"max_players"`
	CreatedBy      int    `json:"created_by"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
	CreatedAt      string `json:"created_at"`

	Participants []Participant `json:"participants,omitempty"`
}

type TournamentCreateRequest struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	TournamentType string `json:"tournament_type"` // e.g., "single_elimination", "round_robin", "league"
	MaxPlayers     int    `json:"max_players"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
}