## Visão Geral

O pomoDojo API é uma CLI em Go que faz upload automático dos vídeos Pomodoro gerados para o YouTube. Ele interpreta o nome do arquivo para gerar títulos, descrições, timestamps e tags — e deleta os arquivos locais após o upload bem-sucedido.

## Stack

- **Go** — CLI, HTTP, fluxo OAuth2
- **YouTube Data API v3** — upload e metadados
- **Google OAuth2** — autenticação

## Como Funciona

```
output/videos/*.mp4 → upload.go
                        ├── parseNome()       # extrai metadados do nome do arquivo
                        ├── gerarTitulo()     # monta título para o YouTube
                        ├── gerarDescricao()  # monta descrição completa
                        ├── gerarTags()       # gera lista de tags
                        └── fazerUpload()     # faz upload + deleta arquivo local
```

O nome do arquivo define tudo. Um arquivo chamado `pomodoro_4h_50_10_pink.mp4` gera automaticamente título completo, descrição com explicação do noise e capítulos com timestamps.

## Instalação

```bash
git clone https://github.com/youruser/PomodojoAPI.git
cd PomodojoAPI
go mod tidy
```

Coloque suas credenciais OAuth do Google em:

```
assets/json/cred.json
```

> Obtenha as credenciais em [console.cloud.google.com](https://console.cloud.google.com) → APIs e Serviços → Credenciais → ID do cliente OAuth 2.0 (App para computador). A YouTube Data API v3 deve estar ativada.

## Uso

```bash
go run .
```

Na primeira execução, um navegador abre para autorização Google. O token é salvo localmente e reutilizado nas execuções seguintes.

## Convenção de Nome de Arquivo

```
pomodoro_{duração}h_{focus}_{break}_{noise}.mp4
```

| Segmento | Exemplo | Descrição |
|----------|---------|-----------|
| `duração` | `4h` | Duração total do vídeo em horas |
| `focus` | `50` | Duração da sessão de foco em minutos |
| `break` | `10` | Duração do intervalo em minutos |
| `noise` | `pink` | Tipo de ruído |

**Tipos de noise:** `pink` `brown` `green` `white` `grey` `blue`

**Exemplo:** `pomodoro_4h_50_10_pink.mp4` →  
Título: `4 Hours Pomodoro Timer | 50/10 | Pink Noise 🎀 | pomoDojo`

## Configuração

Atualize o caminho `pastaVideos` em `main.go` para apontar para sua pasta de vídeos:

```go
pastaVideos := filepath.Join("..", "PomodoroYouTube", "output", "videos")

------------------------------------------------------------------------------------------------------------------------

# 🎋 pomoDojo API

> Automated YouTube uploader for pomoDojo videos.

**[EN]** | [PT-BR](#pt-br)

---

## Overview

pomoDojo API is a Go CLI that automatically uploads generated Pomodoro videos to YouTube. It parses video filenames to generate titles, descriptions, timestamps, and tags — then deletes local files after a successful upload.

## Stack

- **Go** — CLI, HTTP, OAuth2 flow
- **YouTube Data API v3** — video upload and metadata
- **Google OAuth2** — authentication

## How It Works

```
output/videos/*.mp4 → upload.go
                        ├── parseNome()       # extracts metadata from filename
                        ├── gerarTitulo()     # builds YouTube title
                        ├── gerarDescricao()  # builds full description
                        ├── gerarTags()       # generates tag list
                        └── fazerUpload()     # uploads + deletes local file
```

Filename drives everything. A file named `pomodoro_4h_50_10_pink.mp4` automatically generates a complete title, description with noise explanation, and timestamped chapters.

## Setup

```bash
git clone https://github.com/youruser/PomodojoAPI.git
cd PomodojoAPI
go mod tidy
```

Place your Google OAuth credentials at:

```
assets/json/cred.json
```

> Obtain credentials at [console.cloud.google.com](https://console.cloud.google.com) → APIs & Services → Credentials → OAuth 2.0 Client ID (Desktop app). YouTube Data API v3 must be enabled.

## Usage

```bash
go run .
```

On first run, a browser window opens for Google authorization. The token is saved locally and reused on subsequent runs.

## File Naming Convention

```
pomodoro_{duration}h_{focus}_{break}_{noise}.mp4
```

| Segment | Example | Description |
|---------|---------|-------------|
| `duration` | `4h` | Total video duration in hours |
| `focus` | `50` | Focus session length in minutes |
| `break` | `10` | Break length in minutes |
| `noise` | `pink` | Noise type |

**Noise options:** `pink` `brown` `green` `white` `grey` `blue`

**Example:** `pomodoro_4h_50_10_pink.mp4` →  
Title: `4 Hours Pomodoro Timer | 50/10 | Pink Noise 🎀 | pomoDojo`

## Configuration

Update the `pastaVideos` path in `main.go` to point to your videos folder:

```go
pastaVideos := filepath.Join("..", "PomodoroYouTube", "output", "videos")
