package domain

type PlayerStat struct {
	ParticipantID   int `json:"participant_id"`
	MatchesPlayed   int `json:"matches_played"`
	Wins            int `json:"wins"`
	Draws           int `json:"draws"`
	Losses          int `json:"losses"`
	GoalsScored     int `json:"goals_scored"`
	GoalsConceded   int `json:"goals_conceded"`
	GoalDifference  int `json:"goal_difference"`
	Points          int `json:"points"`
}