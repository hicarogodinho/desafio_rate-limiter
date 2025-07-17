package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"desafio_rate-limiter/api"
	"desafio_rate-limiter/config"
	"desafio_rate-limiter/internal/middleware"
	"desafio_rate-limiter/internal/storage/redis"
)

func main() {
	// Carrega variáveis de ambiente
	err := godotenv.Load()
	if err != nil {
		log.Println("Arquivo .env não encontrado")
	}

	cfg := config.Load()

	// Inicializa o Redis
	redisStore, err := redis.NewRedisStore(cfg)
	if err != nil {
		log.Fatalf("Erro ao conectar ao Redis: %v", err)
	}

	// Cria o middleware de rate limiter
	rateLimiter := middleware.RateLimitMiddleware(redisStore, cfg)

	// Define o handler com middleware
	mux := http.NewServeMux()
	mux.HandleFunc("/", api.HomeHandler)

	// Aplica o middleware
	handler := rateLimiter(mux)

	// Inicia o servidor
	port := ":8080"
	fmt.Printf("Servidor rodando na porta %s\n", port)
	if err := http.ListenAndServe(port, handler); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
