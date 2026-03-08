package announcement

import (
	"context"
	"errors"
	"fmt"
	"tournament-manager/infra/ws"
	"tournament-manager/internal/domain"
)

func (s *service) CreateAnnouncement(ctx context.Context, tournamentID int, userID int, req domain.AnnouncementCreateRequest) (*domain.Announcement, error) {
	// Verify if the user is the tournament owner
	isOwner, err := s.announcementRepo.VerifyTournamentOwner(ctx, tournamentID, userID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, errors.New("only tournament owner can create announcements")
	}
	announcement := &domain.Announcement{
		TournamentID:     tournamentID,
		AuthorID:         userID,
		Title:            req.Title,
		Content:          req.Content,
		AnnouncementType: req.AnnouncementType,
		IsPinned:         req.IsPinned,
		IsCommentable:    req.IsCommentable,
	}
	announcement, err = s.announcementRepo.CreateAnnouncement(ctx, announcement)
	if err != nil {
		return nil, err
	}
	participants, err := s.announcementRepo.GetAllParticipant(ctx, tournamentID)
	if err != nil {
		return nil, err
	}
	announcement_messsege := fmt.Sprintf("New announcement: %s", announcement.Title)
	// add notification for each participant
	notification_added := s.announcementRepo.AddAnnouncementNotification(ctx, announcement.ID, announcement_messsege, participants)
	if notification_added != nil {
		return nil, notification_added
	}
	// send websocket notification to participants
	for _, p := range participants {
		s.hub.Broadcast <- ws.Notification{
			UserID:  p.UserID,
			Message: []byte(announcement_messsege),
		}
	}
	return announcement, nil
}

func (s *service) GetAnnouncements(ctx context.Context, tournamentID int, userID int) ([]*domain.Announcement, error) {
	// Verify if the user is the tournament owner or participant
	isOwner, err := s.announcementRepo.VerifyTournamentOwner(ctx, tournamentID, userID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		participants, err := s.announcementRepo.GetAllParticipant(ctx, tournamentID)
		if err != nil {
			return nil, err
		}
		isParticipant := false
		for _, p := range participants {
			if p.UserID == userID {
				isParticipant = true
				break
			}
		}
		if !isParticipant {
			return nil, errors.New("only tournament owner and participants can view announcements")
		}
	}
	return s.announcementRepo.GetAnnouncements(ctx, tournamentID)
}

func (s *service) GetAnnouncementByID(ctx context.Context, tournamentID int, announcementID int, userID int) (*domain.Announcement, error) {
	// Verify if the user is the tournament owner or participant
	isOwner, err := s.announcementRepo.VerifyTournamentOwner(ctx, tournamentID, userID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		participants, err := s.announcementRepo.GetAllParticipant(ctx, tournamentID)
		if err != nil {
			return nil, err
		}
		isParticipant := false
		for _, p := range participants {
			if p.UserID == userID {
				isParticipant = true
				break
			}
		}
		if !isParticipant {
			return nil, errors.New("only tournament owner and participants can view announcements")
		}
	}
	return s.announcementRepo.GetAnnouncementByID(ctx, tournamentID, announcementID, userID)
}

func (s *service) UpdateAnnouncement(ctx context.Context, tournamentID int, announcementID int, userID int, req domain.AnnouncementCreateRequest) (*domain.Announcement, error) {
	// Verify if the user is the tournament owner
	isOwner, err := s.announcementRepo.VerifyTournamentOwner(ctx, tournamentID, userID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, errors.New("only tournament owner can update announcements")
	}
	announcement := &domain.Announcement{
		ID:               announcementID,
		TournamentID:     tournamentID,
		AuthorID:         userID,
		Title:            req.Title,
		Content:          req.Content,
		AnnouncementType: req.AnnouncementType,
		IsPinned:         req.IsPinned,
		IsCommentable:    req.IsCommentable,
	}
	return s.announcementRepo.UpdateAnnouncement(ctx, announcement)
}

func (s *service) DeleteAnnouncement(ctx context.Context, tournamentID int, announcementID int, userID int) error {
	// Verify if the user is the tournament owner
	isOwner, err := s.announcementRepo.VerifyTournamentOwner(ctx, tournamentID, userID)
	if err != nil {
		return err
	}
	if !isOwner {
		return errors.New("only tournament owner can delete announcements")
	}
	return s.announcementRepo.DeleteAnnouncement(ctx, tournamentID, announcementID)
}

func (s *service) GetParticipantsAnnouncementSeenStatus(ctx context.Context, tournamentID int, announcementID int, userID int) (*[]domain.Participant, error) {
	// Verify if the user is the tournament owner
	isOwner, err := s.announcementRepo.VerifyTournamentOwner(ctx, tournamentID, userID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, errors.New("only tournament owner can view participants seen status")
	}
	return s.announcementRepo.GetParticipantsAnnouncementSeenStatus(ctx, tournamentID, announcementID, userID)
}

func (s *service) ReactOnAnnouncement(ctx context.Context, tournamentID int, announcementID int, userID int, reaction string) (*domain.Announcement, error) {
	// Verify user is a participant
	participants, err := s.announcementRepo.GetAllParticipant(ctx, tournamentID)
	if err != nil {
		return nil, err
	}
	isParticipant := false
	for _, p := range participants {
		if p.UserID == userID {
			isParticipant = true
			break
		}
	}
	if !isParticipant {
		return nil, errors.New("only tournament owner and participants can react on announcements")
	}
	if reaction != "like" && reaction != "dislike" {
		return nil, errors.New("invalid reaction type")
	}

	previousReaction, err := s.announcementRepo.GetAnnouncementPrevReaction(ctx, tournamentID, announcementID, userID)
	if err != nil {
		return nil, err
	}

	switch previousReaction {
	case reaction:
		// User wants to remove their reaction
		return s.announcementRepo.RemoveAnnouncementReaction(ctx, tournamentID, announcementID, userID, reaction)
	case "":
		// User is reacting for the first time
		return s.announcementRepo.ReactOnAnnouncement(ctx, tournamentID, announcementID, userID, reaction)
	}

	// User is changing reaction
	_, err = s.announcementRepo.RemoveAnnouncementReaction(ctx, tournamentID, announcementID, userID, previousReaction)
	if err != nil {
		return nil, err
	}
	return s.announcementRepo.ReactOnAnnouncement(ctx, tournamentID, announcementID, userID, reaction)
}
