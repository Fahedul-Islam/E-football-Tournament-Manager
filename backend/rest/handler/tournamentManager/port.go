package tournamentmanager

import "tournament-manager/domain"

type Service interface {
	CreateTournament(created_by int, request domain.TournamentCreateRequest) error
	GetTournamentByID(id int) (*domain.Tournament, error)
	GetAllTournaments(tournament_owner_id int) ([]*domain.Tournament, error)
	UpdateTournament(tournament_owner_id int, tournament_id int, tournament domain.TournamentCreateRequest) error
	DeleteTournament(tournament_owner_id int, tournament_id int) error
	ApproveParticipant(tournament_owner_id int, req domain.ParticipantRequest) error
	RejectParticipant(tournament_owner_id int, req domain.ParticipantRequest) error
	AddParticipant(tournament_owner_id int, participant domain.ParticipantRequest) error
	RemoveParticipant(tournament_owner_id int, req domain.ParticipantRequest) error
	GetAllParticipant(tournament_id int) ([]*domain.Participant, error)
	GetApprovedParticipants(tournament_id int) ([]*domain.Participant, error)
	CreateMatchSchedules(tournament_id int, group_count int, approvedParticipants []*domain.Participant) error
	GenerateGroups(tournament_id int, groupCount int, approvedParticipants []*domain.Participant) error
	GetAllMatches(tournament_id int) ([]*domain.Match, error)
	UpdateScore(tournament_owner_id int,req *domain.UpadateMatchScoreInput) (*domain.UpadateMatchScoreInput,error)
	CheckAndAdvanceRound(tournament_id int, round string) (bool, error)	
	GetGroupCount(tournament_id int) (int, error)
	GenerateKnockoutStage(tournament_id int) (bool, error)
	GenerateQuarterFinals(tournament_id int) (bool, error)
	GenerateSemiFinals(tournament_id int) (bool, error)
	GenerateFinal(tournament_id int) (bool, error)
	GetLeaderboard(tournament_id int) (map[int][]domain.PlayerStat, error)
	GetTournamentType(tournament_id int) (string, error)
	LeagueStyleSchedule(tournament_id int, approvedParticipants []*domain.Participant) error
	VerifyTournamentOwner(tournament_id int, user_id int) (bool, error)
}

type TournamentManagerHandler struct {
	tournamentService Service
}

func NewTournamentManagerHandler(tournamentService Service) *TournamentManagerHandler {
	return &TournamentManagerHandler{
		tournamentService: tournamentService,
	}
}
