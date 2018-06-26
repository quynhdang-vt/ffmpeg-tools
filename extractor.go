package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type Extractor struct {
	inputFile      string
	timingFile     string
	outputFile     string
	outputAudioDir string
}

func (t *Extractor) generateTiming() {
	// run the ./segment_audio.sh samples/plumb/Plumbing-Call-Recording.mp3 timinginfo.txt
	segmentCommand := exec.Command("bash", "./segment_audio.sh",
		t.inputFile,
		t.timingFile)
	err:=segmentCommand.Run()
	if err!=nil {
		log.Fatalf("Failed to generate timing info for the segments, err=%v\n", err)
	}
}
func CreateDir(dir string){
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatalf("Failed to create dir %s, err=%s\n", dir, err)
		}
	}
}
/**
read from timingFile the start, end (in seconds)
formulate a filename, e.g. audio_#_start_end.wav in outputDirectory
Call ffmpeg -i inputfile -ss start -t duration filename
add the filename to outputFile
*/
func (t *Extractor) Extract() {
	CreateDir(t.outputAudioDir)
	t.generateTiming()
	file, err := os.Open(t.timingFile)
	if err != nil {
		log.Fatalf("Failed to open timingFile (%s), err=%s\n", t.timingFile, err)
	}
	defer file.Close()

	outFile, err := os.Create(t.outputFile)
	if err != nil {
		log.Fatalf("Failed to create outputFile (%s), err=%s\n", t.outputFile, err)
	}
	defer outFile.Close()

	fullOutputDir, err := filepath.Abs(t.outputAudioDir)
	if err != nil {
		log.Fatalf("Failed to get full path for outputAudioDir(%s), err=%s\n", t.outputAudioDir, err)
	}
	scanner := bufio.NewScanner(file)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		s := scanner.Text()
		// line should have start, end at the minimum
		ss := strings.Split(s, ",")
		startTime, err := strconv.Atoi(ss[0])
		if err != nil {
			log.Fatalf("Failed to get start time, lineNo=%d, err=%s", lineNo, err)
		}
		endTime, err := strconv.Atoi(ss[1])
		if err != nil {
			log.Fatalf("Failed to get end time, lineNo=%d, err=%s", lineNo, err)
		}
		duration := endTime - startTime
		var fname string
		if duration > 0 {
			fname = fmt.Sprintf("audio_%05d_%d_%d.wav", lineNo, startTime, endTime)
			destination := filepath.Join(fullOutputDir, fname)
			outFile.WriteString(fmt.Sprintf("%s\n", destination))
			// calling ffmpeg
			ffmpegCmd := exec.Command("ffmpeg", "-hide_banner",
				"-ss", strconv.Itoa(startTime),
				"-t", strconv.Itoa(duration),
				"-i", t.inputFile,
				destination,
					"-y")
			err := ffmpegCmd.Run()
			if err != nil {
				log.Fatalf("Failed to split, lineNo=%d, cmd=`ffmpeg -hide_banner -ss %d -t %d -i %s %s`, err=%s\n",
					lineNo, startTime, duration, t.inputFile, destination, err)
			}
			fmt.Printf("Line %d processed..%s\n", lineNo, destination)
		} else {
			fmt.Printf("Line %d skipped (duration=%d)\n", lineNo, duration)
		}
	}
}
