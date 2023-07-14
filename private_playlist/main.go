package main

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	//Read json file (json 파일 읽어옴)
	jsonKey, err := ioutil.ReadFile("client.json")
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

	// Retrieve the playlist ID for the private playlist you want to access (접근을 원하는 비공개 재생목록 ID)
	playlistID := "자신의 계정에 비공개 재생목록의 ID를 넣어주세요"

	// Retrieve the playlist items from the private playlist (비공개 재생목록으로부터 재생목록 아이템들을 불러옵니다.)
	playlistItemsCall := service.PlaylistItems.List([]string{"snippet"}).
		PlaylistId(playlistID). // 재생목록 ID 설정
		MaxResults(50) // Adjust the maximum number of results as per your requirements(가져올 재생목록 item 최대값 설정)

	playlistItemsResponse, err := playlistItemsCall.Do() // "youtube.playlistItems.list" 호출 실행.
	if err != nil {
		log.Fatalf("Unable to retrieve playlist items: %v", err)
	}

	// Process and display the playlist items //가져온 재생목록 아이템들을 제목과 ID 출력
	for _, item := range playlistItemsResponse.Items {
		title := item.Snippet.Title
		videoID := item.Snippet.ResourceId.VideoId
		fmt.Printf("Title: %s, Video ID: %s\n", title, videoID)
	}
}

// getClient retrieves a valid OAuth 2.0 client.(유효한 OAuth 2.0 클라이언트를 불러옵니다.)
// 토큰트로 *http.Client를 생성하고 반환
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
