package auth

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

func (r *Repository) Create(user User) (User, error) {
	r.store.Mu.Lock()
	defer r.store.Mu.Unlock()

	if _, exists := r.store.UsersByEmail[user.Email]; exists {
		return User{}, errors.New("email already exists")
	}

	record := map[string]any{
		"id":            user.ID,
		"email":         user.Email,
		"password_hash": user.PasswordHash,
		"full_name":     user.FullName,
		"avatar_url":    user.AvatarURL,
		"role":          user.Role,
		"plan":          user.Plan,
		"status":        user.Status,
		"created_at":    user.CreatedAt,
		"updated_at":    user.UpdatedAt,
	}

	r.store.UsersByID[user.ID] = record
	r.store.UsersByEmail[user.Email] = user.ID
	return user, nil
}

func (r *Repository) FindByEmail(email string) (User, error) {
	r.store.Mu.RLock()
	defer r.store.Mu.RUnlock()

	id, ok := r.store.UsersByEmail[email]
	if !ok {
		return User{}, errors.New("user not found")
	}

	return mapToUser(r.store.UsersByID[id]), nil
}

func (r *Repository) FindByID(id string) (User, error) {
	r.store.Mu.RLock()
	defer r.store.Mu.RUnlock()

	record, ok := r.store.UsersByID[id]
	if !ok {
		return User{}, errors.New("user not found")
	}

	return mapToUser(record), nil
}

func (r *Repository) SaveRefreshToken(token string, userID string, expiresAt time.Time) {
	r.store.Mu.Lock()
	defer r.store.Mu.Unlock()
	r.store.RefreshTokens[token] = map[string]any{
		"user_id":    userID,
		"expires_at": expiresAt,
	}
}

func (r *Repository) FindRefreshToken(token string) (string, error) {
	r.store.Mu.RLock()
	defer r.store.Mu.RUnlock()

	record, ok := r.store.RefreshTokens[token]
	if !ok {
		return "", errors.New("refresh token not found")
	}

	if time.Now().After(record["expires_at"].(time.Time)) {
		return "", errors.New("refresh token expired")
	}

	return record["user_id"].(string), nil
}

func (r *Repository) DeleteRefreshToken(token string) {
	r.store.Mu.Lock()
	defer r.store.Mu.Unlock()
	delete(r.store.RefreshTokens, token)
}

func mapToUser(record map[string]any) User {
	return User{
		ID:           record["id"].(string),
		Email:        record["email"].(string),
		PasswordHash: record["password_hash"].(string),
		FullName:     record["full_name"].(string),
		AvatarURL:    record["avatar_url"].(string),
		Role:         record["role"].(string),
		Plan:         record["plan"].(string),
		Status:       record["status"].(string),
		CreatedAt:    record["created_at"].(time.Time),
		UpdatedAt:    record["updated_at"].(time.Time),
	}
}
