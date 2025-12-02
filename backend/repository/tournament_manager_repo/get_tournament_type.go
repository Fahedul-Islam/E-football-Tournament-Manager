package tournamentmanagerrepo

func (r *tournamentManagerRepo) GetTournamentType(tournament_id int) (string, error) {
	var tournament_type string
	err := r.db.QueryRow("SELECT tournament_type FROM tournaments WHERE id = $1", tournament_id).Scan(&tournament_type)
	if err != nil {
		return "", err
	}
	return tournament_type, nil
}
