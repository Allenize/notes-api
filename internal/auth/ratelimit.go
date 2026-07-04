package auth

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type bucket struct {
	tokens   float64
	lastSeen time.Time
}

type RateLimiter struct {
	mu        sync.Mutex
	buckets   map[string]*bucket
	rate      float64
	burst     float64
	lastSweep time.Time
}

func NewRateLimiter(requestsPerMinute int, burst int) *RateLimiter {
	return &RateLimiter{
		buckets:   make(map[string]*bucket),
		rate:      float64(requestsPerMinute) / 60.0,
		burst:     float64(burst),
		lastSweep: time.Now(),
	}
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func (rl *RateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	if now.Sub(rl.lastSweep) > 10*time.Minute {
		for k, b := range rl.buckets {
			if now.Sub(b.lastSeen) > 10*time.Minute {
				delete(rl.buckets, k)
			}
		}
		rl.lastSweep = now
	}

	b, ok := rl.buckets[key]
	if !ok {
		b = &bucket{tokens: rl.burst, lastSeen: now}
		rl.buckets[key] = b
	}

	elapsed := now.Sub(b.lastSeen).Seconds()
	b.tokens += elapsed * rl.rate
	if b.tokens > rl.burst {
		b.tokens = rl.burst
	}
	b.lastSeen = now

	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rl.allow(clientIP(r)) {
			w.Header().Set("Retry-After", "5")
			http.Error(w, `{"error":"rate limit exceeded, slow down"}`, http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
