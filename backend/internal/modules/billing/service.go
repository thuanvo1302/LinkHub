package billing

import (
	"errors"
	"fmt"
	"time"

	"linkhub/backend/internal/config"
)

type Service struct {
	repo *Repository
	cfg  config.Config
}

func NewService(repo *Repository, cfg config.Config) *Service {
	return &Service{repo: repo, cfg: cfg}
}

func (s *Service) Plans() []Plan {
	return s.repo.ListPlans()
}

func (s *Service) CreateCheckout(userID string, input CreateCheckoutRequest) (Payment, error) {
	plan, err := s.repo.FindPlan(input.PlanID)
	if err != nil {
		return Payment{}, err
	}
	if plan.ID == "free" {
		return Payment{}, errors.New("free plan does not require checkout")
	}

	now := time.Now()
	paymentID := s.repo.store.NextPaymentID()
	payment := Payment{
		ID:                paymentID,
		UserID:            userID,
		Provider:          "mock",
		ProviderPaymentID: "mock_" + fmt.Sprint(now.UnixNano()),
		Amount:            plan.Price,
		Currency:          plan.Currency,
		Status:            "pending",
		CheckoutURL:       fmt.Sprintf("%s/dashboard/billing?payment_id=%s&mode=mock", s.cfg.FrontendURL, paymentID),
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	return s.repo.CreatePayment(payment), nil
}

func (s *Service) MockSuccess(userID, paymentID string) (Subscription, error) {
	payment, err := s.repo.FindPayment(paymentID)
	if err != nil {
		return Subscription{}, err
	}
	if payment.UserID != userID {
		return Subscription{}, errors.New("payment does not belong to user")
	}
	if payment.Status == "paid" {
		return s.repo.FindCurrentSubscription(userID)
	}

	now := time.Now()
	payment.Status = "paid"
	payment.PaidAt = now
	payment.UpdatedAt = now
	s.repo.UpdatePayment(payment)

	planID := "pro"
	plan, _ := s.repo.FindPlan(planID)
	sub := Subscription{
		ID:        s.repo.store.NextSubscriptionID(),
		UserID:    userID,
		PlanID:    planID,
		Status:    "active",
		StartedAt: now,
		ExpiredAt: now.Add(time.Duration(plan.DurationDays) * 24 * time.Hour),
		CreatedAt: now,
		UpdatedAt: now,
	}
	return s.repo.UpsertSubscription(sub), nil
}

func (s *Service) History(userID string) []Payment {
	return s.repo.ListPaymentsByUserID(userID)
}

func (s *Service) CurrentSubscription(userID string) (Subscription, error) {
	return s.repo.FindCurrentSubscription(userID)
}

func (s *Service) Cancel(userID string) (Subscription, error) {
	sub, err := s.repo.FindCurrentSubscription(userID)
	if err != nil {
		return Subscription{}, err
	}
	now := time.Now()
	sub.Status = "canceled"
	sub.CanceledAt = &now
	sub.UpdatedAt = now
	return s.repo.UpsertSubscription(sub), nil
}
