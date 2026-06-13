package link

import (
	"errors"
	"strings"
	"time"

	"linkhub/backend/internal/pkg/validator"
)

type Service struct {
	repo              *Repository
	profileIDResolver func(userID string) (string, error)
}

func NewService(repo *Repository, profileIDResolver func(userID string) (string, error)) *Service {
	return &Service{repo: repo, profileIDResolver: profileIDResolver}
}

func (s *Service) List(userID string) ([]ProfileLink, error) {
	profileID, err := s.profileIDByUserID(userID)
	if err != nil {
		return nil, err
	}
	return s.repo.ListByProfileID(profileID), nil
}

func (s *Service) Create(userID string, input CreateLinkRequest) (ProfileLink, error) {
	if err := validator.Required(input.Title, "title"); err != nil {
		return ProfileLink{}, err
	}
	if err := validator.URL(input.URL); err != nil {
		return ProfileLink{}, err
	}

	profileID, err := s.profileIDByUserID(userID)
	if err != nil {
		return ProfileLink{}, err
	}

	links := s.repo.ListByProfileID(profileID)
	now := time.Now()
	link := ProfileLink{
		ID:        s.repo.store.NextProfileLinkID(),
		ProfileID: profileID,
		Title:     strings.TrimSpace(input.Title),
		URL:       strings.TrimSpace(input.URL),
		Icon:      strings.TrimSpace(input.Icon),
		Position:  len(links) + 1,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return s.repo.Create(link), nil
}

func (s *Service) Update(userID, id string, input UpdateLinkRequest) (ProfileLink, error) {
	profileID, err := s.profileIDByUserID(userID)
	if err != nil {
		return ProfileLink{}, err
	}

	item, err := s.repo.FindByID(id)
	if err != nil {
		return ProfileLink{}, err
	}
	if item.ProfileID != profileID {
		return ProfileLink{}, errors.New("link does not belong to profile")
	}
	if err := validator.Required(input.Title, "title"); err != nil {
		return ProfileLink{}, err
	}
	if err := validator.URL(input.URL); err != nil {
		return ProfileLink{}, err
	}

	item.Title = strings.TrimSpace(input.Title)
	item.URL = strings.TrimSpace(input.URL)
	item.Icon = strings.TrimSpace(input.Icon)
	item.IsActive = input.IsActive
	item.UpdatedAt = time.Now()
	return s.repo.Update(item), nil
}

func (s *Service) Delete(userID, id string) error {
	profileID, err := s.profileIDByUserID(userID)
	if err != nil {
		return err
	}

	item, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if item.ProfileID != profileID {
		return errors.New("link does not belong to profile")
	}
	return s.repo.Delete(id)
}

func (s *Service) Reorder(userID string, ids []string) error {
	profileID, err := s.profileIDByUserID(userID)
	if err != nil {
		return err
	}

	current := s.repo.ListByProfileID(profileID)
	if len(current) != len(ids) {
		return errors.New("reorder payload size mismatch")
	}
	owned := map[string]struct{}{}
	for _, item := range current {
		owned[item.ID] = struct{}{}
	}
	for _, id := range ids {
		if _, ok := owned[id]; !ok {
			return errors.New("reorder payload contains invalid id")
		}
	}
	s.repo.Reorder(profileID, ids)
	return nil
}

func (s *Service) profileIDByUserID(userID string) (string, error) {
	if s.profileIDResolver == nil {
		return "", errors.New("profile resolver is not configured")
	}
	return s.profileIDResolver(userID)
}
