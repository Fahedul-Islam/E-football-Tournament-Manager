package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type rateLimiterStore struct {
	mu       sync.Mutex
	limiters map[string]*ipLimiter
}

var store = &rateLimiterStore{
	limiters: make(map[string]*ipLimiter),
}

func (s *rateLimiterStore) getLimiter(ip string) *rate.Limiter {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, exists := s.limiters[ip]
	if !exists {
		// Allow 10 requests per second with a burst of 20
		lim := rate.NewLimiter(10, 20)
		s.limiters[ip] = &ipLimiter{limiter: lim, lastSeen: time.Now()}
		return lim
	}
	entry.lastSeen = time.Now()
	return entry.limiter
}

// cleanupStaleEntries removes IPs not seen in the last 3 minutes.
// Call this in a background goroutine to avoid unbounded memory growth.
func (s *rateLimiterStore) cleanupStaleEntries() {
	for {
		time.Sleep(3 * time.Minute)
		s.mu.Lock()
		for ip, entry := range s.limiters {
			if time.Since(entry.lastSeen) > 3*time.Minute {
				delete(s.limiters, ip)
			}
		}
		s.mu.Unlock()
	}
}

func init() {
	go store.cleanupStaleEntries()
}

// RateLimit is a per-IP rate limiting middleware (10 req/s, burst 20).
func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}

		limiter := store.getLimiter(ip)
		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
