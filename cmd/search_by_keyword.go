package main

import (
	"log"

	"google.golang.org/api/youtube/v3"
)

func searchByKeyword(query string, maxResults int64, service *youtube.Service) []string {

	// Make the API call to YouTube.
	call := service.Search.List([]string{"id", "snippet"}).
		Q(query).
		MaxResults(maxResults)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error querying Youtube: %v", err)
	}

	var videoIds []string

	for _, item := range response.Items {
		switch item.Id.Kind {
		case "youtube#video":
			videoIds = append(videoIds, item.Id.VideoId)
		}
	}

	return videoIds
}
