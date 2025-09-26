package model

import "learn/fiber/pkg/enum"

type JwtPayload struct {
	Id   string     `json:"id"`
	Role enum.ERole `json:"role"`
}

type JwtResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"accessToken"`
}
