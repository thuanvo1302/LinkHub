package link

import (
	"errors"
	"slices"
	"time"

	"linkhub/backend/internal/database"
)

type Repository struct {
	store *database.Store
}

func NewRepository(store *database.Store) *Repository {
	return &Repository{store: store}
}

func (r *Repository) ListByProfileID(profileID string) []ProfileLink {
	r.store.Mu.RLock()
	defer r.store.Mu.RUnlock()

	ids := r.store.ProfileLinksOrder[profileID]
	links := make([]ProfileLink, 0, len(ids))
	for _, id := range ids {
		links = append(links, mapToLink(r.store.ProfileLinksByID[id]))
	}
	return links
}

func (r *Repository) Create(link ProfileLink) ProfileLink {
	r.store.Mu.Lock()
	defer r.store.Mu.Unlock()

	r.store.ProfileLinksByID[link.ID] = linkToMap(link)
	r.store.ProfileLinksOrder[link.ProfileID] = append(r.store.ProfileLinksOrder[link.ProfileID], link.ID)
	return link
}

func (r *Repository) FindByID(id string) (ProfileLink, error) {
	r.store.Mu.RLock()
	defer r.store.Mu.RUnlock()

	record, ok := r.store.ProfileLinksByID[id]
	if !ok {
		return ProfileLink{}, errors.New("link not found")
	}
	return mapToLink(record), nil
}

func (r *Repository) Update(link ProfileLink) ProfileLink {
	r.store.Mu.Lock()
	defer r.store.Mu.Unlock()
	r.store.ProfileLinksByID[link.ID] = linkToMap(link)
	return link
}

func (r *Repository) Delete(id string) error {
	r.store.Mu.Lock()
	defer r.store.Mu.Unlock()

	record, ok := r.store.ProfileLinksByID[id]
	if !ok {
		return errors.New("link not found")
	}

	profileID := record["profile_id"].(string)
	delete(r.store.ProfileLinksByID, id)

	ids := r.store.ProfileLinksOrder[profileID]
	index := slices.Index(ids, id)
	if index >= 0 {
		r.store.ProfileLinksOrder[profileID] = append(ids[:index], ids[index+1:]...)
	}
	return nil
}

func (r *Repository) Reorder(profileID string, ids []string) {
	r.store.Mu.Lock()
	defer r.store.Mu.Unlock()
	r.store.ProfileLinksOrder[profileID] = ids
	for index, id := range ids {
		if record, ok := r.store.ProfileLinksByID[id]; ok {
			record["position"] = index + 1
			record["updated_at"] = time.Now()
		}
	}
}

func linkToMap(link ProfileLink) map[string]any {
	return map[string]any{
		"id":         link.ID,
		"profile_id": link.ProfileID,
		"title":      link.Title,
		"url":        link.URL,
		"icon":       link.Icon,
		"position":   link.Position,
		"is_active":  link.IsActive,
		"created_at": link.CreatedAt,
		"updated_at": link.UpdatedAt,
	}
}

func mapToLink(record map[string]any) ProfileLink {
	return ProfileLink{
		ID:        record["id"].(string),
		ProfileID: record["profile_id"].(string),
		Title:     record["title"].(string),
		URL:       record["url"].(string),
		Icon:      record["icon"].(string),
		Position:  record["position"].(int),
		IsActive:  record["is_active"].(bool),
		CreatedAt: record["created_at"].(time.Time),
		UpdatedAt: record["updated_at"].(time.Time),
	}
}
