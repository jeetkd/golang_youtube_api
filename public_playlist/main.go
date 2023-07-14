package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"log"
)

func main() {
	//Set up API Key and playlistID(API 키와 playlistID 설정)
	apiKey := flag.String("api-key", "AIzaSyDwG8f1UeezbbYg0_SQ9KKeQto3iNX2qMk", "YouTube API 키")
	playlistID := flag.String("playlist-id", "PLp-ofiyo_L4kUdOxLd0Jn1-G9neeA4ZLX", "재생목록 ID")
	flag.Parse()

	//Create Service(서비스 생성)
	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(*apiKey))
	reqPlaylist := service.PlaylistItems.List([]string{"snippet"}). // Set up snippet option
									PlaylistId(*playlistID). //Set up playlistID
									MaxResults(50)           //Set up Max Result

	// Do Excute "youtube.playlistItems.list"(PlaylistItemListResponse를 응답으로 받음)
	resPlaylist, err := reqPlaylist.Do()
	if err != nil {
		log.Fatalf("Error fetching playlist items: %v", err.Error())
	}

	// Get playlist Items(재생 목록과 비디오 ID를 가져옴)
	for _, playlistItem := range resPlaylist.Items {
		title := playlistItem.Snippet.Title
		videoId := playlistItem.Snippet.ResourceId.VideoId
		fmt.Printf("%v, (%v)\r\n", title, videoId)
	}

}
