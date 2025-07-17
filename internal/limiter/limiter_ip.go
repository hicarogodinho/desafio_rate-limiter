package limiter

import (
	"net/http"
	"strings"
)

// GetIPKey extrai o IP do cliente e retonra a chave de limitação
func GetIPKey(r *http.Request) string {
	ip := extractClientIP(r)
	return "ip:" + ip
}

// extractClienteIP tenta obter o IP real do cliente
func extractClientIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}

	ip := r.RemoteAddr
	if colon := strings.LastIndex(ip, ":"); colon != -1 {
		ip = ip[:colon]
	}

	return ip
}
