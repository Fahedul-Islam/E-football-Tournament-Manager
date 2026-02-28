package tournament

import (
	"context"
	"errors"
	"fmt"
	"tournament-manager/internal/domain"
)

// LeagueStyleSchedule creates a league style schedule
func (s *service) LeagueStyleSchedule(ctx context.Context, tournamentID int) error {
	var approvedParticipants []*domain.Participant
	approvedParticipants, err := s.GetApprovedParticipants(ctx, tournamentID)
	if err != nil {
		return err
	}
	if len(approvedParticipants) < 2 {
		return errors.New("Not enough approved participants to create match schedules")
	}
	return s.tournamentRepo.LeagueStyleSchedule(ctx, tournamentID, approvedParticipants)
}

// CreateMatchSchedules creates match schedules for a tournament
func (s *service) CreateMatchSchedules(ctx context.Context, tournamentID int, tournament_owner_id int, groupCount int) error {
	// verify permission
	hasPermission, err := s.tournamentRepo.VerifyTournamentOwner(ctx, tournamentID, tournament_owner_id)
	if err != nil {
		return err
	}
	if !hasPermission {
		return errors.New("you don't have permission to create match schedule")
	}
	// get approved participants
	var approvedParticipants []*domain.Participant
	approvedParticipants, err = s.tournamentRepo.GetApprovedParticipants(ctx, tournamentID)
	if err != nil {
		return err
	}
	if len(approvedParticipants) < 2 {
		return errors.New("Not enough approved participants to create match schedules")
	}

	// check tournament type
	tournment_type, err := s.tournamentRepo.GetTournamentType(ctx, tournamentID)
	if err != nil {
		return err
	}
	// handling league style tournament
	if tournment_type == domain.League {
		return s.tournamentRepo.LeagueStyleSchedule(ctx, tournamentID, approvedParticipants)
	}

	if groupCount < 1 || groupCount > 8 || groupCount%2 != 0 {
		return errors.New("Group count must be between 1 and 8 and an even number")
	}
	return s.tournamentRepo.CreateMatchSchedules(ctx, tournamentID, groupCount, approvedParticipants)
}

// GenerateGroups generates groups for a tournament
func (s *service) GenerateGroups(ctx context.Context, tournamentID int, groupCount int, approvedParticipants []*domain.Participant) error {
	return s.tournamentRepo.GenerateGroups(ctx, tournamentID, groupCount, approvedParticipants)
}

// GetAllMatches returns all matches for a tournament
func (s *service) GetAllMatches(ctx context.Context, tournamentID int) ([]*domain.Match, error) {
	return s.tournamentRepo.GetAllMatches(ctx, tournamentID)
}

// UpdateScore updates the score for a match
func (s *service) UpdateScore(ctx context.Context, tournamentOwnerID int, req *domain.UpdateMatchScoreInput) (*domain.UpdateMatchScoreInput, error) {
	if (req.Round == domain.RoundOf16 || req.Round == domain.QuarterFinals || req.Round == domain.Semifinals || req.Round == domain.Final) && (req.ScoreA == req.ScoreB) {
		return nil, errors.New("Score must be different for knockout rounds")
	}

	result, err := s.tournamentRepo.UpdateScore(ctx, tournamentOwnerID, req)
	if err != nil {
		return nil, err
	}
	thisRoundDone, err := s.tournamentRepo.CheckAndAdvanceRound(ctx, req.TournamentID, req.Round)
	if err != nil {
		return nil, err
	}
	tournament_type, err := s.tournamentRepo.GetTournamentType(ctx, req.TournamentID)
	if err != nil {
		return nil, err
	}

	if (thisRoundDone && tournament_type == domain.GroupPlusKnockout) || (thisRoundDone && tournament_type == domain.Knockout) {
		switch req.Round {
		case domain.GroupStage:
			_, err = s.tournamentRepo.GenerateKnockoutStage(ctx, req.TournamentID)
			if err != nil {
				return nil, err
			}

		case domain.RoundOf16:
			_, err = s.tournamentRepo.GenerateQuarterFinals(ctx, req.TournamentID)
		case domain.QuarterFinals:
			_, err = s.tournamentRepo.GenerateSemiFinals(ctx, req.TournamentID)
		case domain.Semifinals:
			_, err = s.tournamentRepo.GenerateFinal(ctx, req.TournamentID)
		case domain.Final:
			fmt.Println("Tournament has concluded.")
		default:
			return nil, errors.New("Unknown round")
		}

		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
