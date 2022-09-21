package main

import (
	"log"
	"os"
	"sync"

	"github.com/tjarratt/babble"
	"google.golang.org/api/youtube/v3"
)

const (
	pathBase     = "videos"
	filename     = "life_remote_control.mp4"
	title        = "Life Remote Control"
	description  = "A series of miscellaneous film footage are put together, to represent life as if watching it in front of a television screen with a remote control."
	privacy      = "unlisted"
	category     = "1"
	maxResults   = 50
	maxDownloads = 30
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

	babbler := babble.NewBabbler()
	babbler.Count = 10
	babbler.Separator = ","
	keywords := babbler.Babble()
	log.Printf(keywords)

	client := getClient(youtube.YoutubeScope)

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	videoIds := searchByKeyword(keywords, maxResults, service)

	var wg sync.WaitGroup
	count := 0
	for _, id := range videoIds {
		if count > maxDownloads {
			break
		}
		wg.Add(1)
		go downloadVideo(id, &err, &wg)
		if err != nil {
			log.Printf("Error: skipping video download: %v", err)
		}
		count++
	}

	wg.Wait()

	e := joinVideos()
	if e != nil {
		log.Fatalf("Error joining videos")
	}

	e = os.RemoveAll(pathBase + "/downloads")
	if e != nil {
		log.Printf("Error deleting files in downloads directory")
	}

	uploadVideo(title, description, category, privacy, keywords, pathBase+"/output/"+filename, *service)

}
