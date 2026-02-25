package tournamentmanagerrepo

func (r *tournamentManagerRepo) VerifyTournamentOwner(tournament_id int, user_id int) (bool, error) {
	var owner_id int
	err := r.db.QueryRow("SELECT tournament_owner_id FROM tournaments WHERE id = $1", tournament_id).Scan(&owner_id)
	if err != nil {
		return false, err
	}
	return owner_id == user_id, nil
}