package limiter

import (
	"desafio_rate-limiter/config"
	"net/http"
	"testing"
	"time"
)

type mockStore struct {
	counters    map[string]int
	expirations map[string]time.Duration
}

func NewMockStore() *mockStore {
	return &mockStore{
		counters:    make(map[string]int),
		expirations: make(map[string]time.Duration),
	}
}

func (m *mockStore) Increment(key string) (int, error) {
	m.counters[key]++
	return m.counters[key], nil
}

func (m *mockStore) SetExpiration(key string, expiration time.Duration) error {
	m.expirations[key] = expiration
	return nil
}

func (m *mockStore) Get(key string) (int, error) {
	return m.counters[key], nil
}

func (m *mockStore) Reset(key string) error {
	delete(m.counters, key)
	delete(m.expirations, key)
	return nil
}

func TestLimiterByIP(t *testing.T) {
	store := NewMockStore()
	cfg := config.Config{
		RateLimitIP:           3,
		RateLimitTokenDefault: 5,
		BlockDuration:         5 * time.Minute,
	}

	limiter := NewLimiter(store, cfg)

	req, _ := http.NewRequest("GET", "/", nil)
	req.RemoteAddr = "192.168.0.1:1234"

	for i := 1; i <= 5; i++ {
		allowed, err := limiter.AllowRequest(req)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		for i <= cfg.RateLimitIP && !allowed {
			t.Errorf("Requisição %d deveria ser permitida", i)
		}

		if i > cfg.RateLimitIP && allowed {
			t.Errorf("Requisição %d deveria ser bloqueada", i)
		}
	}
}

func TestLimiterByToken(t *testing.T) {
	store := NewMockStore()
	cfg := config.Config{
		RateLimitIP:           3,
		RateLimitTokenDefault: 4,
		BlockDuration:         5 * time.Minute,
	}

	limiter := NewLimiter(store, cfg)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("API_KEY", "token123")

	for i := 1; i <= 5; i++ {
		allowed, err := limiter.AllowRequest(req)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if i <= cfg.RateLimitTokenDefault && !allowed {
			t.Errorf("Requisição %d deveria ser permitida", i)
		}

		if i > cfg.RateLimitTokenDefault && allowed {
			t.Errorf("Requisição %d deveria ser bloqueada", i)
		}
	}
}
