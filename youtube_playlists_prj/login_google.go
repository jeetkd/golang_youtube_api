package main

import (
	cu "github.com/Davincible/chromedp-undetected"
	"github.com/chromedp/chromedp"
	"time"
)

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
		WriteLogFile("youtube.log", "cu.New() = Fail to create undetected Chrome executor : ", err)
		//log.Fatalf("cu.New() = Fail to create undetected Chrome executor : %v", err)
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
		WriteLogFile("youtube.log", "chromedp.Run() = Fail to login in chromedp.Run() : ", err)
		//log.Fatalf("chromedp.Run() = Fail to login in chromedp.Run() : %v", err)
		//이메일로 알림!
	}
}
