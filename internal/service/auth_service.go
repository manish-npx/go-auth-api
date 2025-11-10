package service

import (
	"context"
	"errors"
	"time"

	utils "github.com/manish-npx/go-auth-api/internal/config"
	"github.com/manish-npx/go-auth-api/internal/model"
	"github.com/manish-npx/go-auth-api/internal/repository"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthService struct {
	cfg      *utils.Config
	repo     *repository.UserRepository
	tokenTTL time.Duration
}

func NewAuthService(cfg *utils.Config, repo *repository.UserRepository) *AuthService {
	return &AuthService{
		cfg:      cfg,
		repo:     repo,
		tokenTTL: 72 * time.Hour,
	}
}

func (s *AuthService) Register(ctx context.Context, name, email, password string) (*model.User, error) {
	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}
	u := &model.User{
		Name:         name,
		Email:        email,
		PasswordHash: hash,
	}
	if err := s.repo.New(ctx, u); err != nil {
		return nil, err
	}
	return s.repo.Get(ctx, email)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, *model.User, error) {
	u, err := s.repo.Get(ctx, email)
	if err != nil {
		return "", nil, ErrInvalidCredentials
	}
	if !utils.CheckPassword(password, u.PasswordHash) {
		return "", nil, ErrInvalidCredentials
	}
	token, err := utils.GenerateToken(s.cfg.Security.JWTSecret, u.ID, s.tokenTTL)
	if err != nil {
		return "", nil, err
	}
	return token, u, nil
}
