package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	config "github.com/theEricHoang/lovenote/backend/internal"
)

type contextKey string

const UserIDKey contextKey = "userId"

var secretKey = []byte(config.LoadConfig().JWTSecretKey)

func AuthenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			http.Error(w, "Invalid Authorization format", http.StatusUnauthorized)
			return
		}
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return secretKey, nil // look up secret key
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		userIdFloat, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
			return
		}

		expFloat, ok := claims["exp"].(float64)
		if !ok {
			http.Error(w, "Token expiration time missing", http.StatusUnauthorized)
			return
		}

		expTime := time.Unix(int64(expFloat), 0)
		if time.Now().After(expTime) {
			http.Error(w, "Token expired", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, uint(userIdFloat))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
