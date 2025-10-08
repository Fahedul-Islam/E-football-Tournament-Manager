package participant

import "tournament-manager/domain"

type Service interface {
	RequestToJoinTournament(domain.ParticipantRequest) error
}

type ParticipantHandler struct {
	service Service
}

func NewParticipantHandler(service Service) *ParticipantHandler {
	return &ParticipantHandler{service: service}
}
