package middleware

import (
	"desafio_rate-limiter/config"
	"desafio_rate-limiter/internal/limiter"
	"desafio_rate-limiter/internal/storage"
	"log"
	"net/http"
)

func RateLimitMiddleware(store storage.RateLimiterStore, cfg config.Config) func(http.Handler) http.Handler {
	// limiterService := limiter.NewLimiter(store, cfg)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var key string
			var limit int
			//var isToken bool // Não precisamos mais de 'isToken' explicitamente se a lógica for bem definida

			// Tenta obter a chave e o limite baseados no token
			tokenFromHeader := r.Header.Get("API_KEY")
			if tokenFromHeader != "" {
				// Usa o token exato do cabeçalho como a chave para buscar o limite
				// e também para formar a chave do Redis
				if specificLimit, found := cfg.TokenLimits[tokenFromHeader]; found {
					key = "token:" + tokenFromHeader // Chave no Redis para este token específico
					limit = specificLimit
				} else {
					// Token existe no header, mas não é um token específico configurado, usa o padrão
					key = "token:" + tokenFromHeader // Ainda usa o token real como parte da chave
					limit = cfg.RateLimitTokenDefault
				}
			} else {
				// Se não houver token, usa o limite baseado no IP
				key = limiter.GetIPKey(r)
				limit = cfg.RateLimitIP
			}

			// Lógica de Rate Limiting (reusável para IP e Token)
			currentCount, err := store.Increment(key)
			if err != nil {
				log.Printf("Erro ao incrementar contador para %s: %v", key, err)
				http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
				return
			}

			// Define a expiração se for a primeira requisição ou se não estiver definida
			if currentCount == 1 {
				if err := store.SetExpiration(key, cfg.BlockDuration); err != nil {
					log.Printf("Erro ao definir expiração para %s: %v", key, err)
					// Não é um erro fatal, o limite ainda pode funcionar
				}
			}

			// Verifica se o IP/Token está bloqueado
			isBlocked, err := store.IsBlocked(key)
			if err != nil {
				log.Printf("Erro ao verificar bloqueio para %s: %v", key, err)
				http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
				return
			}

			if isBlocked {
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte("Você foi bloqueado devido a muitas requisições."))
				return
			}

			// Verifica se o limite foi excedido
			if currentCount > limit {
				// Se excedeu, bloqueia o IP/Token
				if err := store.Block(key, cfg.BlockDuration); err != nil {
					log.Printf("Erro ao bloquear %s: %v", key, err)
				}
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte("you have reached the maximum number of requests or actions allowed within a certain time frame"))
				return
			}

			// Se não foi bloqueado nem excedeu o limite, continua para o próximo handler
			next.ServeHTTP(w, r)
		})
	}
}
