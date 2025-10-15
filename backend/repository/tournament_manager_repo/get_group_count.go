package tournamentmanagerrepo

func (r *tournamentManagerRepo) GetGroupCount(tournament_id int) (int, error) {
	var group_count int
	query := `SELECT COUNT(ID) FROM groups WHERE tournament_id = $1`
	err := r.db.QueryRow(query, tournament_id).Scan(&group_count)
	if err != nil {
		return 0, err
	}
	return group_count, nil
}