package main

import (
	"log"
	"net/http"
)

func main() {
	loadConfig()

	porta := ":5000"
	log.Printf("Servidor proxy rodando na porta %s...", porta)

	if err := http.ListenAndServe(porta, http.HandlerFunc(proxyHandler)); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}

}
