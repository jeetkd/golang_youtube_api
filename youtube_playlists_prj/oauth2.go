package main

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"log"
	"net/http"
)

// getClient retrieves a valid OAuth 2.0 client.(유효한 OAuth 2.0 클라이언트를 불러옵니다.)
// 토큰을 *http.Client에 넣어주고 반환
func getClient(config *oauth2.Config) *http.Client {
	// Retrieve a token, if it exists, or prompts the user to authenticate.
	token := getTokenFromWeb(config)
	return config.Client(context.Background(), token)
}

// getTokenFromWeb uses the provided OAuth 2.0 Config to request a Token. (제공된 Auth 2.0 구성을 사용해서 토큰을 요청)
// It returns the retrieved Token.(유효한 토큰을 반환)
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser, then type the "+
		"authorization code: \n%v\n", authURL)

	//토큰 입력
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return token
}
