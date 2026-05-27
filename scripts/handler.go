package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	targetURL := strings.TrimPrefix(r.RequestURI, "/")

	//Verifica se a URL está vazia
	if targetURL == "" {
		http.Error(w, "URL de destino não fornecida. Exemplo: /http://www.example.com", http.StatusBadRequest)
		return
	}

	// Parsing da URL para extrair o Host
	parsedURL, err := url.Parse(targetURL)
	if err != nil || parsedURL.Host == "" {
		http.Error(w, "URL de destino inválida", http.StatusBadRequest)
		return
	}

	// Verificação de Bloqueio
	configMutex.RLock() //Mutex apenas para leitura
	isBlocked := blockedDomains[parsedURL.Host]
	configMutex.RUnlock()

	//  Se a URL estiver na lista de bloqueados, Retorna um HTML personalizado
	if isBlocked {
		logAccess(targetURL, "bloqueado")
		http.ServeFile(w, r, "html/blocked.html")
		return
	}

	// Requisição para a URL desejada
	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao criar requisição: %v", err), http.StatusInternalServerError)
		return
	}

	// Copiar headers, mas tratar casos especiais
	for name, values := range r.Header {
		// ALguns elementos dos headers são remvidos para evitar conflitos ou erros
		if strings.EqualFold(name, "Host") || strings.EqualFold(name, "Accept-Encoding") ||
			strings.EqualFold(name, "If-Modified-Since") || strings.EqualFold(name, "If-None-Match") {
			continue
		}
		for _, value := range values {
			req.Header.Add(name, value)
		}
	}

	// Esse header serve para evitar Gzip como resposta, queremos apenas texto
	req.Header.Set("Accept-Encoding", "identity")

	client := &http.Client{}
	resp, err := client.Do(req) // Faz a requisição
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao acessar o servidor de origem: %v", err), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Preparação para resposta ao cliente
	for name, values := range resp.Header {
		// Pula alguns headers
		if strings.EqualFold(name, "Content-Length") ||
			strings.EqualFold(name, "Content-Encoding") ||
			strings.EqualFold(name, "Transfer-Encoding") {
			continue
		}
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	contentType := strings.ToLower(resp.Header.Get("Content-Type"))

	if strings.Contains(contentType, "text/html") {
		// Filtro de Conteúdo do HTML
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Erro ao ler corpo da resposta", http.StatusInternalServerError)
			return
		}

		body := string(bodyBytes)
		isModified := false

		configMutex.RLock()
		for _, filter := range regexFilters {
			if filter.regex.MatchString(body) {
				body = filter.regex.ReplaceAllString(body, filter.replacement)
				isModified = true
			}
		}
		configMutex.RUnlock()

		if isModified {
			// Se entrou aqui, o HTML foi filtrado
			logAccess(targetURL, "filtrado")
			w.Header().Del("Content-Length")
		} else {
			// Se entrou aqui, o HTML não foi filtrado
			logAccess(targetURL, "permitido")
		}

		w.WriteHeader(resp.StatusCode)
		io.WriteString(w, body)
	} else {
		// Repasse Direto (não HTML - imagens, scripts, etc)
		logAccess(targetURL, "permitido")
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}
