package link

import "time"

type ProfileLink struct {
	ID        string    `json:"id"`
	ProfileID string    `json:"profile_id"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	Icon      string    `json:"icon"`
	Position  int       `json:"position"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateLinkRequest struct {
	ProfileID string `json:"profile_id"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	Icon      string `json:"icon"`
}

type UpdateLinkRequest struct {
	Title    string `json:"title"`
	URL      string `json:"url"`
	Icon     string `json:"icon"`
	IsActive bool   `json:"is_active"`
}

type ReorderRequest struct {
	IDs []string `json:"ids"`
}
