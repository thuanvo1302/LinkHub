package shortener

import (
	"errors"
	"strings"
	"time"

	"linkhub/backend/internal/config"
	"linkhub/backend/internal/pkg/slug"
	"linkhub/backend/internal/pkg/validator"
)

type Service struct {
	repo *Repository
	cfg  config.Config
}

func NewService(repo *Repository, cfg config.Config) *Service {
	return &Service{repo: repo, cfg: cfg}
}

func (s *Service) List(userID string) []ShortLink {
	return s.repo.ListByUserID(userID)
}

func (s *Service) Create(userID string, input CreateRequest) (ShortLink, error) {
	if err := validator.URL(input.OriginalURL); err != nil {
		return ShortLink{}, err
	}

	code := slug.Normalize(strings.TrimSpace(input.Code))
	if code == "" {
		code = slug.Random(6)
	}
	if slug.Reserved(code) {
		return ShortLink{}, errors.New("code is reserved")
	}
	if s.repo.CodeTaken(code, "") {
		return ShortLink{}, errors.New("code already exists")
	}

	now := time.Now()
	item := ShortLink{
		ID:          s.repo.store.NextShortLinkID(),
		UserID:      userID,
		Code:        code,
		OriginalURL: strings.TrimSpace(input.OriginalURL),
		Title:       strings.TrimSpace(input.Title),
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	return s.repo.Create(item), nil
}

func (s *Service) Get(userID, id string) (ShortLink, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return ShortLink{}, err
	}
	if item.UserID != userID {
		return ShortLink{}, errors.New("short link does not belong to user")
	}
	return item, nil
}

func (s *Service) Update(userID, id string, input UpdateRequest) (ShortLink, error) {
	item, err := s.Get(userID, id)
	if err != nil {
		return ShortLink{}, err
	}
	if err := validator.URL(input.OriginalURL); err != nil {
		return ShortLink{}, err
	}

	item.OriginalURL = strings.TrimSpace(input.OriginalURL)
	item.Title = strings.TrimSpace(input.Title)
	item.IsActive = input.IsActive
	item.UpdatedAt = time.Now()
	return s.repo.Update(item), nil
}

func (s *Service) Delete(userID, id string) error {
	item, err := s.Get(userID, id)
	if err != nil {
		return err
	}
	return s.repo.Delete(item.ID)
}

func (s *Service) Resolve(code string) (ShortLink, error) {
	item, err := s.repo.FindByCode(slug.Normalize(code))
	if err != nil {
		return ShortLink{}, err
	}
	if !item.IsActive {
		return ShortLink{}, errors.New("short link is disabled")
	}
	if item.ExpiresAt != nil && time.Now().After(*item.ExpiresAt) {
		return ShortLink{}, errors.New("short link is expired")
	}
	s.repo.IncrementClicks(item.ID)
	return item, nil
}

func (s *Service) Overview(userID string) map[string]any {
	return map[string]any{
		"total_short_links": len(s.repo.ListByUserID(userID)),
		"total_clicks":      s.repo.CountClicks(userID),
		"top_links":         s.repo.TopLinks(userID),
		"app_base_url":      s.cfg.AppBaseURL,
	}
}
