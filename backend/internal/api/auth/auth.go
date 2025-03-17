package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	config "github.com/theEricHoang/lovenote/backend/internal"
	db "github.com/theEricHoang/lovenote/backend/internal/pkg"
	"golang.org/x/crypto/bcrypt"
)

var (
	AccessTokenExpiry  = 15 * time.Minute
	RefreshTokenExpiry = 7 * 24 * time.Hour
	SecretKey          = []byte(config.LoadConfig().JWTSecretKey)
)

type AuthService struct {
	DB *db.Database
}

type Claims struct {
	UserId uint `json:"user_id"`
	jwt.RegisteredClaims
}

func NewAuthService(db *db.Database) *AuthService {
	return &AuthService{DB: db}
}

func insertRefreshToken(ctx context.Context, db db.Database, userID uint, tokenStr string, exp time.Time) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO UPDATE
		SET token = EXCLUDED.token, expires_at = EXCLUDED.expires_at
	`
	_, err := db.Pool.Exec(ctx, query, userID, tokenStr, exp)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) GetRefreshToken(ctx context.Context, userID uint) (string, error) {
	query := `
		SELECT token FROM refresh_tokens WHERE user_id = $1
	`

	var token string
	err := s.DB.Pool.QueryRow(ctx, query, userID).Scan(&token)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func (s *AuthService) CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *AuthService) GenerateTokens(ctx context.Context, userId uint) (string, string, error) {
	accessClaims := Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenExpiry)),
		},
	}
	accessToken, err := generateJWT(accessClaims)
	if err != nil {
		return "", "", err
	}

	refreshClaims := Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenExpiry)),
		},
	}
	refreshToken, err := generateJWT(refreshClaims)
	if err != nil {
		return "", "", err
	}

	err = insertRefreshToken(ctx, *s.DB, userId, refreshToken, time.Now().Add(RefreshTokenExpiry))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return SecretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func generateJWT(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}
