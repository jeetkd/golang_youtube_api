package main

import (
	"context"
	"golang.org/x/oauth2"
	"net/http"
	"sync"
)

func getClient(config *oauth2.Config) *http.Client {
	// Retrieve a token, if it exists, or prompts the user to authenticate.
	token := getTokenFromWeb(config)
	return config.Client(context.Background(), token)
}

// getTokenFromWeb uses the provided OAuth 2.0 Config to request a Token. (제공된 Auth 2.0 구성을 사용해서 토큰을 요청)
// It returns the retrieved Token.(유효한 토큰을 반환)
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	var code string
	var wg sync.WaitGroup
	wg.Add(1)

	//로그인 인증 페이지 반환
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	// 로그인 인증 후 code 요청 처리 핸들러 등록
	http.HandleFunc("/login/callback", func(w http.ResponseWriter, req *http.Request) {
		value := req.URL.Query()
		// code를 가져옴
		code = value.Get("code")
		wg.Done()
	})

	//코드를 받아올 서버 실행(고루틴)
	go http.ListenAndServe(":8080", nil)
	//로그인 인증
	GoogleLoginAuth(authURL)
	wg.Wait()

	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		WriteLogFile("youtube.log", "Unable to retrieve token from web: ", err)
		//log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return token
}
