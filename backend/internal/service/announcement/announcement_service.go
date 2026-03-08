package announcement

import (
	"tournament-manager/infra/ws"
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
	hub              *ws.Hub
}

func NewAnnouncementService(announcementRepo AnnouncementRepo, hub *ws.Hub) Service {
	return &service{
		announcementRepo: announcementRepo,
		hub:              hub,
	}
}
