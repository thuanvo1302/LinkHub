package billing

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

func (h *Handler) Plans(w http.ResponseWriter, r *http.Request) {
	response.OK(w, h.service.Plans(), "OK")
}

func (h *Handler) CreateCheckout(w http.ResponseWriter, r *http.Request) {
	var input CreateCheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
		return
	}
	data, err := h.service.CreateCheckout(middleware.UserIDFromContext(r.Context()), input)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "PAYMENT_FAILED", err.Error(), nil)
		return
	}
	response.OK(w, data, "Checkout created")
}

func (h *Handler) MockSuccess(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		PaymentID string `json:"payment_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
		return
	}
	data, err := h.service.MockSuccess(middleware.UserIDFromContext(r.Context()), payload.PaymentID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "PAYMENT_FAILED", err.Error(), nil)
		return
	}
	response.OK(w, data, "Subscription activated")
}

func (h *Handler) History(w http.ResponseWriter, r *http.Request) {
	response.OK(w, h.service.History(middleware.UserIDFromContext(r.Context())), "OK")
}

func (h *Handler) Current(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.CurrentSubscription(middleware.UserIDFromContext(r.Context()))
	if err != nil {
		response.Error(w, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}
	response.OK(w, data, "OK")
}

func (h *Handler) Cancel(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.Cancel(middleware.UserIDFromContext(r.Context()))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}
	response.OK(w, data, "Subscription canceled")
}
