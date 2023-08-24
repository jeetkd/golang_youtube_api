package main

import (
	"context"
	cu "github.com/Davincible/chromedp-undetected"
	"github.com/chromedp/chromedp"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"sync"
	"time"
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
	//fmt.Println("code :", code)
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return token
}

// GoogleLoginAuth : 구글 로그인 인증 자동화
func GoogleLoginAuth(url string) {
	var email = ""    // 로그인 이메일
	var password = "" // 이메일 비밀번호

	ctx, cancel, err := cu.New(cu.NewConfig(
		// Remove this if you want to see a browser window.
		//cu.WithHeadless(),
		// If the webelement is not found within 300 seconds, timeout.
		cu.WithTimeout(300 * time.Second),
	))
	defer cancel()

	if err != nil {
		log.Fatalf("cu.New() = Fail to create undetected Chrome executor : %v", err)
	}

	if err := chromedp.Run(ctx,
		// url 열기
		chromedp.Navigate(url),
		chromedp.WaitVisible(`#identifierId`),
		// ID 입력, 다음 클릭
		chromedp.SetAttributeValue("#identifierId", "value", email, chromedp.NodeVisible), //id
		chromedp.Sleep(5*time.Second),
		chromedp.Click(`#identifierNext > div > button > div.VfPpkd-RLmnJb`, chromedp.NodeVisible), //다음

		// password 입력, 다음 클릭
		chromedp.WaitVisible(`#password > div.aCsJod.oJeWuf > div > div.Xb9hP > input`),
		chromedp.SetAttributeValue("#password > div.aCsJod.oJeWuf > div > div.Xb9hP > input", "value", password, chromedp.NodeVisible), // 비밀번호
		chromedp.Sleep(5*time.Second),
		chromedp.Click(`#passwordNext > div > button > div.VfPpkd-RLmnJb`, chromedp.NodeVisible), // 다음
		chromedp.Sleep(10*time.Second),
	); err != nil {
		log.Fatalf("chromedp.Run() = Fail to login in chromedp.Run() : %v", err)
		//이메일로 알림!
	}
}
