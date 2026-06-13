package profile

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

func (r *Repository) Upsert(profile Profile) Profile {
	r.store.Mu.Lock()
	defer r.store.Mu.Unlock()

	if existing, ok := r.store.ProfilesByUserID[profile.UserID]; ok {
		if oldUsername, ok := existing["username"].(string); ok && oldUsername != "" && oldUsername != profile.Username {
			delete(r.store.ProfilesByName, oldUsername)
		}
	}

	r.store.ProfilesByUserID[profile.UserID] = map[string]any{
		"id":           profile.ID,
		"user_id":      profile.UserID,
		"username":     profile.Username,
		"display_name": profile.DisplayName,
		"bio":          profile.Bio,
		"avatar_url":   profile.AvatarURL,
		"theme":        profile.Theme,
		"is_public":    profile.IsPublic,
		"created_at":   profile.CreatedAt,
		"updated_at":   profile.UpdatedAt,
	}

	if profile.Username != "" {
		r.store.ProfilesByName[profile.Username] = profile.UserID
	}

	return profile
}

func (r *Repository) FindByUserID(userID string) (Profile, error) {
	r.store.Mu.RLock()
	defer r.store.Mu.RUnlock()

	record, ok := r.store.ProfilesByUserID[userID]
	if !ok {
		return Profile{}, errors.New("profile not found")
	}

	return mapToProfile(record), nil
}

func (r *Repository) FindByUsername(username string) (Profile, error) {
	r.store.Mu.RLock()
	defer r.store.Mu.RUnlock()

	userID, ok := r.store.ProfilesByName[username]
	if !ok {
		return Profile{}, errors.New("profile not found")
	}

	return mapToProfile(r.store.ProfilesByUserID[userID]), nil
}

func (r *Repository) UsernameTaken(username, ownerUserID string) bool {
	r.store.Mu.RLock()
	defer r.store.Mu.RUnlock()

	userID, ok := r.store.ProfilesByName[username]
	return ok && userID != ownerUserID
}

func mapToProfile(record map[string]any) Profile {
	return Profile{
		ID:          record["id"].(string),
		UserID:      record["user_id"].(string),
		Username:    record["username"].(string),
		DisplayName: record["display_name"].(string),
		Bio:         record["bio"].(string),
		AvatarURL:   record["avatar_url"].(string),
		Theme:       record["theme"].(string),
		IsPublic:    record["is_public"].(bool),
		CreatedAt:   record["created_at"].(time.Time),
		UpdatedAt:   record["updated_at"].(time.Time),
	}
}
