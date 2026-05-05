package security

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/leandrogutierrez148/acomm/auth-server/internal/models"
)

func GenerateJWT(user *models.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "supersecretkey_change_in_production" // fallback
	}

	claims := jwt.MapClaims{
		"sub":       user.ID.String(),
		"email":     user.Email,
		"user_type": string(user.UserType),
		"exp":       time.Now().Add(time.Hour * 24).Unix(), // 24 hour expiration
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
