package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

const TOKEN_FILE = "json/token.json"

func getClient(config *oauth2.Config) *http.Client {
	token, err := lerToken()

	if err != nil {
		token = autorizarNoNavegador(config)
		salvarToken(token)
	}

	if !token.Valid() {
		tokenSource := config.TokenSource(context.Background(), token)
		novoToken, err := tokenSource.Token()
		if err != nil {
			panic("Erro ao renovar token: " + err.Error())
		}
		salvarToken(novoToken)
		token = novoToken
	}

	return config.Client(context.Background(), token)
}

func lerToken() (*oauth2.Token, error) {
	f, err := os.Open(TOKEN_FILE)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	return token, err
}

func salvarToken(token *oauth2.Token) {
	f, err := os.Create(TOKEN_FILE)
	if err != nil {
		panic("Erro ao salvar token: " + err.Error())
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
	fmt.Println("Token salvo!")
}

func autorizarNoNavegador(config *oauth2.Config) *oauth2.Token {
	url := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Printf("Abre esse link no navegador:\n%v\n\n", url)
	fmt.Print("Cola o código aqui: ")

	var code string
	fmt.Scan(&code)

	token, err := config.Exchange(context.TODO(), code)
	if err != nil {
		panic("Erro ao obter token: " + err.Error())
	}

	return token
}
