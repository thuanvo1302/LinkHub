package profile

import (
	"errors"
	"strings"
	"time"

	"linkhub/backend/internal/pkg/slug"
	"linkhub/backend/internal/pkg/validator"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetMe(userID string) (Profile, error) {
	return s.repo.FindByUserID(userID)
}

func (s *Service) UpdateMe(userID string, input UpdateProfileRequest) (Profile, error) {
	current, err := s.repo.FindByUserID(userID)
	if err != nil {
		current = Profile{
			ID:        s.repo.store.NextProfileID(),
			UserID:    userID,
			Theme:     "sunset-grid",
			IsPublic:  true,
			CreatedAt: time.Now(),
		}
	}

	username := slug.Normalize(input.Username)
	if username != "" {
		if err := validator.Username(username); err != nil {
			return Profile{}, err
		}
		if slug.Reserved(username) {
			return Profile{}, errors.New("username is reserved")
		}
		if s.repo.UsernameTaken(username, userID) {
			return Profile{}, errors.New("username already taken")
		}
	}

	current.Username = username
	current.DisplayName = strings.TrimSpace(input.DisplayName)
	current.Bio = strings.TrimSpace(input.Bio)
	current.AvatarURL = strings.TrimSpace(input.AvatarURL)
	if input.Theme != "" {
		current.Theme = input.Theme
	}
	current.IsPublic = input.IsPublic
	current.UpdatedAt = time.Now()

	return s.repo.Upsert(current), nil
}

func (s *Service) GetPublic(username string) (Profile, error) {
	profile, err := s.repo.FindByUsername(slug.Normalize(username))
	if err != nil {
		return Profile{}, err
	}
	if !profile.IsPublic {
		return Profile{}, errors.New("profile is private")
	}
	return profile, nil
}
