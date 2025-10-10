package domain

type Match struct {
	ID             int    `json:"id"`
	TournamentID   int    `json:"tournament_id"`
	GroupID        int    `json:"group_id"`
	Round          string `json:"round"` // e.g., "Group Stage", "Quarterfinal", etc.
	ParticipantAID int    `json:"participant_a_id"`
	ParticipantBID int    `json:"participant_b_id"`
	ScoreA         *int   `json:"score_a,omitempty"`    // Nullable, will be set when the match is concluded
	ScoreB         *int   `json:"score_b,omitempty"`    // Nullable, will be set when the match is concluded
	MatchDate      *string `json:"match_date,omitempty"` // Nullable, can be set when scheduling the match
	Status         string  `json:"status"`               // e.g., "Scheduled", "Completed", "Pending"
	CreatedAt      *string  `json:"created_at"`
	WinnerID       *int    `json:"winner_id,omitempty"`  // Nullable, will be set when the match is concluded
}
