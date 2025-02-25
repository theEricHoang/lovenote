package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	config "github.com/theEricHoang/lovenote/backend/internal"
	"golang.org/x/crypto/bcrypt"
)

var (
	AccessTokenExpiry  = 15 * time.Minute
	RefreshTokenExpiry = 7 * 24 * time.Hour
	SecretKey          = []byte(config.LoadConfig().JWTSecretKey)
)

type Claims struct {
	UserId uint `json:"user_id"`
	jwt.RegisteredClaims
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GenerateTokens(userId uint) (string, string, error) {
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

	return accessToken, refreshToken, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
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
