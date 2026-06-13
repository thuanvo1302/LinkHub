package middleware

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"linkhub/backend/internal/pkg/response"
)

type limitBucket struct {
	Count       int
	WindowStart time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	buckets map[string]limitBucket
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		buckets: map[string]limitBucket{},
	}
}

func (l *RateLimiter) Allow(key string, limit int, window time.Duration) bool {
	now := time.Now()
	l.mu.Lock()
	defer l.mu.Unlock()

	bucket, ok := l.buckets[key]
	if !ok || now.Sub(bucket.WindowStart) >= window {
		l.buckets[key] = limitBucket{
			Count:       1,
			WindowStart: now,
		}
		return true
	}

	if bucket.Count >= limit {
		return false
	}

	bucket.Count++
	l.buckets[key] = bucket
	return true
}

func (l *RateLimiter) IP(limit int, window time.Duration, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := clientIP(r)
		if !l.Allow("ip:"+ip+":"+r.URL.Path, limit, window) {
			response.Error(w, http.StatusTooManyRequests, "RATE_LIMITED", "Too many requests", nil)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func (l *RateLimiter) User(limit int, window time.Duration, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := UserIDFromContext(r.Context())
		if strings.TrimSpace(userID) == "" {
			response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing user context", nil)
			return
		}
		if !l.Allow("user:"+userID+":"+r.URL.Path, limit, window) {
			response.Error(w, http.StatusTooManyRequests, "RATE_LIMITED", "Too many requests", nil)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func clientIP(r *http.Request) string {
	if forwarded := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); forwarded != "" {
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}
	if realIP := strings.TrimSpace(r.Header.Get("X-Real-IP")); realIP != "" {
		return realIP
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && host != "" {
		return host
	}
	return r.RemoteAddr
}
