package tournamentmanager

import "tournament-manager/domain"

type Service interface {
	CreateTournament(created_by int, request domain.TournamentCreateRequest) error
	GetTournamentByID(id int) (*domain.Tournament, error)
	GetAllTournaments() ([]*domain.Tournament, error)
	UpdateTournament(id int, tournament domain.TournamentCreateRequest) error
	DeleteTournament(tournament_owner_id int,tournament_id int) error
	ApproveParticipant(tournament_owner_id int,req domain.ParticipantRequest) error
	RejectParticipant(req domain.ParticipantRequest) error
	AddParticipant(tournament_owner_id int, participant domain.ParticipantRequest) error
	RemoveParticipant(req domain.ParticipantRequest) error
	GetAllParticipant(tournament_id int) ([]*domain.Participant, error)
}

type TournamentManagerHandler struct {
	tournamentService Service
}

func NewTournamentManagerHandler(tournamentService Service) *TournamentManagerHandler {
	return &TournamentManagerHandler{
		tournamentService: tournamentService,
	}
}
