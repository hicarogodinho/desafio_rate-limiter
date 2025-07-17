package limiter

import (
	"net/http"
	"strings"
)

// GetIPKey extrai o IP do cliente e retonra a chave de limitação
func GetTokenKey(r *http.Request) (string, bool) {
	token := r.Header.Get("API_KEY")
	if token == "" {
		return "", false
	}
	return "token:" + strings.TrimSpace(token), true
}
