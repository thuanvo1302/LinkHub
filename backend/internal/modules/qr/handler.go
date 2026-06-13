package qr

import (
	"fmt"
	"net/http"
	"strings"

	"linkhub/backend/internal/config"
	"linkhub/backend/internal/modules/profile"
	"linkhub/backend/internal/modules/shortener"
)

type Handler struct {
	cfg         config.Config
	profileRepo *profile.Repository
	shortRepo   *shortener.Repository
}

func NewHandler(cfg config.Config, profileRepo *profile.Repository, shortRepo *shortener.Repository) *Handler {
	return &Handler{cfg: cfg, profileRepo: profileRepo, shortRepo: shortRepo}
}

func (h *Handler) ProfileSVG(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimSpace(r.URL.Query().Get("username"))
	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}
	target := fmt.Sprintf("%s/%s", h.cfg.FrontendURL, username)
	writeSVG(w, target, "Profile QR")
}

func (h *Handler) ShortLinkSVG(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/qr/short-links/")
	item, err := h.shortRepo.FindByID(id)
	if err != nil {
		http.Error(w, "short link not found", http.StatusNotFound)
		return
	}
	target := fmt.Sprintf("%s/%s", h.cfg.AppBaseURL, item.Code)
	writeSVG(w, target, "Short Link QR")
}

func writeSVG(w http.ResponseWriter, value string, title string) {
	escaped := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;").Replace(value)
	w.Header().Set("Content-Type", "image/svg+xml")
	_, _ = w.Write([]byte(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="360" height="360" viewBox="0 0 360 360">
  <rect width="360" height="360" rx="28" fill="#f4efe7"/>
  <rect x="24" y="24" width="312" height="312" rx="22" fill="#17181f"/>
  <rect x="48" y="48" width="264" height="264" rx="18" fill="#fff"/>
  <text x="50%%" y="82" text-anchor="middle" font-family="Arial" font-size="18" fill="#17181f">%s</text>
  <text x="50%%" y="180" text-anchor="middle" font-family="Arial" font-size="16" fill="#17181f">Scan or copy</text>
  <foreignObject x="62" y="196" width="236" height="90">
    <div xmlns="http://www.w3.org/1999/xhtml" style="font-family:Arial,sans-serif;font-size:13px;color:#17181f;text-align:center;line-height:1.5;word-break:break-word;">%s</div>
  </foreignObject>
  <rect x="102" y="104" width="32" height="32" fill="#f26a4b"/>
  <rect x="226" y="104" width="32" height="32" fill="#4aa3a2"/>
  <rect x="102" y="228" width="32" height="32" fill="#4aa3a2"/>
  <rect x="226" y="228" width="32" height="32" fill="#f26a4b"/>
  <rect x="156" y="132" width="48" height="96" fill="#17181f"/>
  <rect x="132" y="156" width="96" height="48" fill="#17181f"/>
</svg>`, title, escaped)))
}
