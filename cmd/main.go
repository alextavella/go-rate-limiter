package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alextavella/rate-limiter/pkg"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	rateLimiter := pkg.NewRateLimiter(pkg.RateLimiterConfig{
		MaxRequests:  5,
		LockedTime:   time.Second * 5,
		ApiKeyHeader: "API_KEY",
		ApiKeys: map[string]int{
			"ABC": 10,
			"DEF": 8,
		},
	})

	svr := http.Server{
		Addr:                         ":8080",
		Handler:                      rateLimiter.Handler(mux),
		DisableGeneralOptionsHandler: false,
	}

	fmt.Println("🚀 Server running on port 8080")
	if err := svr.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			panic(err)
		}
	}
}
