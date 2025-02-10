package main

import (
	"fmt"
	"net/http"

	"github.com/alextavella/rate-limiter/pkg"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	rateLimiter := pkg.NewRateLimiter(pkg.RateLimiterConfig{
		RateLimiterValue: 5,
		Keys: map[string]pkg.RateLimiterValue{
			"ABC": 10,
			"DEF": 20,
		},
	})

	svr := http.Server{
		Addr:                         ":8080",
		Handler:                      rateLimiter.Handler(mux),
		DisableGeneralOptionsHandler: false,
	}

	fmt.Println("Server running on port 8080")
	if err := svr.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			panic(err)
		}
	}
}
