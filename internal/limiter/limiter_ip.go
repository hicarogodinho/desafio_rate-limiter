package limiter

import (
	"net/http"
	"strings"
)

// GetIPKey extrai o IP real do cliente e retorna a chave de limitação para o Redis.
// Ela prioriza o cabeçalho X-Forwarded-For para cenários com proxies.
func GetIPKey(r *http.Request) string {
	ip := extractIP(r) // <--- Agora, GetIPKey chama a função auxiliar mais robusta
	return "ip:" + ip
}

// extractIP tenta obter o IP real do cliente, considerando X-Forwarded-For
// e removendo a porta de RemoteAddr.
// Esta é a sua antiga 'extractClientIP', renomeada e tornada privada (minúscula).
func extractIP(r *http.Request) string {
	// Primeiro, verifica o cabeçalho X-Forwarded-For
	// Isso é crucial se sua aplicação estiver atrás de um proxy ou load balancer.
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		// X-Forwarded-For pode conter uma lista de IPs (ex: client, proxy1, proxy2).
		// O primeiro IP na lista é geralmente o do cliente original.
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}

	// Se X-Forwarded-For não estiver presente, usa r.RemoteAddr.
	// r.RemoteAddr geralmente inclui a porta (ex: "192.168.1.1:12345").
	ip := r.RemoteAddr
	// Encontra a última ocorrência de ':' para separar o IP da porta.
	if colon := strings.LastIndex(ip, ":"); colon != -1 {
		ip = ip[:colon] // Pega a parte do IP antes do último ':'
	}

	return strings.TrimSpace(ip) // Retorna o IP sem espaços em branco
}
