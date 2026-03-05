package tournament

import (
	"context"
	"errors"
	"tournament-manager/internal/domain"
)

func (s *service) CreateAnnouncement(ctx context.Context, tournamentID int, userID int, req domain.AnnouncementCreateRequest) (*domain.Announcement, error) {
	// Verify if the user is the tournament owner
	isOwner, err := s.tournamentRepo.VerifyTournamentOwner(ctx, tournamentID, userID)
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
	return s.tournamentRepo.CreateAnnouncement(ctx, announcement)
}

func (s *service) GetAnnouncements(ctx context.Context, tournamentID int, userID int) ([]*domain.Announcement, error) {
	// Verify if the user is the tournament owner or participant
	isOwner, err := s.tournamentRepo.VerifyTournamentOwner(ctx, tournamentID, userID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		participants, err := s.tournamentRepo.GetAllParticipant(ctx, tournamentID)
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
	return s.tournamentRepo.GetAnnouncements(ctx, tournamentID)
}

func (s *service) GetAnnouncementByID(ctx context.Context, tournamentID int, announcementID int, userID int) (*domain.Announcement, error) {
	// Verify if the user is the tournament owner or participant
	isOwner, err := s.tournamentRepo.VerifyTournamentOwner(ctx, tournamentID, userID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		participants, err := s.tournamentRepo.GetAllParticipant(ctx, tournamentID)
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
	return s.tournamentRepo.GetAnnouncementByID(ctx, tournamentID, announcementID, userID)
}

func (s *service) UpdateAnnouncement(ctx context.Context, tournamentID int, announcementID int, userID int, req domain.AnnouncementCreateRequest) (*domain.Announcement, error) {
	// Verify if the user is the tournament owner
	isOwner, err := s.tournamentRepo.VerifyTournamentOwner(ctx, tournamentID, userID)
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
	return s.tournamentRepo.UpdateAnnouncement(ctx, announcement)
}

func (s *service) DeleteAnnouncement(ctx context.Context, tournamentID int, announcementID int, userID int) error {
	// Verify if the user is the tournament owner
	isOwner, err := s.tournamentRepo.VerifyTournamentOwner(ctx, tournamentID, userID)
	if err != nil {
		return err
	}
	if !isOwner {
		return errors.New("only tournament owner can delete announcements")
	}
	return s.tournamentRepo.DeleteAnnouncement(ctx, tournamentID, announcementID)
}

func (s *service) GetParticipantsAnnouncementSeenStatus(ctx context.Context, tournamentID int, announcementID int, userID int) (*[]domain.Participant, error) {
	// Verify if the user is the tournament owner
	isOwner, err := s.tournamentRepo.VerifyTournamentOwner(ctx, tournamentID, userID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, errors.New("only tournament owner can view participants seen status")
	}
	return s.tournamentRepo.GetParticipantsAnnouncementSeenStatus(ctx, tournamentID, announcementID, userID)
}