package announcement

import (
	"context"
	"errors"
	"tournament-manager/internal/domain"
)

func (s *service) CreateComment(ctx context.Context, tournamentID int, announcementID int, userID int, req domain.CommentCreateRequest) (*domain.AnnouncementComment, error) {
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
			return nil, errors.New("user is not a participant of the tournament")
		}
	}
	return s.announcementRepo.AddComment(ctx, announcementID, userID, req.ParentCommentID, &req.Content)
}

func (s *service) GetComments(ctx context.Context, tournamentID int, announcementID int, parentCommentID *int, userID int) ([]*domain.AnnouncementComment, error) {
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
			return nil, errors.New("user is not a participant of the tournament")
		}
	}

	return s.announcementRepo.GetComments(ctx, announcementID, parentCommentID)
}

func (s *service) DeleteComment(ctx context.Context, tournamentID int, userID int, commentID int) error {
	isOwner, err := s.announcementRepo.VerifyTournamentOwner(ctx, tournamentID, userID)
	if err != nil {
		return err
	}
	if !isOwner {
		participants, err := s.announcementRepo.GetAllParticipant(ctx, tournamentID)
		if err != nil {
			return err
		}
		isParticipant := false
		for _, p := range participants {
			if p.UserID == userID {
				isParticipant = true
				break
			}
		}
		if !isParticipant {
			return errors.New("user is not a participant of the tournament")
		}
	}

	return s.announcementRepo.DeleteComment(ctx, commentID)
}

func (s *service) EditComment(ctx context.Context, tournamentID int, userID int, commentID int, req domain.CommentCreateRequest) (*domain.AnnouncementComment, error) {
	isCommentOwner, err := s.announcementRepo.VerifyCommentOwner(ctx, commentID, userID)
	if err != nil {
		return nil, err
	}
	if !isCommentOwner {
		return nil, errors.New("user is not the owner of the comment")
	}
	return s.announcementRepo.EditComment(ctx, commentID, req.Content)
}

func (s *service) ReactToComment(ctx context.Context, tournamentID int, commentID int, userID int, reaction string) (*domain.AnnouncementComment, error) {
	// Validate reaction type
	if reaction != "like" && reaction != "dislike" {
		return nil, errors.New("invalid reaction type: must be 'like' or 'dislike'")
	}

	// Verify comment exists and belongs to the tournament
	exists, err := s.announcementRepo.VerifyCommentBelongsToTournament(ctx, tournamentID, commentID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("comment not found in this tournament")
	}

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
			return nil, errors.New("user is not a participant of the tournament")
		}
	}
	previousReaction, err := s.announcementRepo.GetCommentPrevReaction(ctx, tournamentID, commentID, userID)
	if err != nil {
		return nil, err
	}
	if previousReaction == reaction {
		return s.announcementRepo.RemoveReactionFromComment(ctx, commentID, userID)
	} else if previousReaction != "" {
		_, err = s.announcementRepo.RemoveReactionFromComment(ctx, commentID, userID)
		if err != nil {
			return nil, err
		}
	}
	return s.announcementRepo.ReactToComment(ctx, commentID, userID, reaction)
}
