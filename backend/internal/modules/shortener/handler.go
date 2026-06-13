package shortener

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
	response.OK(w, h.service.List(middleware.UserIDFromContext(r.Context())), "OK")
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
		return
	}

	data, err := h.service.Create(middleware.UserIDFromContext(r.Context()), input)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}
	response.Created(w, data, "Short link created")
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/short-links/")
	data, err := h.service.Get(middleware.UserIDFromContext(r.Context()), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}
	response.OK(w, data, "OK")
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/short-links/")
	var input UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
		return
	}

	data, err := h.service.Update(middleware.UserIDFromContext(r.Context()), id, input)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}
	response.OK(w, data, "Short link updated")
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/short-links/")
	if err := h.service.Delete(middleware.UserIDFromContext(r.Context()), id); err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}
	response.OK(w, map[string]string{"id": id}, "Short link deleted")
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/")
	if code == "" {
		response.OK(w, map[string]string{"message": "welcome"}, "OK")
		return
	}

	item, err := h.service.Resolve(code)
	if err != nil {
		response.Error(w, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}
	http.Redirect(w, r, item.OriginalURL, http.StatusFound)
}

func (h *Handler) Overview(w http.ResponseWriter, r *http.Request) {
	response.OK(w, h.service.Overview(middleware.UserIDFromContext(r.Context())), "OK")
}

func (h *Handler) LinkAnalytics(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/analytics/short-links/")
	item, err := h.service.Get(middleware.UserIDFromContext(r.Context()), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}
	response.OK(w, map[string]any{
		"id":           item.ID,
		"code":         item.Code,
		"title":        item.Title,
		"clicks":       h.service.repo.CountClicksForLink(item.ID),
		"original_url": item.OriginalURL,
	}, "OK")
}
