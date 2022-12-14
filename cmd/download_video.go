package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const (
	maxDurationSeconds    = 10 * 60 * 60
	clipDurationStrFfmpeg = "1"
)

func joinVideos() error {
	files, err := ioutil.ReadDir(pathBase + "/downloads")
	if err != nil {
		log.Fatalf("Error reading files in download directory: %v", err)
	}

	command := "ffmpeg -y -vsync 2 "
	nfiles := 0
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".mkv") {
			command += "-i " + pathBase + "/downloads/" + f.Name() + " "
			nfiles++
		}
	}

	command += "-filter_complex '"

	for ii := 0; ii < nfiles; ii++ {
		command += fmt.Sprintf("[%v:v]scale=1280x720,setdar=1280/720[v%v]; ", ii, ii)
	}
	for ii := 0; ii < nfiles; ii++ {
		command += fmt.Sprintf("[v%v] [%v:a] ", ii, ii)
	}

	command += fmt.Sprintf("concat=n=%v:v=1:a=1 [v] [a]' ", nfiles)

	command += "-map '[v]' -map '[a]' " + pathBase + "/output/" + filename

	log.Printf("Running join command with: %v", command)

	out, e := exec.Command("bash", "-c", command).Output()
	log.Printf(string(out))

	return e

}

func downloadVideo(id string, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Printf("Getting duration for video ID: %v", id)
	command := "youtube-dl --get-duration " + id
	log.Printf("Running command: %v", command)
	out, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		log.Printf("error when getting duration: %v", err)
		return
	}
	outSplit := strings.Split(strings.TrimSpace(string(out)), ":")
	if len(outSplit) == 0 {
		log.Printf("video duration not accessible")
		return
	}
	clipDurationSplit := strings.Split(clipDurationStrFfmpeg, ":")
	units := []string{"s", "m", "h"}
	lenStr := ""
	for ii := 0; ii < len(outSplit); ii++ {
		jj := len(outSplit) - ii - 1
		lenStr = outSplit[jj] + units[ii] + lenStr
	}
	clipDurationStr := ""
	for ii := 0; ii < len(clipDurationSplit); ii++ {
		jj := len(clipDurationSplit) - ii - 1
		clipDurationStr = clipDurationSplit[jj] + units[ii] + clipDurationStr
	}
	duration, err := time.ParseDuration(lenStr)
	if err != nil {
		log.Printf("Error: parsing duration failed: %v", err)
		return
	}
	durationSeconds := duration.Seconds()
	if durationSeconds > maxDurationSeconds {
		log.Printf("Error: duration of video exceeds max duration")
		return
	}
	clipDuration, err := time.ParseDuration(clipDurationStr)
	if err != nil {
		log.Printf("Error: parsing clip duration failed: %v", err)
		return
	}
	clipDurationSeconds := clipDuration.Seconds()
	var start string
	if durationSeconds < clipDurationSeconds {
		start = "00:00:00"
	} else {
		startTimeFloat := rand.Float64() * (durationSeconds - 1)
		startTimeStr := fmt.Sprintf("%fs", startTimeFloat)
		startTime, err := time.ParseDuration(startTimeStr)
		if err != nil {
			log.Printf("Error: could not parse start time: %v", err)
			return
		}
		startHours := startTime.Hours()
		startMinutes := startTime.Minutes()
		startSeconds := startTime.Seconds()
		start = fmt.Sprintf("%02d:%02d:%02d", int(startHours), int(startMinutes-float64(int(startHours)*60)), int(startSeconds-float64(int(startMinutes)*60)))
	}
	filename := pathBase + "/downloads/" + id + ".mkv"
	command = "youtube-dl --youtube-skip-dash-manifest -g 'https://www.youtube.com/watch?v=" + id + "'"
	log.Printf("Running command: %v", command)
	out, err = exec.Command("bash", "-c", command).Output()
	if err != nil {
		log.Printf("Error: getting video and audio streams failed: %v", err)
		return
	}
	outStr := string(out)
	outSplit = strings.Split(outStr, "\n")
	videoUrl := outSplit[0]
	audioUrl := outSplit[1]
	command = "ffmpeg -xerror -ss " + start + " -i '" + videoUrl + "' -ss " + start + " -i '" + audioUrl + "' -map 0:v -map 1:a -ss 0 -t " + clipDurationStrFfmpeg + " -c:v libx264 -c:a aac " + filename
	log.Printf("Running command: %v", command)
	out, err = exec.Command("bash", "-c", command).Output()
	if err != nil {
		log.Printf("error streaming video to file")
		err := os.Remove(filename)
		if err != nil {
			log.Printf("Error deleting file: %v", err)
		}
		return
	}
}
