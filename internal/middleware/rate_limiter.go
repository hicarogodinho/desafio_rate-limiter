package middleware

import (
	"desafio_rate-limiter/config"
	"desafio_rate-limiter/internal/limiter"
	"desafio_rate-limiter/internal/storage"
	"net/http"
)

func RateLimitMiddleware(store storage.RateLimiterStore, cfg config.Config) func(http.Handler) http.Handler {
	limiterService := limiter.NewLimiter(store, cfg)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			allowed, err := limiterService.AllowRequest(r)
			if err != nil {
				http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
				return
			}

			if !allowed {
				http.Error(w, "Limite de requisições excedido", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
