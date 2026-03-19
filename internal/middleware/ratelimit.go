package middleware

import (
	"net"
	"net/http"
	"strings"
	"sync"

	"service-courier/pkg/ratelimit"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	rateLimitExceededTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "rate_limit_exceeded_total",
			Help: "Total number of rate limit exceeded events",
		},
	)
)

type RateLimiter struct {
	globalLimiter *ratelimit.TokenBucket
	ipLimiters    map[string]*ratelimit.TokenBucket
	ipRPS         int
	mu            sync.RWMutex
}

func NewRateLimiter(globalRPS, ipRPS int) *RateLimiter {
	return &RateLimiter{
		globalLimiter: ratelimit.NewTokenBucket(globalRPS, globalRPS),
		ipLimiters:    make(map[string]*ratelimit.TokenBucket),
		ipRPS:         ipRPS,
	}
}

func (rl *RateLimiter) allowIP(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.ipLimiters[ip]
	if !exists {
		limiter = ratelimit.NewTokenBucket(rl.ipRPS, rl.ipRPS)
		rl.ipLimiters[ip] = limiter
	}

	return limiter.Allow()
}

func getClientIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		ips := strings.Split(forwarded, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return host
}

func RateLimitMiddleware(limiter *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := getClientIP(r)

			if !limiter.globalLimiter.Allow() {
				rateLimitExceededTotal.Inc()
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			if !limiter.allowIP(clientIP) {
				rateLimitExceededTotal.Inc()
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
