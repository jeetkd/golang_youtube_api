package main

import (
	"log"
	"os"
)

func OpenFile(name string) *os.File {
	//O_CREATE : 파일이 존재 하지 않으면 생성, O_APPEND : 이어서 작성, O_WRONLY : 쓰기전용
	f, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func WriteLogFile(name, message string, err error) {
	//파일 열기
	f := OpenFile(name)
	// 출력을 f로 설정
	log.SetOutput(f)
	// 로그 기록
	log.Println(message, err)
	//파일 닫기
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
