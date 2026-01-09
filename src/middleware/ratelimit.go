package middleware

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"quizzler/cache"
)

var RateLimitWindow = time.Minute

func RateLimit(maxRequests int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getClientIP(r)
			key := fmt.Sprintf("ratelimit:%s:%s", r.URL.Path, ip)

			// Check if rate limiting is available (Redis connected)
			if !cache.IsAvailable() {
				// If Redis is not available, allow the request
				next.ServeHTTP(w, r)
				return
			}

			// Get current count
			count, err := cache.Increment(key, RateLimitWindow)
			if err != nil {
				// If there's an error, allow the request
				next.ServeHTTP(w, r)
				return
			}

			// Set rate limit headers
			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", maxRequests))
			w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", max(0, maxRequests-int(count))))

			if int(count) > maxRequests {
				w.Header().Set("Retry-After", fmt.Sprintf("%d", int(RateLimitWindow.Seconds())))
				http.Error(w, `{"error": "Too many requests. Please try again later."}`, http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// getClientIP extracts the client IP from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (for proxies/load balancers)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
