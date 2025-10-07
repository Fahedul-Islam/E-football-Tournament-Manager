package tournamentmanager

import "tournament-manager/domain"

type Service interface {
	CreateTournament(created_by int, request domain.TournamentCreateRequest) error
	GetTournamentByID(id int) (*domain.Tournament, error)
	GetAllTournaments() ([]*domain.Tournament, error)
	UpdateTournament(id int, tournament domain.TournamentCreateRequest) error
	DeleteTournament(id int) error
	ApproveParticipant(req domain.ParticipantRequest) error
	RejectParticipant(req domain.ParticipantRequest) error
	AddParticipant(participant domain.ParticipantRequest) error
	RemoveParticipant(req domain.ParticipantRequest) error
}

type TournamentManagerHandler struct {
	tournamentService Service
}

func NewTournamentManagerHandler(tournamentService Service) *TournamentManagerHandler {
	return &TournamentManagerHandler{
		tournamentService: tournamentService,
	}
}