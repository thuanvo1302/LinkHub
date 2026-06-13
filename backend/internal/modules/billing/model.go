package billing

import "time"

type Plan struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Price        int       `json:"price"`
	Currency     string    `json:"currency"`
	DurationDays int       `json:"duration_days"`
	Features     []string  `json:"features"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Subscription struct {
	ID         string     `json:"id"`
	UserID     string     `json:"user_id"`
	PlanID     string     `json:"plan_id"`
	Status     string     `json:"status"`
	StartedAt  time.Time  `json:"started_at"`
	ExpiredAt  time.Time  `json:"expired_at"`
	CanceledAt *time.Time `json:"canceled_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type Payment struct {
	ID                string    `json:"id"`
	UserID            string    `json:"user_id"`
	Provider          string    `json:"provider"`
	ProviderPaymentID string    `json:"provider_payment_id"`
	Amount            int       `json:"amount"`
	Currency          string    `json:"currency"`
	Status            string    `json:"status"`
	CheckoutURL       string    `json:"checkout_url"`
	PaidAt            time.Time `json:"paid_at,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type CreateCheckoutRequest struct {
	PlanID string `json:"plan_id"`
}
