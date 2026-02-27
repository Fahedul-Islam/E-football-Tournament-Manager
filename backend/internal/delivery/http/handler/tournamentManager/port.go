package tournamentmanager

import (
	"context"
	"tournament-manager/internal/domain"
)

type Service interface {
	CreateTournament(ctx context.Context, created_by int, request domain.TournamentCreateRequest) error
	GetTournamentByID(ctx context.Context, id int) (*domain.Tournament, error)
	GetAllTournaments(ctx context.Context, tournament_owner_id int) ([]*domain.Tournament, error)
	UpdateTournament(ctx context.Context, tournament_owner_id int, tournament_id int, tournament domain.TournamentCreateRequest) error
	DeleteTournament(ctx context.Context, tournament_owner_id int, tournament_id int) error
	ApproveParticipant(ctx context.Context, tournament_owner_id int, req domain.ParticipantRequest) error
	RejectParticipant(ctx context.Context, tournament_owner_id int, req domain.ParticipantRequest) error
	AddParticipant(ctx context.Context, tournament_owner_id int, participant domain.ParticipantRequest) error
	RemoveParticipant(ctx context.Context, tournament_owner_id int, req domain.ParticipantRequest) error
	GetAllParticipant(ctx context.Context, tournament_id int) ([]*domain.Participant, error)
	GetApprovedParticipants(ctx context.Context, tournament_id int) ([]*domain.Participant, error)
	CreateMatchSchedules(ctx context.Context, tournament_id int, tournament_owner_id int, groupCount int) error
	GenerateGroups(ctx context.Context, tournament_id int, groupCount int, approvedParticipants []*domain.Participant) error
	GetAllMatches(ctx context.Context, tournament_id int) ([]*domain.Match, error)
	UpdateScore(ctx context.Context,tournament_owner_id int,req *domain.UpdateMatchScoreInput) (*domain.UpdateMatchScoreInput,error)
	CheckAndAdvanceRound(ctx context.Context,tournament_id int, round string) (bool, error)	
	GetGroupCount(ctx context.Context, tournament_id int) (int, error)
	GetLeaderboard(ctx context.Context, tournament_id int) (map[int][]domain.PlayerStat, error)
	GetTournamentType(ctx context.Context, tournament_id int) (string, error)
	LeagueStyleSchedule(ctx context.Context, tournament_id int) error
	VerifyTournamentOwner(ctx context.Context, tournament_id int, user_id int) (bool, error)
}

type TournamentManagerHandler struct {
	tournamentService Service
}

func NewTournamentManagerHandler(tournamentService Service) *TournamentManagerHandler {
	return &TournamentManagerHandler{
		tournamentService: tournamentService,
	}
}
