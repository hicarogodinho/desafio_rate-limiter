package redis

import (
	"context"
	"desafio_rate-limiter/config"
	"desafio_rate-limiter/internal/storage"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisStore(cfg config.Config) (storage.RateLimiterStore, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.RedisAddr,
		Password:     cfg.RedisPassword,
		DB:           cfg.RedisDB,
		DialTimeout:  10 * time.Second, // Tempo máximo para estabelecer a conexão TCP
		ReadTimeout:  5 * time.Second,  // Tempo máximo para ler a resposta do servidor
		WriteTimeout: 5 * time.Second,  // Tempo máximo para escrever a requisição no servidor
		PoolTimeout:  10 * time.Second, // Tempo máximo para pegar uma conexão do pool
	})

	ctx := context.Background()

	// Implementa a lógica de retry para a conexão Redis
	maxRetries := 5               // Número máximo de tentativas
	retryDelay := 2 * time.Second // Atraso entre as tentativas

	for i := 0; i < maxRetries; i++ {
		// Adicione um timeout para cada tentativa de ping
		pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second) // Ping com timeout de 3 segundos
		err := rdb.Ping(pingCtx).Err()
		cancel() // Libere o recurso do contexto imediatamente

		if err == nil {
			fmt.Println("Conexão com Redis estabelecida com sucesso!")
			return &RedisStore{
				client: rdb,
				ctx:    ctx,
			}, nil
		}

		fmt.Printf("Tentativa %d/%d: Erro ao conectar ao Redis: %v. Tentando novamente em %v...\n", i+1, maxRetries, err, retryDelay)
		time.Sleep(retryDelay) // Espera antes de tentar novamente
	}

	// Se todas as tentativas falharem
	return nil, fmt.Errorf("falha em conectar ao Redis após %d tentativas", maxRetries)
	// if err := rdb.Ping(ctx).Err(); err != nil {
	// 	return nil, fmt.Errorf("erro ao conectar ao Redis: %w", err)
	// }

	// return &RedisStore{
	// 	client: rdb,
	// 	ctx:    ctx,
	// }, nil
}

func (r *RedisStore) Increment(key string) (int, error) {
	count, err := r.client.Incr(r.ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *RedisStore) SetExpiration(key string, expiration time.Duration) error {
	return r.client.Expire(r.ctx, key, expiration).Err()
}

func (r *RedisStore) Get(key string) (int, error) {
	val, err := r.client.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	count, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *RedisStore) Reset(key string) error {
	return r.client.Del(r.ctx, key).Err()
}
