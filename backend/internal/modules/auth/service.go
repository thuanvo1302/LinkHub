package auth

import (
	"errors"
	"strings"
	"time"

	"linkhub/backend/internal/config"
	"linkhub/backend/internal/pkg/jwt"
	"linkhub/backend/internal/pkg/password"
	"linkhub/backend/internal/pkg/validator"
)

type Service struct {
	repo *Repository
	cfg  config.Config
}

func NewService(repo *Repository, cfg config.Config) *Service {
	return &Service{repo: repo, cfg: cfg}
}

func (s *Service) Register(input RegisterRequest) (AuthResponse, error) {
	if err := validator.Required(input.Email, "email"); err != nil {
		return AuthResponse{}, err
	}
	if err := validator.Email(input.Email); err != nil {
		return AuthResponse{}, err
	}
	if err := validator.Required(input.Password, "password"); err != nil {
		return AuthResponse{}, err
	}
	if err := validator.Password(input.Password); err != nil {
		return AuthResponse{}, err
	}
	if err := validator.Required(input.FullName, "full_name"); err != nil {
		return AuthResponse{}, err
	}

	now := time.Now()
	user := User{
		ID:           s.repo.store.NextUserID(),
		Email:        strings.ToLower(strings.TrimSpace(input.Email)),
		PasswordHash: password.Hash(input.Password),
		FullName:     strings.TrimSpace(input.FullName),
		Role:         "user",
		Plan:         "free",
		Status:       "active",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	created, err := s.repo.Create(user)
	if err != nil {
		return AuthResponse{}, err
	}

	accessToken, err := s.issueToken(created)
	if err != nil {
		return AuthResponse{}, err
	}

	refreshToken := s.issueRefreshToken(created)

	return AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         created,
	}, nil
}

func (s *Service) Login(input LoginRequest) (AuthResponse, error) {
	user, err := s.repo.FindByEmail(strings.ToLower(strings.TrimSpace(input.Email)))
	if err != nil {
		return AuthResponse{}, errors.New("invalid credentials")
	}

	if !password.Compare(input.Password, user.PasswordHash) {
		return AuthResponse{}, errors.New("invalid credentials")
	}

	accessToken, err := s.issueToken(user)
	if err != nil {
		return AuthResponse{}, err
	}

	refreshToken := s.issueRefreshToken(user)

	return AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (s *Service) Me(userID string) (User, error) {
	return s.repo.FindByID(userID)
}

func (s *Service) Refresh(input RefreshRequest) (AuthResponse, error) {
	userID, err := s.repo.FindRefreshToken(strings.TrimSpace(input.RefreshToken))
	if err != nil {
		return AuthResponse{}, errors.New("invalid refresh token")
	}

	user, err := s.repo.FindByID(userID)
	if err != nil {
		return AuthResponse{}, err
	}

	s.repo.DeleteRefreshToken(input.RefreshToken)
	accessToken, err := s.issueToken(user)
	if err != nil {
		return AuthResponse{}, err
	}

	refreshToken := s.issueRefreshToken(user)
	return AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (s *Service) Logout(input RefreshRequest) {
	s.repo.DeleteRefreshToken(strings.TrimSpace(input.RefreshToken))
}

func (s *Service) issueToken(user User) (string, error) {
	return jwt.Sign(jwt.Claims{
		UserID: user.ID,
		Email:  user.Email,
		Exp:    time.Now().Add(24 * time.Hour).Unix(),
	}, s.cfg.JWTSecret)
}

func (s *Service) issueRefreshToken(user User) string {
	token := user.ID + "_" + jwtTokenEntropy()
	s.repo.SaveRefreshToken(token, user.ID, time.Now().Add(7*24*time.Hour))
	return token
}

func jwtTokenEntropy() string {
	return strings.ReplaceAll(time.Now().Format("20060102150405.000000000"), ".", "")
}
