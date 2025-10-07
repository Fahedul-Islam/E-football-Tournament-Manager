package domain

type Participant struct {
	ID           int    `json:"id"`
	UserID       int    `json:"user_id"`
	TournamentID int    `json:"tournament_id"`
	TeamName     string `json:"team_name"`
	Status       string `json:"status"` // e.g., "pending", "approved", "rejected"	
	CreatedAt    string `json:"created_at"`
}

type ParticipantRequest struct {
	UserID       int    `json:"user_id"`
	TournamentID int    `json:"tournament_id"`
	TeamName     string `json:"team_name"`
}