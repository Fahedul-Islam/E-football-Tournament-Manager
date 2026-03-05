package participantrepo

import (
	"database/sql"
	"tournament-manager/internal/domain/repository"
)

type participantRepo struct {
	db *sql.DB
}

func NewParticipantRepo(db *sql.DB) repository.ParticipantRepository {
	return &participantRepo{db: db}
}
