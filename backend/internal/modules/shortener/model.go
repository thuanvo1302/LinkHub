package shortener

import "time"

type ShortLink struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	Code        string     `json:"code"`
	OriginalURL string     `json:"original_url"`
	Title       string     `json:"title"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateRequest struct {
	Code        string `json:"code"`
	OriginalURL string `json:"original_url"`
	Title       string `json:"title"`
}

type UpdateRequest struct {
	OriginalURL string `json:"original_url"`
	Title       string `json:"title"`
	IsActive    bool   `json:"is_active"`
}
