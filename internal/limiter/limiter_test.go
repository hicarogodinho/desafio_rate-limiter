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
	blocks      map[string]time.Time
}

func NewMockStore() *mockStore {
	return &mockStore{
		counters:    make(map[string]int),
		expirations: make(map[string]time.Duration),
		blocks:      make(map[string]time.Time),
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

// IsBlocked verifica se a chave está bloqueada no mock store
func (m *mockStore) IsBlocked(key string) (bool, error) {
	if unlockTime, ok := m.blocks[key]; ok {
		// Se o tempo atual for antes do tempo de desbloqueio, ainda está bloqueado
		return time.Now().Before(unlockTime), nil
	}
	return false, nil // Não está na lista de bloqueios, então não está bloqueado
}

// Block bloqueia a chave no mock store por uma duração específica
func (m *mockStore) Block(key string, duration time.Duration) error {
	// Calcula o tempo em que o bloqueio deve expirar
	m.blocks[key] = time.Now().Add(duration)
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
