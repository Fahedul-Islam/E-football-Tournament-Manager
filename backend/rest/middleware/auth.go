package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			jwtSecret := os.Getenv("JWT_SECRET")

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
				return
			}

			token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(jwtSecret), nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			// Check expiration
			if exp, ok := claims["exp"].(float64); ok && time.Now().Unix() > int64(exp) {
				http.Error(w, "Token expired", http.StatusUnauthorized)
				return
			}

			// Role validation
			if requiredRole != "" && claims["roles"] != requiredRole {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			userId, err := normalizeUserID(claims["user_id"])
			if err != nil {
				http.Error(w, "Invalid user ID", http.StatusBadRequest)
				return
			}
			// Add claims to context
			ctx := context.WithValue(r.Context(), "user_id", userId)
			ctx = context.WithValue(ctx, "email", claims["email"])
			ctx = context.WithValue(ctx, "roles", claims["roles"])

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func normalizeUserID(id interface{}) (string, error) {
	switch v := id.(type) {
	case string:
		return v, nil
	case float64: // numbers from JWT claims are float64
		return strconv.Itoa(int(v)), nil
	default:
		return "", fmt.Errorf("unsupported user_id type: %T", v)
	}
}
