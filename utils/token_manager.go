package utils

import (
	"errors"
	"fmt"
	"time"

	"learn/fiber/config"
	"learn/fiber/pkg/enum"
	"learn/fiber/pkg/model"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(jwtPayload model.JwtPayload) (string, error) {
	secret := config.JWT_SECRET_ACCESS_TOKEN.GetValue()

	if secret == "" {
		return "", errors.New("secret not found in environment variables")
	}

	return generateToken(jwtPayload, secret, 15*time.Minute)
}

func GenerateRefreshToken(jwtPayload model.JwtPayload) (string, error) {
	secret := config.JWT_SECRET_REFRESH_TOKEN.GetValue()

	if secret == "" {
		return "", errors.New("secret key not found in environment variables")
	}

	return generateToken(jwtPayload, secret, 7*24*time.Hour)
}

func GenerateNewAccessToken(token string) (string, error) {
	secret := config.JWT_SECRET_REFRESH_TOKEN.GetValue()

	payload, err := ValidateToken(token, secret)

	if err != nil {
		return "", err
	}

	return GenerateAccessToken(payload)
}

func ValidateToken(token string, secret string) (model.JwtPayload, error) {
	if secret == "" {
		return model.JwtPayload{}, errors.New("secret Key not found in environment variables")
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

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok || !parsedToken.Valid {
		return model.JwtPayload{}, errors.New("invalid token")
	}

	return model.JwtPayload{
		Id:   claims["id"].(string),
		Role: enum.ERole(claims["role"].(string)),
	}, nil
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
