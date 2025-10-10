package participant

import "tournament-manager/domain"

type Service interface {
	RequestToJoinTournament(domain.ParticipantRequest) error
	IsApprovedParticipant(tournament_id int, user_id int) (bool, error)
	GetGroupDistribution(tournament_id int) ([]*domain.Group, error)
	SeeMatchSchedule(tournament_id int) ([]*domain.Match, error)
}

type ParticipantHandler struct {
	service Service
}

func NewParticipantHandler(service Service) *ParticipantHandler {
	return &ParticipantHandler{service: service}
}
