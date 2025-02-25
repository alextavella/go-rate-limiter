package pkg

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type (
	RateLimiter struct {
		config  RateLimiterConfig
		clients map[string]*RateLimiterClient
	}
	RateLimiterClient struct {
		RequestTime time.Time
		Limiter     *rate.Limiter
	}
	RateLimiterConfig struct {
		ApiKeyHeader string
		ApiKeys      map[string]int
		MaxRequests  int
		LockedTime   time.Duration
	}
)

func NewRateLimiter(c RateLimiterConfig) RateLimiter {
	return RateLimiter{
		config:  c,
		clients: make(map[string]*RateLimiterClient),
	}
}

func (rl *RateLimiter) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientKey, limitValue := rl.getRequestConfig(r)

		var client, ok = rl.clients[clientKey]
		if !ok {
			client = &RateLimiterClient{
				RequestTime: time.Now(),
				Limiter:     rate.NewLimiter(rate.Limit(limitValue), int(limitValue)),
			}
			rl.clients[clientKey] = client
		}

		isFreeze := time.Since(client.RequestTime) < rl.config.LockedTime

		if !client.Limiter.Allow() && isFreeze {
			w.WriteHeader(http.StatusTooManyRequests)
			fmt.Fprintf(w, "429 - Too Many Requests")
			return
		}

		client.RequestTime = time.Now()
		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) getRequestConfig(r *http.Request) (string, rate.Limit) {
	maxReq := rate.Limit(rl.config.MaxRequests)
	if token := r.Header.Get(rl.config.ApiKeyHeader); token != "" {
		// Se houver token, verifica se ele está na lista de tokens
		if max, ok := rl.config.ApiKeys[token]; ok {
			return token, rate.Limit(max)
		}
		// Se não estiver, usa o valor padrão
		return token, maxReq
	}
	// Se não houver token, usa o IP
	return getIP(r), maxReq
}

func getIP(r *http.Request) string {
	// Obtém o IP diretamente do RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "Erro ao obter IP"
	}
	// Se estiver rodando localmente, pode ser [::1], então converter para 127.0.0.1
	if ip == "::1" {
		ip = "127.0.0.1"
	}
	return ip
}
