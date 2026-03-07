package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
	"github.com/ericolvr/sec-back-v2/internal/core/utils"
)

type AuthService struct {
	userRepo domain.UserRepository
}

func NewAuthService(userRepo domain.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) Login(ctx context.Context, mobile, password string) (*domain.User, error) {
	if mobile == "" {
		return nil, errors.New("mobile is required")
	}
	if password == "" {
		return nil, errors.New("password is required")
	}

	// Get user by mobile (search globally across all tenants)
	user, err := s.userRepo.GetByMobileGlobal(ctx, mobile)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	// Check if user is active
	if !user.Active {
		return nil, errors.New("user is inactive")
	}

	// Compare password
	if err := utils.ComparePassword(user.Password, password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
