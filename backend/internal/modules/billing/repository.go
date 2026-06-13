package billing

import (
	"errors"
	"time"

	"linkhub/backend/internal/database"
)

type Repository struct {
	store *database.Store
}

func NewRepository(store *database.Store) *Repository {
	repo := &Repository{store: store}
	repo.seedPlans()
	return repo
}

func (r *Repository) seedPlans() {
	r.store.Mu.Lock()
	defer r.store.Mu.Unlock()
	if len(r.store.PlansByID) > 0 {
		return
	}

	now := time.Now()
	r.store.PlansByID["free"] = map[string]any{
		"id":            "free",
		"name":          "Free",
		"price":         0,
		"currency":      "VND",
		"duration_days": 30,
		"features":      []string{"1 profile", "5 profile links", "10 short links", "7 days analytics"},
		"created_at":    now,
		"updated_at":    now,
	}
	r.store.PlansByID["pro"] = map[string]any{
		"id":            "pro",
		"name":          "Pro",
		"price":         99000,
		"currency":      "VND",
		"duration_days": 30,
		"features":      []string{"Unlimited links", "Custom slug", "QR code", "Expire link", "Full analytics"},
		"created_at":    now,
		"updated_at":    now,
	}
}

func (r *Repository) ListPlans() []Plan {
	r.store.Mu.RLock()
	defer r.store.Mu.RUnlock()
	plans := []Plan{}
	for _, record := range r.store.PlansByID {
		plans = append(plans, mapToPlan(record))
	}
	return plans
}

func (r *Repository) FindPlan(id string) (Plan, error) {
	r.store.Mu.RLock()
	defer r.store.Mu.RUnlock()
	record, ok := r.store.PlansByID[id]
	if !ok {
		return Plan{}, errors.New("plan not found")
	}
	return mapToPlan(record), nil
}

func (r *Repository) CreatePayment(payment Payment) Payment {
	r.store.Mu.Lock()
	defer r.store.Mu.Unlock()
	r.store.PaymentsByID[payment.ID] = paymentToMap(payment)
	return payment
}

func (r *Repository) FindPayment(id string) (Payment, error) {
	r.store.Mu.RLock()
	defer r.store.Mu.RUnlock()
	record, ok := r.store.PaymentsByID[id]
	if !ok {
		return Payment{}, errors.New("payment not found")
	}
	return mapToPayment(record), nil
}

func (r *Repository) UpdatePayment(payment Payment) Payment {
	r.store.Mu.Lock()
	defer r.store.Mu.Unlock()
	r.store.PaymentsByID[payment.ID] = paymentToMap(payment)
	return payment
}

func (r *Repository) ListPaymentsByUserID(userID string) []Payment {
	r.store.Mu.RLock()
	defer r.store.Mu.RUnlock()
	items := []Payment{}
	for _, record := range r.store.PaymentsByID {
		if record["user_id"].(string) == userID {
			items = append(items, mapToPayment(record))
		}
	}
	return items
}

func (r *Repository) UpsertSubscription(sub Subscription) Subscription {
	r.store.Mu.Lock()
	defer r.store.Mu.Unlock()
	r.store.SubscriptionsByID[sub.ID] = subscriptionToMap(sub)
	if user, ok := r.store.UsersByID[sub.UserID]; ok {
		if sub.Status == "active" {
			user["plan"] = sub.PlanID
		} else if sub.Status == "canceled" {
			user["plan"] = "free"
		}
		user["updated_at"] = time.Now()
	}
	return sub
}

func (r *Repository) FindCurrentSubscription(userID string) (Subscription, error) {
	r.store.Mu.RLock()
	defer r.store.Mu.RUnlock()
	var latest Subscription
	found := false
	for _, record := range r.store.SubscriptionsByID {
		if record["user_id"].(string) == userID {
			sub := mapToSubscription(record)
			if !found || sub.CreatedAt.After(latest.CreatedAt) {
				latest = sub
				found = true
			}
		}
	}
	if !found {
		return Subscription{}, errors.New("subscription not found")
	}
	return latest, nil
}

func mapToPlan(record map[string]any) Plan {
	return Plan{
		ID:           record["id"].(string),
		Name:         record["name"].(string),
		Price:        record["price"].(int),
		Currency:     record["currency"].(string),
		DurationDays: record["duration_days"].(int),
		Features:     record["features"].([]string),
		CreatedAt:    record["created_at"].(time.Time),
		UpdatedAt:    record["updated_at"].(time.Time),
	}
}

func paymentToMap(payment Payment) map[string]any {
	return map[string]any{
		"id":                  payment.ID,
		"user_id":             payment.UserID,
		"provider":            payment.Provider,
		"provider_payment_id": payment.ProviderPaymentID,
		"amount":              payment.Amount,
		"currency":            payment.Currency,
		"status":              payment.Status,
		"checkout_url":        payment.CheckoutURL,
		"paid_at":             payment.PaidAt,
		"created_at":          payment.CreatedAt,
		"updated_at":          payment.UpdatedAt,
	}
}

func mapToPayment(record map[string]any) Payment {
	return Payment{
		ID:                record["id"].(string),
		UserID:            record["user_id"].(string),
		Provider:          record["provider"].(string),
		ProviderPaymentID: record["provider_payment_id"].(string),
		Amount:            record["amount"].(int),
		Currency:          record["currency"].(string),
		Status:            record["status"].(string),
		CheckoutURL:       record["checkout_url"].(string),
		PaidAt:            record["paid_at"].(time.Time),
		CreatedAt:         record["created_at"].(time.Time),
		UpdatedAt:         record["updated_at"].(time.Time),
	}
}

func subscriptionToMap(sub Subscription) map[string]any {
	return map[string]any{
		"id":          sub.ID,
		"user_id":     sub.UserID,
		"plan_id":     sub.PlanID,
		"status":      sub.Status,
		"started_at":  sub.StartedAt,
		"expired_at":  sub.ExpiredAt,
		"canceled_at": sub.CanceledAt,
		"created_at":  sub.CreatedAt,
		"updated_at":  sub.UpdatedAt,
	}
}

func mapToSubscription(record map[string]any) Subscription {
	var canceledAt *time.Time
	if record["canceled_at"] != nil {
		canceledAt = record["canceled_at"].(*time.Time)
	}
	return Subscription{
		ID:         record["id"].(string),
		UserID:     record["user_id"].(string),
		PlanID:     record["plan_id"].(string),
		Status:     record["status"].(string),
		StartedAt:  record["started_at"].(time.Time),
		ExpiredAt:  record["expired_at"].(time.Time),
		CanceledAt: canceledAt,
		CreatedAt:  record["created_at"].(time.Time),
		UpdatedAt:  record["updated_at"].(time.Time),
	}
}
