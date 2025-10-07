package domain


type User struct {
	ID           int    "json:\"id\""
	Username     string "json:\"username\""
	Email        string "json:\"email\""
	PasswordHash string "json:\"-\""
	Role         string "json:\"role\""
	CreatedAt    string "json:\"created_at\""
}
