package participant

import (
	"context"
	"errors"
	"tournament-manager/internal/domain"
)

func (s *service) ReactOnAnnouncement(ctx context.Context, tournamentID int, announcementID int, userID int, reaction string) (*domain.Announcement, error) {
	participants, err := s.participantRepo.GetAllParticipant(ctx, tournamentID)
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

	previousReaction, err := s.participantRepo.GetAnnouncementPrevReaction(ctx, tournamentID, announcementID, userID)

	if err != nil {
		return nil, err
	}
	switch previousReaction {
	case reaction:
		// User wants to remove their reaction- need to implement
		return s.participantRepo.RemoveAnnouncementReaction(ctx, tournamentID, announcementID, userID, reaction)
	case "":
		// User is reacting for the first time
		return s.participantRepo.ReactOnAnnouncement(ctx, tournamentID, announcementID, userID, reaction)
	}
	_, err = s.participantRepo.RemoveAnnouncementReaction(ctx, tournamentID, announcementID, userID, previousReaction)
	if err != nil {
		return nil, err
	}
	return s.participantRepo.ReactOnAnnouncement(ctx, tournamentID, announcementID, userID, reaction)
}
