package utils

import (
	"errors"
	"strings"
	"time"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/config"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(config.Load().JWT_SECRET)
var signingMethod = jwt.SigningMethodHS256
var tokenExpirationTime = 5 * time.Hour             // el tiempo de expiración de los access tokens
var refreshTokenExpirationTime = 7 * 24 * time.Hour // el tiempo de expiración de los refresh tokens (7 días)
var nowFunc = time.Now

func CreateToken(userId string, role enums.UserRole) (string, error) {
	now := nowFunc()
	token := jwt.NewWithClaims(signingMethod,
		jwt.MapClaims{
			"sub":  userId,
			"iat":  now.Unix(),
			"exp":  now.Add(tokenExpirationTime).Unix(),
			"role": role,
			"type": "access",
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// CreateRefreshToken genera un refresh token con mayor tiempo de expiración
func CreateRefreshToken(userId string, role enums.UserRole) (string, error) {
	now := nowFunc()
	token := jwt.NewWithClaims(signingMethod,
		jwt.MapClaims{
			"sub":  userId,
			"iat":  now.Unix(),
			"exp":  now.Add(refreshTokenExpirationTime).Unix(),
			"role": role,
			"type": "refresh",
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// RefreshAccessToken verifica el refresh token y genera nuevos tokens
func RefreshAccessToken(refreshToken string) (accessToken string, newRefreshToken string, err error) {
	trimmed := strings.TrimSpace(refreshToken)
	if trimmed == "" {
		return "", "", errors.New("refresh token inválido o expirado")
	}

	claims, err := VerifyToken(trimmed)
	if err != nil {
		return "", "", errors.New("refresh token inválido o expirado")
	}

	// Verificar que sea un refresh token
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return "", "", errors.New("token no es un refresh token válido")
	}

	userId, ok := claims["sub"].(string)
	if !ok || strings.TrimSpace(userId) == "" {
		return "", "", errors.New("no se pudo extraer el userId del token")
	}

	roleStr, ok := claims["role"].(string)
	if !ok || strings.TrimSpace(roleStr) == "" {
		return "", "", errors.New("no se pudo extraer el role del token")
	}
	role := enums.UserRole(roleStr)

	// Generar nuevo access token
	accessToken, err = CreateToken(userId, role)
	if err != nil {
		return "", "", err
	}

	// Generar nuevo refresh token (rotación de tokens)
	newRefreshToken, err = CreateRefreshToken(userId, role)
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	trimmed := strings.TrimSpace(tokenString)
	if trimmed == "" {
		return jwt.MapClaims{}, errors.New("token vacío")
	}

	token, err := jwt.Parse(
		trimmed,
		func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		},
		jwt.WithValidMethods([]string{signingMethod.Alg()}),
		jwt.WithExpirationRequired(),
		jwt.WithTimeFunc(nowFunc),
	)

	if err != nil {
		return jwt.MapClaims{}, err
	}

	if !token.Valid {
		return jwt.MapClaims{}, errors.New("token inválido")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("no se pudieron extraer los claims")
	}

	return claims, nil
}
