package profile

import (
	"encoding/json"
	"net/http"
	"strings"

	"linkhub/backend/internal/middleware"
	"linkhub/backend/internal/modules/link"
	"linkhub/backend/internal/pkg/response"
)

type Handler struct {
	service  *Service
	linkRepo *link.Repository
}

func NewHandler(service *Service, linkRepo *link.Repository) *Handler {
	return &Handler{service: service, linkRepo: linkRepo}
}

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.GetMe(middleware.UserIDFromContext(r.Context()))
	if err != nil {
		response.Error(w, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}

	response.OK(w, data, "OK")
}

func (h *Handler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	var input UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
		return
	}

	data, err := h.service.UpdateMe(middleware.UserIDFromContext(r.Context()), input)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	response.OK(w, data, "Profile updated")
}

func (h *Handler) GetPublic(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimPrefix(r.URL.Path, "/api/v1/public/profiles/")
	data, err := h.service.GetPublic(username)
	if err != nil {
		response.Error(w, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}

	links := h.linkRepo.ListByProfileID(data.ID)
	active := make([]link.ProfileLink, 0, len(links))
	for _, item := range links {
		if item.IsActive {
			active = append(active, item)
		}
	}
	response.OK(w, map[string]any{
		"profile": data,
		"links":   active,
	}, "OK")
}
