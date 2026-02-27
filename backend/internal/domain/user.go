package domain


type User struct {
	ID           int    "json:\"id\""
	Username     string "json:\"username\""
	Email        string "json:\"email\""
	PasswordHash string "json:\"-\""
	Role         string "json:\"role\""
	CreatedAt    string "json:\"created_at\""
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    string `json:"expires_in"`
	TokenType    string `json:"token_type"`
}