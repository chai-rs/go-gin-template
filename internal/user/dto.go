package user

import (
	"github.com/0xanonydxck/simple-bookstore/internal/model"
	"github.com/0xanonydxck/simple-bookstore/pkg/crypto"
)

// LoginRequestDTO represents the request payload for user login.
type LoginRequestDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponseDTO represents the response payload after successful login.
type LoginResponseDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// RegisterRequestDTO represents the request payload for user registration.
type RegisterRequestDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// ToUser converts RegisterRequestDTO to a User model with hashed password.
func (r *RegisterRequestDTO) ToUser() (*model.User, error) {
	hashedPassword, err := crypto.HashPassword(r.Password)
	if err != nil {
		return nil, err
	}

	return &model.User{
		Email:          r.Email,
		HashedPassword: hashedPassword,
	}, nil
}

// RegisterResponseDTO represents the response payload after successful registration.
type RegisterResponseDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// RefreshTokenRequestDTO represents the request payload for refreshing access token.
type RefreshTokenRequestDTO struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshTokenResponseDTO represents the response payload after successful token refresh.
type RefreshTokenResponseDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
