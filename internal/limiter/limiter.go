package limiter

import (
	"desafio_rate-limiter/config"
	"desafio_rate-limiter/internal/storage"
	"net/http"
	"strings"
)

type Limiter struct {
	Store  storage.RateLimiterStore
	Config config.Config
}

func NewLimiter(store storage.RateLimiterStore, cfg config.Config) *Limiter {
	return &Limiter{
		Store:  store,
		Config: cfg,
	}
}

func (l *Limiter) getKeyAndLimit(r *http.Request) (string, int) {
	token := r.Header.Get("API_KEY")
	if token != "" {
		return "token:" + token, l.Config.RateLimitTokenDefault
	}

	ip := getIP(r)
	return "ip:" + ip, l.Config.RateLimitIP
}

func (l *Limiter) AllowRequest(r *http.Request) (bool, error) {
	key, limit := l.getKeyAndLimit(r)

	count, err := l.Store.Increment(key)
	if err != nil {
		return false, err
	}

	if count == 1 {
		// Primeira requisição, vai definir a expiração
		_ = l.Store.SetExpiration(key, l.Config.BlockDuration)
	}

	if count > limit {
		return false, nil
	}

	return true, nil
}

func getIP(r *http.Request) string {
	fowarded := r.Header.Get("X-Forwarded-For")
	if fowarded != "" {
		parts := strings.Split(fowarded, ",")
		return strings.TrimSpace(parts[0])
	}
	ip := r.RemoteAddr
	if colon := strings.LastIndex(ip, ":"); colon != -1 {
		ip = ip[:colon]
	}
	return ip
}
