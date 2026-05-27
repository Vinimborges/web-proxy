package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	URL       string `json:"url"`
	Action    string `json:"action"`
}

var logMutex sync.Mutex

func logAccess(targetURL string, action string) {
	logMutex.Lock()
	defer logMutex.Unlock()

	// Carrega o fuso do Brasil
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		fmt.Println("Erro ao carregar o timezone:", err)
		return
	}

	timeNowBR := time.Now().In(loc) //Pega a hora exata

	entry := LogEntry{
		Timestamp: timeNowBR.Format("02/01/2006 15:04:05"),
		URL:       targetURL,
		Action:    action,
	}

	// Print no terminal para visualização em tempo real
	fmt.Printf("Timestamp: %s | Ação: %-10s | URL: %s\n", entry.Timestamp, entry.Action, entry.URL)

	f, err := os.OpenFile("json/log.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Erro ao abrir arquivo de logs para registro: %v", err)
		return
	}
	defer f.Close()

	data, _ := json.Marshal(entry) // Conversão para Json
	f.Write(data)
	f.WriteString("\n")
}
