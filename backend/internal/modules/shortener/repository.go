package shortener

import (
	"errors"
	"time"

	"linkhub/backend/internal/database"
)

type Repository struct {
	store *database.Store
}

func NewRepository(store *database.Store) *Repository {
	return &Repository{store: store}
}

func (r *Repository) ListByUserID(userID string) []ShortLink {
	r.store.Mu.RLock()
	defer r.store.Mu.RUnlock()

	items := []ShortLink{}
	for _, record := range r.store.ShortLinksByID {
		if record["user_id"].(string) == userID {
			items = append(items, mapToShortLink(record))
		}
	}
	return items
}

func (r *Repository) Create(item ShortLink) ShortLink {
	r.store.Mu.Lock()
	defer r.store.Mu.Unlock()
	r.store.ShortLinksByID[item.ID] = shortLinkToMap(item)
	r.store.ShortLinksByCode[item.Code] = item.ID
	return item
}

func (r *Repository) FindByID(id string) (ShortLink, error) {
	r.store.Mu.RLock()
	defer r.store.Mu.RUnlock()

	record, ok := r.store.ShortLinksByID[id]
	if !ok {
		return ShortLink{}, errors.New("short link not found")
	}
	return mapToShortLink(record), nil
}

func (r *Repository) FindByCode(code string) (ShortLink, error) {
	r.store.Mu.RLock()
	defer r.store.Mu.RUnlock()

	id, ok := r.store.ShortLinksByCode[code]
	if !ok {
		return ShortLink{}, errors.New("short link not found")
	}
	return mapToShortLink(r.store.ShortLinksByID[id]), nil
}

func (r *Repository) CodeTaken(code, currentID string) bool {
	r.store.Mu.RLock()
	defer r.store.Mu.RUnlock()
	id, ok := r.store.ShortLinksByCode[code]
	return ok && id != currentID
}

func (r *Repository) Update(item ShortLink) ShortLink {
	r.store.Mu.Lock()
	defer r.store.Mu.Unlock()

	if currentID, ok := r.store.ShortLinksByCode[item.Code]; ok && currentID != item.ID {
		delete(r.store.ShortLinksByCode, item.Code)
	}
	r.store.ShortLinksByID[item.ID] = shortLinkToMap(item)
	r.store.ShortLinksByCode[item.Code] = item.ID
	return item
}

func (r *Repository) Delete(id string) error {
	r.store.Mu.Lock()
	defer r.store.Mu.Unlock()

	record, ok := r.store.ShortLinksByID[id]
	if !ok {
		return errors.New("short link not found")
	}
	delete(r.store.ShortLinksByCode, record["code"].(string))
	delete(r.store.ShortLinksByID, id)
	return nil
}

func (r *Repository) IncrementClicks(id string) {
	r.store.Mu.Lock()
	defer r.store.Mu.Unlock()
	r.store.ClickCounts[id]++
}

func (r *Repository) CountClicks(userID string) int {
	r.store.Mu.RLock()
	defer r.store.Mu.RUnlock()

	total := 0
	for id, record := range r.store.ShortLinksByID {
		if record["user_id"].(string) == userID {
			total += r.store.ClickCounts[id]
		}
	}
	return total
}

func (r *Repository) CountClicksForLink(id string) int {
	r.store.Mu.RLock()
	defer r.store.Mu.RUnlock()
	return r.store.ClickCounts[id]
}

func (r *Repository) TopLinks(userID string) []map[string]any {
	r.store.Mu.RLock()
	defer r.store.Mu.RUnlock()

	results := []map[string]any{}
	for id, record := range r.store.ShortLinksByID {
		if record["user_id"].(string) == userID {
			results = append(results, map[string]any{
				"id":     id,
				"code":   record["code"].(string),
				"title":  record["title"].(string),
				"clicks": r.store.ClickCounts[id],
			})
		}
	}
	return results
}

func shortLinkToMap(item ShortLink) map[string]any {
	return map[string]any{
		"id":           item.ID,
		"user_id":      item.UserID,
		"code":         item.Code,
		"original_url": item.OriginalURL,
		"title":        item.Title,
		"expires_at":   item.ExpiresAt,
		"is_active":    item.IsActive,
		"created_at":   item.CreatedAt,
		"updated_at":   item.UpdatedAt,
	}
}

func mapToShortLink(record map[string]any) ShortLink {
	var expiresAt *time.Time
	if record["expires_at"] != nil {
		expiresAt = record["expires_at"].(*time.Time)
	}
	return ShortLink{
		ID:          record["id"].(string),
		UserID:      record["user_id"].(string),
		Code:        record["code"].(string),
		OriginalURL: record["original_url"].(string),
		Title:       record["title"].(string),
		ExpiresAt:   expiresAt,
		IsActive:    record["is_active"].(bool),
		CreatedAt:   record["created_at"].(time.Time),
		UpdatedAt:   record["updated_at"].(time.Time),
	}
}
