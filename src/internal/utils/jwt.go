package utils

import (
	"errors"
	"time"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/config"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(config.Load().JWT_SECRET)
var signingMethod = jwt.SigningMethodHS256
var tokenExpirationTime = 5 * time.Hour // el tiempo de expiración de los tokens

func CreateToken(userId string, role enums.UserRole) (string, error) {
	token := jwt.NewWithClaims(signingMethod,
		jwt.MapClaims{
			"sub":  userId,
			"iat":  time.Now().Unix(),
			"exp":  time.Now().Add(tokenExpirationTime).Unix(),
			"role": role,
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return jwt.MapClaims{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("no se pudieron extraer los claims")
	}

	return claims, nil
}
