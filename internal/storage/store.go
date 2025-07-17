package storage

import "time"

// Define uma interface para mecanismos de persistência de rate limiting
type RateLimiterStore interface {
	// Incrementa o contador associado à chave e retorna o novo valor
	Increment(key string) (int, error)

	// Define o tempo de expiração para a chave
	SetExpiration(key string, exxpiration time.Duration) error

	// Obtém o valor atual do contador
	Get(key string) (int, error)

	// Reseta o contador da chave
	Reset(key string) error
}
