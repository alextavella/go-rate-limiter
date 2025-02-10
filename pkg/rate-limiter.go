package pkg

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

type (
	RateLimiter struct {
		config  RateLimiterConfig
		clients map[string]*rate.Limiter
	}
	RateLimiterValue  rate.Limit
	RateLimiterConfig struct {
		RateLimiterValue
		Keys map[string]RateLimiterValue
	}
)

func NewRateLimiter(c RateLimiterConfig) RateLimiter {
	return RateLimiter{
		config:  c,
		clients: make(map[string]*rate.Limiter),
	}
}

func (rl *RateLimiter) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()
		clientKey, limitValue := rl.getRequestConfig(r)

		var limiter, ok = rl.clients[clientKey]
		if !ok {
			limiter = rate.NewLimiter(rate.Limit(limitValue), 1)
		}

		rl.clients[clientKey] = limiter

		if !limiter.Allow() {
			w.WriteHeader(http.StatusTooManyRequests)
			fmt.Fprintf(w, "429 - Too Many Requests")
			return
		}

		next.ServeHTTP(w, r)
		fmt.Println(time.Since(begin), "|", clientKey, "|", limitValue)
	})
}

func (rl *RateLimiter) getRequestConfig(r *http.Request) (string, RateLimiterValue) {
	dcf := rl.config.RateLimiterValue
	if token := getToken(r); token != "" {
		if cf, ok := rl.config.Keys[token]; ok {
			return token, cf
		}
		return token, dcf
	}
	return getIP(r), dcf
}

func getToken(r *http.Request) string {
	return r.Header.Get("API_KEY")
}

func getIP(r *http.Request) string {
	// Verifica o cabeçalho X-Forwarded-For
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0]) // Pega o primeiro IP da lista
	}

	// Verifica o cabeçalho X-Real-IP
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Usa o RemoteAddr como fallback
	ip := r.RemoteAddr
	// Remove a porta se estiver presente
	if strings.Contains(ip, ":") {
		ip = strings.Split(ip, ":")[0]
	}

	return ip
}
