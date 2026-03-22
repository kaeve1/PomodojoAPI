package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"google.golang.org/api/youtube/v3"
)

// ─── ESTRUTURA COM INFO DO VÍDEO ──────────────────────────
type VideoInfo struct {
	Nome          string
	WorkMinutes   int
	BreakMinutes  int
	Ciclos        int
	TotalMinutes  int
	TotalDuration string
	NoiseType     string
	NoiseEmoji    string
}

// ─── PARSE DO NOME DO ARQUIVO ─────────────────────────────
// ex: "pomodoro_4h_50_10_pink" ou "pomodoro_4h_50_10"
func parseNome(nome string) VideoInfo {
	parts := strings.Split(nome, "_")

	info := VideoInfo{Nome: nome}

	// extrai work e break
	for i, p := range parts {
		if i > 0 && isNumber(p) {
			if info.WorkMinutes == 0 {
				info.WorkMinutes, _ = strconv.Atoi(p)
			} else if info.BreakMinutes == 0 {
				info.BreakMinutes, _ = strconv.Atoi(p)
				break
			}
		}
	}

	// extrai ciclos da parte "Xh" ou "Xc"
	for _, p := range parts {
		if strings.HasSuffix(p, "h") {
			horas, err := strconv.Atoi(strings.TrimSuffix(p, "h"))
			if err == nil && info.WorkMinutes > 0 && info.BreakMinutes > 0 {
				cicloMinutos := info.WorkMinutes + info.BreakMinutes
				info.Ciclos = (horas * 60) / cicloMinutos
				info.TotalMinutes = horas * 60
			}
		}
	}

	// fallback se não achou ciclos
	if info.Ciclos == 0 {
		info.Ciclos = 4
	}
	if info.TotalMinutes == 0 {
		info.TotalMinutes = (info.WorkMinutes + info.BreakMinutes) * info.Ciclos
	}

	// formata duração total
	horas := info.TotalMinutes / 60
	minutos := info.TotalMinutes % 60
	if minutos == 0 {
		info.TotalDuration = fmt.Sprintf("%d Hour", horas)
		if horas > 1 {
			info.TotalDuration += "s"
		}
	} else {
		info.TotalDuration = fmt.Sprintf("%dh%02d", horas, minutos)
	}

	// extrai tipo de noise
	info.NoiseType = "Pink"
	info.NoiseEmoji = "🎀"
	for _, p := range parts {
		switch strings.ToLower(p) {
		case "pink":
			info.NoiseType = "Pink"
			info.NoiseEmoji = "🎀"
		case "brown":
			info.NoiseType = "Brown"
			info.NoiseEmoji = "🤎"
		case "green":
			info.NoiseType = "Green"
			info.NoiseEmoji = "🌿"
		case "white":
			info.NoiseType = "White"
			info.NoiseEmoji = "🤍"
		case "grey", "gray":
			info.NoiseType = "Grey"
			info.NoiseEmoji = "🩶"
		case "blue":
			info.NoiseType = "Blue"
			info.NoiseEmoji = "💙"
		}
	}

	return info
}

