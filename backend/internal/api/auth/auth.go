package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	config "github.com/theEricHoang/lovenote/backend/internal"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte(config.LoadConfig().JWTSecretKey)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GenerateJWT(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 24 hour expiration
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}
