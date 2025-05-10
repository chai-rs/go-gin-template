package user

import (
	"github.com/0xanonydxck/simple-bookstore/internal/model"
	"github.com/0xanonydxck/simple-bookstore/pkg/crypto"
)

type LoginRequestDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponseDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RegisterRequestDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

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

type RegisterResponseDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenRequestDTO struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenResponseDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
