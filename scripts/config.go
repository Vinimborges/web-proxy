package main

import (
	"encoding/json"
	"os"
	"regexp"
	"sync"
)

// Struct que armazena os sites bloqueados
type BlockedConfig struct {
	Bloqueados []string `json:"bloqueados"`
}

var (
	blockedDomains map[string]bool   // Armazena os sites bloqueados
	wordFilters    map[string]string // Armazena os filtros de palavras
	regexFilters   []struct {
		regex       *regexp.Regexp // Armazena os padrões de regex
		replacement string         // Armazena o texto de substituição
	}
	configMutex sync.RWMutex // Mutex para proteger o acesso às variáveis compartilhadas
)

func loadConfig() {
	configMutex.Lock()
	defer configMutex.Unlock()

	// Carrega sites bloqueados
	blockedDomains = make(map[string]bool)
	blockedData, err := os.ReadFile("json/blocked.json")
	if err == nil {
		var bc BlockedConfig
		if err := json.Unmarshal(blockedData, &bc); err == nil {
			for _, domain := range bc.Bloqueados {
				blockedDomains[domain] = true
			}
		}
	}

	// Carrega filtro de palavras
	wordFilters = make(map[string]string)
	regexFilters = nil
	wordsData, err := os.ReadFile("json/words.json")
	if err == nil {
		if err := json.Unmarshal(wordsData, &wordFilters); err == nil {
			for word, replacement := range wordFilters {
				re, err := regexp.Compile("(?i)\\b" + regexp.QuoteMeta(word) + "\\b") // Compila regex case-insensitive e com limites de palavra
				if err == nil {
					regexFilters = append(regexFilters, struct {
						regex       *regexp.Regexp
						replacement string
					}{re, replacement})
				}
			}
		}
	}
}