func isNumber(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// ─── GERA TIMESTAMPS ──────────────────────────────────────
func gerarTimestamps(info VideoInfo) string {
	var sb strings.Builder
	cursor := 0

	for c := 1; c <= info.Ciclos; c++ {
		h := cursor / 60
		m := cursor % 60
		sb.WriteString(fmt.Sprintf("%02d:%02d:00 ⏰ FOCUS %d\n", h, m, c))
		cursor += info.WorkMinutes

		if c < info.Ciclos {
			h = cursor / 60
			m = cursor % 60
			sb.WriteString(fmt.Sprintf("%02d:%02d:00 ☕️ BREAK %d\n", h, m, c))
			cursor += info.BreakMinutes
		}
	}

	return strings.TrimSpace(sb.String())
}

// ─── GERA TÍTULO ──────────────────────────────────────────
func gerarTitulo(info VideoInfo) string {
	return fmt.Sprintf(
		"%s Pomodoro Timer | %d/%d | %s Noise %s | pomoDojo",
		info.TotalDuration,
		info.WorkMinutes,
		info.BreakMinutes,
		info.NoiseType,
		info.NoiseEmoji,
	)
}

// ─── DESCRIÇÃO DO NOISE ───────────────────────────────────
func descricaoNoise(info VideoInfo) string {
	switch strings.ToLower(info.NoiseType) {
	case "pink":
		return `🎀 About Pink Noise
Pink noise sits between white and brown noise in frequency.
It's gentle, consistent, and scientifically shown to improve focus and reduce distractions — making it ideal for study sessions, deep work, and ADHD-friendly environments.`
	case "brown":
		return `🤎 About Brown Noise
Brown noise has a deeper, lower frequency than pink or white noise.
It resembles the sound of a strong wind or distant thunder — powerful, grounding, and incredibly effective for blocking distractions and entering a state of deep focus.
Particularly popular among people with ADHD.`
	case "green":
		return `🌿 About Green Noise
Green noise sits in the middle of the sound spectrum, resembling natural environments like forests, rivers, and rainfall.
It's soothing, organic, and helps create a calm atmosphere that supports sustained concentration without mental fatigue.`
	case "white":
		return `🤍 About White Noise
White noise contains all frequencies at equal intensity, creating a consistent sound that masks background noise effectively.
It helps maintain focus by creating a neutral sonic environment — ideal for open spaces, offices, or noisy study areas.`
	case "grey":
		return `🩶 About Grey Noise
Grey noise is spectrally flat to human perception — meaning every frequency sounds equally loud to your ears.
It's one of the most balanced and comfortable noises to listen to for extended periods, making it perfect for long study sessions.`
	case "blue":
		return `💙 About Blue Noise
Blue noise emphasizes higher frequencies, creating a bright, airy sound.
It's energizing and stimulating — ideal for tasks that require alertness and active thinking rather than deep relaxation.`
	default:
		return `🎧 About the Noise
This ambient noise is carefully designed to help you maintain focus and reduce distractions during your study sessions.`
	}
}

// ─── GERA DESCRIÇÃO ─────O NOME DO vIDEO DEVE SER ASSIM:
// "Formato: pomodoro_{DURACAO}h_{FOCUS}_{BREAK}_{NOISE} | Ex: pomodoro_4h_50_10_pink"
func gerarDescricao(info VideoInfo) string {
	timestamps := gerarTimestamps(info)

	return fmt.Sprintf(`⏰ Pomodoro Timer — %d/%d | %s
No music. No distractions. Just focus.
Pomodojo is built to help you train focus and consistency through structured work sessions.
This video follows the Pomodoro Technique: work in focused intervals, then take short breaks to recover and maintain high performance over time.
🔔 A bell will ring at the end of each session.
━━━━━━━━━━━━━━━━━━
📌 Session Info
- Work: %d minutes
- Break: %d minutes
- Cycles: %d
- Total Duration: %s
━━━━━━━━━━━━━━━━━━
📌 How to Use
1. Define one clear task before starting
2. Focus completely during work sessions
3. Step away during breaks
4. Repeat the cycle and build momentum
━━━━━━━━━━━━━━━━━━
%s
━━━━━━━━━━━━━━━━━━
💡 Benefits
- Improves concentration
- Reduces fatigue
- Builds discipline
- Helps prevent procrastination
━━━━━━━━━━━━━━━━━━
🧠 Ideal For
- Studying
- Coding
- Reading
- Writing
- Deep work
- ADHD-friendly focus sessions
━━━━━━━━━━━━━━━━━━
🚀 About the Channel
Pomodojo is focused on one thing: helping you get into deep work.
No music. No distractions. Just time, structure, and execution.
━━━━━━━━━━━━━━━━━━
⏱ Timestamps
%s
━━━━━━━━━━━━━━━━━━
#pomodoro #studytimer #focus #deepwork #productivity #discipline #timemanagement #study #adhd #pomodojo`,
		info.WorkMinutes, info.BreakMinutes, info.TotalDuration,
		info.WorkMinutes, info.BreakMinutes, info.Ciclos, info.TotalDuration,
		descricaoNoise(info),
		timestamps,
	)
}

// ─── GERA TAGS ────────────────────────────────────────────
func gerarTags(info VideoInfo) []string {
	tags := []string{
		"pomodoro", "study timer", "focus", "pomodojo",
		"deep work", "study with me", "productivity",
		"discipline", "time management", "adhd focus",
		fmt.Sprintf("%d/%d pomodoro", info.WorkMinutes, info.BreakMinutes),
		fmt.Sprintf("%s noise", strings.ToLower(info.NoiseType)),
		fmt.Sprintf("%s pomodoro timer", info.TotalDuration),
	}
	return tags
}

// ─── UPLOAD DE TODOS OS VÍDEOS ────────────────────────────
func uploadTodosVideos(service *youtube.Service, pastaVideos string) {

	entries, err := os.ReadDir(pastaVideos)
	if err != nil {
		panic("Pasta não encontrada: " + err.Error())
	}

	total := 0
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".mp4") {
			total++
		}
	}

	if total == 0 {
		fmt.Println(" Nenhum vídeo .mp4 encontrado!")
		return
	}

	fmt.Printf("📁 %d vídeo(s) encontrado(s)\n\n", total)

	count := 0
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".mp4") {
			continue
		}

		count++
		caminho := filepath.Join(pastaVideos, entry.Name())
		nome := strings.TrimSuffix(entry.Name(), ".mp4")
		info := parseNome(nome)

		fmt.Printf("📤 [%d/%d] Uploading: %s\n", count, total, entry.Name())
		fmt.Printf("  Título: %s\n", gerarTitulo(info))

		videoId, err := fazerUpload(service, caminho, info)
		if err != nil {
			fmt.Printf("Erro no upload de %s: %v\n", entry.Name(), err)
			continue
		}

		fmt.Printf("✅ Publicado: https://youtube.com/watch?v=%s\n", videoId)

		//deleta o arquivo após upload bem sucedido
		err = os.Remove(caminho)
		if err != nil {
			fmt.Printf("Não foi possível deletar %s: %v\n", entry.Name(), err)
		} else {
			fmt.Printf("%s deletado!\n\n", entry.Name())
		}
	}

	fmt.Println("🏁 Todos os uploads concluídos!")
}

// ─── UPLOAD DE UM VÍDEO ───────────────────────────────────
func fazerUpload(service *youtube.Service, caminho string, info VideoInfo) (string, error) {

	file, err := os.Open(caminho)
	if err != nil {
		return "", fmt.Errorf("erro ao abrir vídeo: %w", err)
	}
	defer file.Close()

	video := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       gerarTitulo(info),
			Description: gerarDescricao(info),
			Tags:        gerarTags(info),
			CategoryId:  "27",
		},
		Status: &youtube.VideoStatus{
			PrivacyStatus: "public",
		},
	}

	call := service.Videos.Insert([]string{"snippet", "status"}, video)
	resp, err := call.Media(file).Do()
	if err != nil {
		return "", fmt.Errorf("erro no upload: %w", err)
	}

	return resp.Id, nil
}
