package test

import (
	"desafio_rate-limiter/api"
	"desafio_rate-limiter/config"
	"desafio_rate-limiter/internal/middleware"
	"desafio_rate-limiter/internal/storage/redis"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiterIntegration_IP(t *testing.T) {
	cfg := config.Load()
	store, err := redis.NewRedisStore(cfg)
	if err != nil {
		t.Fatalf("failed to create redis store: %v", err)
	}

	handler := middleware.RateLimitMiddleware(store, cfg)(http.HandlerFunc(api.HomeHandler))

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "192.168.1.100:1234"

	for i := 1; i <= cfg.RateLimitIP+2; i++ {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if i <= cfg.RateLimitIP && rr.Code != http.StatusOK {
			t.Errorf("expected status OK for request %d, got %d", i, rr.Code)
		}

		if i > cfg.RateLimitIP && rr.Code != http.StatusTooManyRequests {
			t.Errorf("expected status TooManyRequests for request %d, got %d", i, rr.Code)
		}
	}

	// Aguarda expiração para garantir que o IP seja bloqueado
	time.Sleep(cfg.BlockDuration + time.Second)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status OK after expiration, got %d", rr.Code)
	}
}

func TestRateLimiterIntegration_Token(t *testing.T) {
	cfg := config.Load()
	store, err := redis.NewRedisStore(cfg)
	if err != nil {
		t.Fatalf("Erro ao conectar ao Redis: %v", err)
	}

	handler := middleware.RateLimitMiddleware(store, cfg)(http.HandlerFunc(api.HomeHandler))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("API_KEY", "test-token-123")

	for i := 1; i <= cfg.RateLimitTokenDefault+2; i++ {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if i <= cfg.RateLimitTokenDefault && rr.Code != http.StatusOK {
			t.Errorf("Esperado 200 OK na requisição %d, mas recebeu %d", i, rr.Code)
		}

		if i > cfg.RateLimitTokenDefault && rr.Code != http.StatusTooManyRequests {
			t.Errorf("Esperado 429 Too Many Requests na requisição %d, mas recebeu %d", i, rr.Code)
		}
	}
}
