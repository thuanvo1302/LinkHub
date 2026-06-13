package auth

import (
	"encoding/json"
	"net/http"

	"linkhub/backend/internal/middleware"
	"linkhub/backend/internal/pkg/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var input RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
		return
	}

	data, err := h.service.Register(input)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	response.Created(w, data, "Registered successfully")
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var input LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
		return
	}

	data, err := h.service.Login(input)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", err.Error(), nil)
		return
	}

	response.OK(w, data, "Login successful")
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var input RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
		return
	}

	data, err := h.service.Refresh(input)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", err.Error(), nil)
		return
	}

	response.OK(w, data, "Token refreshed")
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	var input RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
		return
	}

	h.service.Logout(input)
	response.OK(w, map[string]bool{"logged_out": true}, "Logout successful")
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	user, err := h.service.Me(middleware.UserIDFromContext(r.Context()))
	if err != nil {
		response.Error(w, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}

	response.OK(w, user, "OK")
}
