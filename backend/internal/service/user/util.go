package user

import (
	"strconv"
	"time"
	"tournament-manager/internal/domain"

	"github.com/golang-jwt/jwt/v4"
)

func (s *service) generateToken(user *domain.User) (string, string, error) {
	now := time.Now()
	accessClaims := jwt.MapClaims{
		"user_id": strconv.Itoa(int(user.ID)),
		"exp":     now.Add(s.cfg.JWT.TokenExpiry).Unix(),
		"iat":     now.Unix(),
		"email":   user.Email,
		"roles":   user.Role,
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refreshClaims := jwt.MapClaims{
		"user_id": strconv.Itoa(int(user.ID)),
		"exp":     now.Add(s.cfg.JWT.RefreshExpiry).Unix(),
		"iat":     now.Unix(),
		"email":   user.Email,
		"roles":   user.Role,
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	accessTokenString, err := accessToken.SignedString(s.cfg.JWT.Secret)
	if err != nil {
		return "", "", err
	}
	refreshTokenString, err := refreshToken.SignedString(s.cfg.JWT.Secret)
	if err != nil {
		return "", "", err
	}
	return accessTokenString, refreshTokenString, nil
}
