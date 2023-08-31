package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"io/ioutil"
	"net/smtp"
	"os"
	"strings"
)

const (
	maxResult = 50
)

type YoutubeInfo struct {
	title string
	id    string
}

type UpdateInfo struct {
	addlist    youtubeinfolists
	deletelist youtubeinfolists
}
type youtubeinfolists []YoutubeInfo

func main() {
	// 현재 경로 읽어옴
	currentDir, err := os.Getwd()
	if err != nil {
		WriteLogFile("youtube.log", "Failed to read path: ", err)
	}

	//Read json file (json 파일 읽어옴)
	jsonKey, err := ioutil.ReadFile(currentDir + "/client.json")
	if err != nil {
		WriteLogFile("youtube.log", "Failed to read JSON file: ", err)
	}

	// Set up the OAuth 2.0 configuration (OAuth 2.0 구성 환경 설정)
	config, err := google.ConfigFromJSON(jsonKey, youtube.YoutubeReadonlyScope)
	if err != nil {
		WriteLogFile("youtube.log", "Unable to parse client secret file: ", err)
	}

	// Create a new OAuth 2.0 client (새로운 OAuth 2.0 클라이언트를 생성 합니다.)
	client := getClient(config)

	// Set up the YouTube API client using the authenticated client (인증된 클라이언트를 사용하여 유튜브 API 클라이언트를 설정)
	service, err := youtube.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		WriteLogFile("youtube.log", "Unable to create YouTube service: ", err)
	}

	var playlistsOrigin youtubeinfolists //내 유튜브 계정에서 가져온 리스트
	var playlistsLocal youtubeinfolists  //내부 플레이 리스트.txt

	playlistID := "" //재생목록 ID
	nextPageToken := ""

	for {
		// Retrieve the playlist items from the private playlist (비공개 재생목록으로부터 재생목록 아이템들을 불러옵니다.)
		playlistItemsResponse := PlayListRes(service, "snippet", maxResult, nextPageToken, playlistID)

		playlistsOrigin = append(playlistsOrigin, BringPlaylists(playlistItemsResponse)...)
		nextPageToken = playlistItemsResponse.NextPageToken //최대 가져올수 있는 아이템들이 50이므로 다음 토큰으로 넘어와서 가져와야함.
		if nextPageToken == "" {
			break
		}
	}

	//로컬에 있는 재생목록들을 저장한 파일들을 읽어옴
	playlistsLocal.ReadFile("playlistsLocal")

	//로컬에 있는 재생목록들을 기록한 파일과 유튜브에서 가져온 재생목록 파일을 비교 후 업데이트할 리스트 반환
	updatelist := playlistsLocal.CheckPlaylists(playlistsOrigin)
	//삭제 또는 추가할 리스트들을 로컬 파일에 업데이트
	playlistsOrigin.UpdatePlaylists(updatelist)
	WriteLogFile("youtube.log", "Success to execute file", errors.New("파일이 성공적으로 실행 후 완료 되었습니다"))
}

// WriteFile : id와 제목을 playlistsLocal.txt 파일로 만듬
func (y *youtubeinfolists) WriteFile(name string) int {
	var num int //작성한 라인 수 반환(song numbers)

	currentDir, err := os.Getwd()
	if err != nil {
		WriteLogFile("youtube.log", "Failed to read path: ", err)
	}

	file, err := os.Create(currentDir + "/" + name + ".txt")
	if err != nil {
		WriteLogFile("youtube.log", "Failed to Create file: ", err)
	}
	defer file.Close()

	playLists := *y // id와 title 복사
	for _, list := range playLists {
		s := fmt.Sprintf("%s, %s\n", list.id, list.title)
		_, err := file.Write([]byte(s))
		if err != nil {
			WriteLogFile("youtube.log", "Failed to write file", err)
		}
		num = num + 1
	}
	return num
}

// ReadFile : playlistsLocal.txt 파일을 읽어와서 *y에 넣어주고 읽어온 라인 수 반환
func (y *youtubeinfolists) ReadFile(name string) int {
	var s []string
	var num int
	tempLists := YoutubeInfo{}

	currentDir, err := os.Getwd()
	if err != nil {
		WriteLogFile("youtube.log", "Failed to read path: ", err)
	}

	file, err := os.Open(currentDir + "/" + name + ".txt")
	if err != nil {
		WriteLogFile("youtube.log", "Failed to Open file", err)
	}
	defer file.Close()

	//Get title, id
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break //완료 시 탈출
		}

		s = strings.Split(line, ", ")
		tempLists.id = s[0]
		tempLists.title = strings.TrimRight(s[1], "\n")
		*y = append(*y, YoutubeInfo{id: tempLists.id, title: tempLists.title}) // 읽어온 데이터 추가
		num = num + 1
	}
	return num
}

