package api

import (
	"fmt"
	"net/http"
)

// Homehandler responde com uma mensagem para tesatr o rate limiter
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Requisição aceita com sucesso!")
}
