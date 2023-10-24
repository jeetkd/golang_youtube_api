package main

import (
	"fmt"
	"net/smtp"
)

// SendMail : 업데이트된 재생목록을 메일로 알림
func (y *youtubeinfolists) SendMail() {
	var s string
	username := "" // 이메일을 보낼 구글 계정 입력
	passwd := ""   // 구글 앱 비밀번호
	auth := smtp.PlainAuth("", username, passwd, "smtp.gmail.com")
	subject := "Subject: [제목] 유튜브 재생목록 삭제 목록 알림\r\n"

	from := username         // 보내는 사람
	to := []string{username} //받는 사람
	msg := []byte(subject)   //보낼 메시지

	for _, v := range *y {
		s = fmt.Sprintf("%s, %s\n", v.id, v.title)
		msg = append(msg, []byte(s)...) // 메시지 추가
	}

	err := smtp.SendMail("smtp.gmail.com:587", auth, from, to, msg)
	if err != nil {
		WriteLogFile("youtube.log", "Failed to Send mail", err)
	}
}