// CheckPlaylists : 저장된 원본 리스트(y)와 가져온 리스트(playlists)를 비교 후 삭제와 추가할 재생목록들을 반환
func (y youtubeinfolists) CheckPlaylists(playlists youtubeinfolists) UpdateInfo {
	var updateListZ UpdateInfo //zero(0)을 가진 채널 목록들
	var updateList UpdateInfo  // 삭제 또는 추가할 채널 목록들

	//변경되지 않은 목록들의 id를 "0"으로 변경 후 UpdateInfo 구조체 반환
	updateListZ = ComparePlaylistsZ(y, playlists)

	// 유튜브 재생목록에서 추가된 재생목록 추가
	for _, v := range updateListZ.addlist {
		if v.id != "0" {
			updateList.addlist = append(updateList.addlist, v)
		}
	}
	// 유튜브 재생목록에서 삭제된 재생목록 추가
	for _, v := range updateListZ.deletelist {
		if v.id != "0" {
			updateList.deletelist = append(updateList.deletelist, v)
		}
	}

	return updateList
}

// UpdatePlaylists : 재생목록 업데이트 (y = 가져온 재생목록)
func (y youtubeinfolists) UpdatePlaylists(info UpdateInfo) {
	//삭제된 재생목록 존재시 실행
	if info.deletelist != nil {
		y.WriteFile("playlistsLocal")
		fmt.Println("Delete playlists : ", info.deletelist)
		info.deletelist.SendMail() // 삭제된 재생목록을 메일로 알림
	}

	//추가된 재생목록 존재시 실행
	if info.addlist != nil {
		y.WriteFile("playlistsLocal")
		fmt.Println("Add playlists : ", info.addlist)
	}
}

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

// BringPlaylists : 재생목록을 가져옵니다.
func BringPlaylists(res *youtube.PlaylistItemListResponse) []YoutubeInfo {
	playLists := make([]YoutubeInfo, 0)
	//numlists := 0
	for _, item := range res.Items {
		songTitle := item.Snippet.Title
		videoId := item.Snippet.ResourceId.VideoId
		info := YoutubeInfo{title: songTitle, id: videoId}
		playLists = append(playLists, info)
		//numlists = numlists + 1
	}
	return playLists
}

// PlayListRes : 재생목록을 가지고 있는 응답 리스트를 반환
func PlayListRes(service *youtube.Service, part string, maxResults int64, pageToken string, playlistId string) *youtube.PlaylistItemListResponse {
	// Retrieve the playlist items from the private playlist (비공개 재생목록으로부터 재생목록 아이템들을 불러옵니다.)
	call := service.PlaylistItems.List([]string{part}).
		PlaylistId(playlistId). // 재생목록 ID 설정
		MaxResults(maxResult). // 가져올 재생목록 item 최대값 설정
		PageToken(pageToken) // 다음 재생목록에 대한 토큰

	// "youtube.playlistItems.list" 호출 실행.
	response, err := call.Do()
	if err != nil {
		WriteLogFile("youtube.log", "Unable to retrieve playlist items: ", err)
	}
	return response
}

// ComparePlaylistsZ : 두 재생 목록 비교 후 삭제, 추가할지 재생목록 구분 후 UpdateInfo구조체 반환
func ComparePlaylistsZ(playlists1, playlists2 youtubeinfolists) UpdateInfo {
	var updateList UpdateInfo
	copyLists1 := make([]YoutubeInfo, len(playlists1))
	copyLists2 := make([]YoutubeInfo, len(playlists2))

	copy(copyLists1, playlists1)
	copy(copyLists2, playlists2)

	//삭제 또는 추가할 재생목록들을 구분(변경되지 않은 목록들의 id를 "0"으로 변경)
	for i1, v1 := range playlists1 {
		for i2, v2 := range playlists2 {
			if v1.id == v2.id {
				copyLists1[i1].id = "0"
				copyLists2[i2].id = "0"
				continue
			}
		}
	}

	updateList.deletelist = copyLists1
	updateList.addlist = copyLists2
	return updateList
}
