package main

import (
	"log"
	"os"
	"sync"

	"google.golang.org/api/youtube/v3"
)

const (
	pathBase    = "videos"
	filename    = "life_remote_control.mp4"
	title       = "Life Remote Control"
	description = "A series of miscellaneous film footage are put together, to represent life as if watching it in front of a television screen with a remote control."
	privacy     = "public"
	category    = "1"
)

func main() {
	err := os.MkdirAll(pathBase+"/downloads", os.ModePerm)
	if err != nil {
		log.Printf("Error creating downloads path: %v", err)
	}
	err = os.MkdirAll(pathBase+"/output", os.ModePerm)
	if err != nil {
		log.Printf("Error creating upload path: %v", err)

	}

	client := getClient(youtube.YoutubeScope)

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	videoIds, keywords := getVideoIdsAndKeywords(service)

	var wg sync.WaitGroup
	for _, id := range videoIds {
		wg.Add(1)
		go downloadVideo(id, &wg)
	}

	wg.Wait()

	err = joinVideos()
	if err != nil {
		log.Fatalf("Error joining videos: %v", err)
	}

	err = os.RemoveAll(pathBase + "/downloads")
	if err != nil {
		log.Printf("Error deleting files in downloads directory: %v", err)
	}

	log.Printf("Title: %v", title)
	log.Printf("Description: %v", description)
	log.Printf("Category: %v", category)
	log.Printf("Privacy: %v", privacy)
	log.Printf("Keywords: %v", keywords)
	log.Printf("Filename: %v", filename)

	uploadVideo(title, description, category, privacy, keywords, pathBase+"/output/"+filename, *service)

	err = os.RemoveAll(pathBase + "/output")
	if err != nil {
		log.Printf("Error deleting files in output directory: %v", err)
	}

}
