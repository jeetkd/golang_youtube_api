package main

import (
	"bufio"
	"context"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

const (
	maxResult = 50
)

type YoutubeInfo struct {
	title string
	id    string
}

func main() {
	// 현재 경로 읽어옴
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("오류:", err)
		return
	}

	//Read json file (json 파일 읽어옴)
	jsonKey, err := ioutil.ReadFile(currentDir + "/youtube_playlists_prj/client.json")
	if err != nil {
		log.Fatalf("Failed to read JSON file: %v", err)
	}

	// Set up the OAuth 2.0 configuration (OAuth 2.0 구성 환경 설정)
	config, err := google.ConfigFromJSON(jsonKey, youtube.YoutubeReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file: %v", err)
	}

	// Create a new OAuth 2.0 client (새로운 OAuth 2.0 클라이언트를 생성 합니다.)
	client := getClient(config)

	// Set up the YouTube API client using the authenticated client (인증된 클라이언트를 사용하여 유튜브 API 클라이언트를 설정)
	service, err := youtube.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to create YouTube service: %v", err)
	}

	for {
		playlists := make([]YoutubeInfo, 0)
		playlistID := "PLRHMBOXyVHiYNfq95tKlEPqxpAKw9lMb2"
		nextPageToken := ""
		for {
			// Retrieve the playlist items from the private playlist (비공개 재생목록으로부터 재생목록 아이템들을 불러옵니다.)
			playlistItemsResponse := playlistsList(service, "snippet", maxResult, nextPageToken, playlistID)

			playlists = append(playlists, bringPlaylists(playlistItemsResponse)...)
			nextPageToken = playlistItemsResponse.NextPageToken //최대 가져올수 있는 아이템들이 50이므로 다음 토큰으로 넘어와서 가져와야함.
			if nextPageToken == "" {
				break
			}
		}
		fmt.Println(len(playlists), cap(playlists))
		//fmt.Println(readFile("playlists.txt"))
		writeFile(playlists)
		time.Sleep(30 * time.Minute)
	}
}

// id와 제목을 playlists.txt 파일로 만듬
func writeFile(playlists []YoutubeInfo) {
	file, err := os.Create("playlists.txt")
	if err != nil {
		fmt.Println("오류:", err)
		return
	}
	defer file.Close()

	for _, list := range playlists {
		s := fmt.Sprintf("%s, %s\n", list.id, list.title)
		_, err := file.Write([]byte(s))
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}
}

//playlists.txt 파일을 읽어와서 []YoutubeInfo에 넣어주고 반환
func readFile(name string) []YoutubeInfo {
	playlists := make([]YoutubeInfo, 0)
	templists := YoutubeInfo{}
	var s []string

	file, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	//Get title, id
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		s = strings.Split(line, ", ")
		templists.id = s[0]
		templists.title = strings.TrimRight(s[1], "\n")
		playlists = append(playlists, YoutubeInfo{id: templists.id, title: templists.title})
	}
	return playlists
}

//재생목록을 가져옵니다.
func bringPlaylists(res *youtube.PlaylistItemListResponse) []YoutubeInfo {
	playlists := make([]YoutubeInfo, 0)
	//numPlaylist := 0
	for _, item := range res.Items {
		songTitle := item.Snippet.Title
		videoId := item.Snippet.ResourceId.VideoId
		info := YoutubeInfo{title: songTitle, id: videoId}
		playlists = append(playlists, info)
		//numPlaylist = numPlaylist + 1
	}
	return playlists
}

func playlistsList(service *youtube.Service, part string, maxResults int64, pageToken string, playlistId string) *youtube.PlaylistItemListResponse {
	// Retrieve the playlist items from the private playlist (비공개 재생목록으로부터 재생목록 아이템들을 불러옵니다.)
	call := service.PlaylistItems.List([]string{part}).
		PlaylistId(playlistId). // 재생목록 ID 설정
		MaxResults(maxResult).  // 가져올 재생목록 item 최대값 설정
		PageToken(pageToken)

	// "youtube.playlistItems.list" 호출 실행.
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Unable to retrieve playlist items: %v", err)
	}
	return response
}

// 재생목록 업데이트
func updateFile() {

}

// 재생목록 삭제
func deleteList() {

}

// 메일 보내기
func sendMail() {

}
