# Web Proxy com Controle de Conteúdo

Este projeto é uma implementação de um servidor proxy HTTP didático desenvolvido para a disciplina de Sistemas para Internet 2. O proxy permite o repasse de requisições, bloqueio de domínios específicos e filtragem de conteúdo HTML (substituição de palavras).

## Estrutura do Projeto

O projeto está organizado nas seguintes pastas:
- **`scripts/`**: Contém todo o código-fonte em Go (`main.go`, `config.go`, `handler.go`, `logger.go`).
- **`json/`**: Contém os arquivos de configuração e logs (`blocked.json`, `words.json`, `log.json`).
- **`html/`**: Contém os arquivos de html (`blocked.html`).

## Como Instalar e Executar

### Pré-requisitos
- [Go](https://golang.org/dl/) instalado (versão 1.16 ou superior recomendada).

### Execução
1. Clone o repositório ou baixe os arquivos.
2. No terminal, na raiz do projeto, execute:
   ```bash
   go run ./scripts
   ```
3. O servidor iniciará na porta `5000`.

### Uso
O proxy aceita requisições no formato:
`http://localhost:5000/http://www.exemplo.com`

Exemplos:
- **Acesso Transparente**: `http://localhost:5000/http://www.google.com`
- **Acesso Bloqueado**: `http://localhost:5000/http://www.facebook.com` (ou qualquer domínio em `json/blocked.json`)
- **Filtro de Conteúdo**: Acesse qualquer página HTTP que contenha palavras cadastradas em `json/words.json`.

