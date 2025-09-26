package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"learn/fiber/pkg/model"
	"learn/fiber/types/enum"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(jwtPayload model.JwtPayload) (string, error) {
	secret := os.Getenv("JWT_SECRET_ACCESS_TOKEN")
	if secret == "" {
		return "", errors.New("secret not found in environment variables")
	}

	return generateToken(jwtPayload, secret, 15*time.Minute)
}

func GenerateRefreshToken(jwtPayload model.JwtPayload) (string, error) {
	secret := os.Getenv("JWT_SECRET_REFRESH_TOKEN")
	if secret == "" {
		return "", errors.New("secret key not found in environment variables")
	}

	return generateToken(jwtPayload, secret, 7*24*time.Hour)
}

func generateToken(jwtPayload model.JwtPayload, secret string, expTime time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"id":   jwtPayload.Id,
		"role": jwtPayload.Role,
		"exp":  time.Now().Add(expTime).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return t, nil
}

func ValidateToken(token string) (model.JwtPayload, error) {
	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		return model.JwtPayload{}, errors.New("JWT_SECRET not found in environment variables")
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return model.JwtPayload{}, fmt.Errorf("failed to parse token: %w", err)
	}

	if !parsedToken.Valid {
		return model.JwtPayload{}, errors.New("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok || !parsedToken.Valid {
		return model.JwtPayload{}, errors.New("invalid token")
	}

	return model.JwtPayload{
		Id:   claims["id"].(string),
		Role: enum.ERole(claims["role"].(string)),
	}, nil
}
