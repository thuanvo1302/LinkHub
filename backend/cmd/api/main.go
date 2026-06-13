package main

import (
	"log"
	"net/http"
	"time"

	"linkhub/backend/internal/config"
	"linkhub/backend/internal/database"
	"linkhub/backend/internal/middleware"
	"linkhub/backend/internal/modules/auth"
	"linkhub/backend/internal/modules/billing"
	"linkhub/backend/internal/modules/link"
	"linkhub/backend/internal/modules/profile"
	"linkhub/backend/internal/modules/qr"
	"linkhub/backend/internal/modules/shortener"
	"linkhub/backend/internal/pkg/response"
)

func main() {
	cfg := config.Load()
	store := database.NewStore()
	rateLimiter := middleware.NewRateLimiter()

	authRepo := auth.NewRepository(store)
	authService := auth.NewService(authRepo, cfg)
	authHandler := auth.NewHandler(authService)

	profileRepo := profile.NewRepository(store)
	profileService := profile.NewService(profileRepo)

	linkRepo := link.NewRepository(store)
	linkService := link.NewService(linkRepo, func(userID string) (string, error) {
		item, err := profileRepo.FindByUserID(userID)
		if err != nil {
			return "", err
		}
		return item.ID, nil
	})
	linkHandler := link.NewHandler(linkService)
	profileHandler := profile.NewHandler(profileService, linkRepo)

	shortRepo := shortener.NewRepository(store)
	shortService := shortener.NewService(shortRepo, cfg)
	shortHandler := shortener.NewHandler(shortService)

	billingRepo := billing.NewRepository(store)
	billingService := billing.NewService(billingRepo, cfg)
	billingHandler := billing.NewHandler(billingService)

	qrHandler := qr.NewHandler(cfg, profileRepo, shortRepo)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		response.OK(w, map[string]any{
			"service": cfg.AppName,
			"status":  "healthy",
		}, "OK")
	})

	mux.HandleFunc("POST /api/v1/auth/register", rateLimiter.IP(5, time.Minute, authHandler.Register))
	mux.HandleFunc("POST /api/v1/auth/login", rateLimiter.IP(5, time.Minute, authHandler.Login))
	mux.HandleFunc("POST /api/v1/auth/refresh", rateLimiter.IP(10, time.Minute, authHandler.Refresh))
	mux.HandleFunc("POST /api/v1/auth/logout", authHandler.Logout)
	mux.Handle("GET /api/v1/me", middleware.Auth(cfg, authHandler.Me))

	mux.Handle("GET /api/v1/profiles/me", middleware.Auth(cfg, profileHandler.GetMe))
	mux.Handle("PUT /api/v1/profiles/me", middleware.Auth(cfg, profileHandler.UpdateMe))
	mux.HandleFunc("GET /api/v1/public/profiles/", rateLimiter.IP(300, time.Minute, profileHandler.GetPublic))

	mux.Handle("GET /api/v1/profile-links", middleware.Auth(cfg, linkHandler.List))
	mux.Handle("POST /api/v1/profile-links", middleware.Auth(cfg, linkHandler.Create))
	mux.Handle("PATCH /api/v1/profile-links/reorder", middleware.Auth(cfg, linkHandler.Reorder))
	mux.Handle("PUT /api/v1/profile-links/", middleware.Auth(cfg, linkHandler.Update))
	mux.Handle("DELETE /api/v1/profile-links/", middleware.Auth(cfg, linkHandler.Delete))

	mux.Handle("GET /api/v1/short-links", middleware.Auth(cfg, shortHandler.List))
	mux.Handle("POST /api/v1/short-links", middleware.Auth(cfg, http.HandlerFunc(rateLimiter.User(30, time.Hour, shortHandler.Create))))
	mux.Handle("GET /api/v1/short-links/", middleware.Auth(cfg, shortHandler.Get))
	mux.Handle("PUT /api/v1/short-links/", middleware.Auth(cfg, shortHandler.Update))
	mux.Handle("DELETE /api/v1/short-links/", middleware.Auth(cfg, shortHandler.Delete))
	mux.Handle("GET /api/v1/analytics/overview", middleware.Auth(cfg, shortHandler.Overview))
	mux.Handle("GET /api/v1/analytics/short-links/", middleware.Auth(cfg, shortHandler.LinkAnalytics))

	mux.HandleFunc("GET /api/v1/plans", billingHandler.Plans)
	mux.Handle("POST /api/v1/payments/create-checkout", middleware.Auth(cfg, billingHandler.CreateCheckout))
	mux.Handle("POST /api/v1/payments/mock-success", middleware.Auth(cfg, billingHandler.MockSuccess))
	mux.Handle("GET /api/v1/billing/history", middleware.Auth(cfg, billingHandler.History))
	mux.Handle("GET /api/v1/subscription/current", middleware.Auth(cfg, billingHandler.Current))
	mux.Handle("POST /api/v1/subscription/cancel", middleware.Auth(cfg, billingHandler.Cancel))

	mux.HandleFunc("GET /api/v1/qr/profile", rateLimiter.IP(120, time.Minute, qrHandler.ProfileSVG))
	mux.HandleFunc("GET /api/v1/qr/short-links/", rateLimiter.IP(120, time.Minute, qrHandler.ShortLinkSVG))

	// Catch-all redirect route should stay last.
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		response.OK(w, map[string]any{
			"name":        cfg.AppName,
			"description": "Link-in-bio and URL shortener MVP API",
		}, "OK")
	})
	mux.HandleFunc("GET /{code}", rateLimiter.IP(300, time.Minute, shortHandler.Redirect))

	handler := middleware.CORS(cfg, middleware.SecurityHeaders(middleware.Logging(mux)))

	log.Printf("starting %s on :%s", cfg.AppName, cfg.AppPort)
	if err := http.ListenAndServe(":"+cfg.AppPort, handler); err != nil {
		log.Fatal(err)
	}
}
