package announcement

import (
	"database/sql"
	"tournament-manager/internal/domain/repository"
)

type announcementRepo struct {
	db *sql.DB
}

func NewAnnouncementRepo(db *sql.DB) repository.AnnouncementRepository {
	return &announcementRepo{db: db}
}
