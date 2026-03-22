package main

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func main() {
	credenciais, err := os.ReadFile("assets/json/cred.json")
	if err != nil {
		panic("client_secret.json não encontrado!")
	}

	config, err := google.ConfigFromJSON(credenciais, youtube.YoutubeUploadScope)
	if err != nil {
		panic("Erro ao ler credenciais: " + err.Error())
	}

	client := getClient(config)

	service, err := youtube.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		panic("Erro ao criar serviço: " + err.Error())
	}

	pastaVideos := `C:\Users\ke947\Kaeve1\PomodoroYouTube\output\videos`

	fmt.Println("Iniciando uploads...")
	uploadTodosVideos(service, pastaVideos)
}
