package announcement

import (
	announcementhandler "tournament-manager/internal/delivery/http/handler/announcement"
	"tournament-manager/internal/domain/repository"
)

type Service interface {
	announcementhandler.Service
}

type AnnouncementRepo interface {
	repository.AnnouncementRepository
}

type service struct {
	announcementRepo AnnouncementRepo
}

func NewAnnouncementService(announcementRepo AnnouncementRepo) Service {
	return &service{
		announcementRepo: announcementRepo,
	}
}
