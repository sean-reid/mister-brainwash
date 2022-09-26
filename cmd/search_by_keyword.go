package main

import (
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"strings"
	"time"

	"google.golang.org/api/youtube/v3"
)

const (
	keywordsPerQuery  = 5
	maxQueries        = 10
	maxVideosPerQuery = 3
	maxResults        = 50
	wordsList         = "common-words-clean.txt"
)

func getRandomWords() string {
	command := fmt.Sprintf("shuf %v | head -n%v | paste -sd ',' | head -c -1", wordsList, keywordsPerQuery)
	out, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		log.Fatalf("error generating random words: %v", err)
	}
	words := string(out)
	return words
}

func searchByKeyword(videoIds chan string, keywords chan string, service *youtube.Service) {

	query := getRandomWords()
	log.Printf("Query: %v", query)

	keywords <- query

	// Make the API call to YouTube.
	call := service.Search.List([]string{"id", "snippet"}).
		Q(query).
		MaxResults(maxResults)
	response, err := call.Do()
	if err != nil {
		log.Printf("Error querying Youtube: %v", err)
		return
	}

	var videoIdsQueried []string

	for _, item := range response.Items {
		switch item.Id.Kind {
		case "youtube#video":
			videoIdsQueried = append(videoIdsQueried, item.Id.VideoId)
		}
	}
	if len(videoIdsQueried) <= 0 {
		log.Printf("no video IDs found for this query: %v", query)
		return
	}

	rand.Shuffle(len(videoIdsQueried), func(i, j int) {
		videoIdsQueried[i], videoIdsQueried[j] = videoIdsQueried[j], videoIdsQueried[i]
	})

	for ii := 0; ii < maxVideosPerQuery; ii++ {
		videoIds <- videoIdsQueried[ii]
	}
}

func getVideoIdsAndKeywords(service *youtube.Service) ([]string, string) {
	rand.Seed(time.Now().Unix())
	videoIds := make(chan string, maxQueries)
	keywords := make(chan string, maxQueries)

	for ii := 0; ii < maxQueries; ii++ {
		go searchByKeyword(videoIds, keywords, service)
	}

	var allKeywordsSlice []string
	var allVideoIdsSlice []string

	for ii := 0; ii < maxQueries*maxVideosPerQuery; ii++ {
		allVideoIdsSlice = append(allVideoIdsSlice, <-videoIds)
	}

	for ii := 0; ii < maxQueries; ii++ {
		allKeywordsSlice = append(allKeywordsSlice, <-keywords)
	}

	allKeywords := strings.Join(allKeywordsSlice, ",")

	return allVideoIdsSlice, allKeywords
}
