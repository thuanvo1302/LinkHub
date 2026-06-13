package link

import (
	"encoding/json"
	"net/http"
	"strings"

	"linkhub/backend/internal/middleware"
	"linkhub/backend/internal/pkg/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.List(middleware.UserIDFromContext(r.Context()))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}
	response.OK(w, data, "OK")
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input CreateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
		return
	}

	data, err := h.service.Create(middleware.UserIDFromContext(r.Context()), input)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}
	response.Created(w, data, "Link created")
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/profile-links/")
	var input UpdateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
		return
	}

	data, err := h.service.Update(middleware.UserIDFromContext(r.Context()), id, input)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}
	response.OK(w, data, "Link updated")
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/profile-links/")
	if err := h.service.Delete(middleware.UserIDFromContext(r.Context()), id); err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}
	response.OK(w, map[string]string{"id": id}, "Link deleted")
}

func (h *Handler) Reorder(w http.ResponseWriter, r *http.Request) {
	var input ReorderRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
		return
	}

	if err := h.service.Reorder(middleware.UserIDFromContext(r.Context()), input.IDs); err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}
	response.OK(w, map[string]any{"ids": input.IDs}, "Links reordered")
}

var _ = middleware.UserIDFromContext
