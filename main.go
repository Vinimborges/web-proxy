package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Captura da URL de destino
	// O trabalho pede o formato: http://localhost:5000/http://www.exemplo.com
	// O r.URL.Path trará "/http://www.exemplo.com". Precisamos remover a barra inicial.
	// Substitua esta linha:
	// targetURL := strings.TrimPrefix(r.URL.Path, "/")

	// Por esta linha:
	targetURL := strings.TrimPrefix(r.RequestURI, "/")

	if targetURL == "" {
		http.Error(w, "URL de destino não fornecida. Exemplo de uso: /http://www.exemplo.com", http.StatusBadRequest)
		return
	}

	// 2. Montando a requisição para o servidor de origem
	// Repassamos o método original (GET, POST) e o corpo da requisição do cliente
	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao criar requisição: %v", err), http.StatusInternalServerError)
		return
	}

	// Copiamos os headers do cliente original para a nova requisição (opcional, mas recomendado)
	for name, values := range r.Header {
		for _, value := range values {
			req.Header.Add(name, value)
		}
	}

	// 3. Executando a requisição (Repasse)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao acessar o servidor de origem: %v", err), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// 4. Devolvendo a resposta ao cliente
	// Primeiro, copiamos os headers da resposta do servidor de origem
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	// Depois, repassamos o HTTP Status Code (ex: 200 OK, 404 Not Found)
	w.WriteHeader(resp.StatusCode)

	// Por fim, copiamos o corpo da resposta (o HTML real da página) de volta para o cliente
	io.Copy(w, resp.Body)

	// Aqui (ou no início do handler) será um bom lugar para inserir a lógica de log no futuro
	log.Printf("Acesso permitido: %s", targetURL)
}

func main() {
	// 1. Remova (ou comente) esta linha:
	// http.HandleFunc("/", proxyHandler)

	porta := ":5000"
	log.Printf("Servidor proxy rodando na porta %s...", porta)

	// 2. Passe a sua função diretamente no ListenAndServe, em vez de 'nil':
	if err := http.ListenAndServe(porta, http.HandlerFunc(proxyHandler)); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}

}
