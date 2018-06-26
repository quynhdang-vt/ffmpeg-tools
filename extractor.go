package main

import (
	"log"
	"os"
	"bufio"
	"strings"
	"strconv"
	"fmt"
	"path/filepath"
	"os/exec"
)

type Extractor struct {
	inputFile string
	timingFile string
	outputFile string
	outputAudioDir string
}

/**
read from timingFile the start, end (in seconds)
formulate a filename, e.g. audio_#_start_end.wav in outputDirectory
Call ffmpeg -i inputfile -ss start -t duration filename
add the filename to outputFile
 */
func (t *Extractor) Extract () {
	file, err := os.Open(t.timingFile)
	if err!=nil {
		log.Fatalf("Failed to open timingFile (%s), err=%s\n", t.timingFile, err)
	}
	defer file.Close()

	outFile, err := os.Create(t.outputFile)
	if err!=nil {
		log.Fatalf("Failed to create outputFile (%s), err=%s\n", t.outputFile, err)
	}
	defer outFile.Close()

	fullOutputDir, err := filepath.Abs(t.outputAudioDir)
	if err!=nil {
		log.Fatalf("Failed to get full path for outputAudioDir(%s), err=%s\n", t.outputAudioDir, err)
	}
	scanner := bufio.NewScanner(file)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		s:=scanner.Text()
		// line should have start, end at the minimum
		ss:=strings.Split(s,",")
		startTime,err:=strconv.Atoi(ss[0])
		if err!=nil {
			log.Fatalf("Failed to get start time, lineNo=%d, err=%s", lineNo, err)
		}
		endTime,err:=strconv.Atoi(ss[1])
		if err!=nil {
			log.Fatalf("Failed to get end time, lineNo=%d, err=%s", lineNo, err)
		}
		duration:=endTime-startTime
		var fname string
		if duration>0 {
			fname = fmt.Sprintf("audio_%05d_%d_%d.wav", lineNo, startTime, endTime)
			destination := filepath.Join(fullOutputDir, fname)
			outFile.WriteString(fmt.Sprintf("%s\n", destination))
			// calling ffmpeg
			ffmpegCmd := exec.Command("ffmpeg", "-hide_banner",
				"-ss", strconv.Itoa(startTime),
				"-t", strconv.Itoa(duration),
				"-i", t.inputFile,
				destination)
			err:=ffmpegCmd.Run()
			if err!=nil {
				log.Fatalf("Failed to split, lineNo=%d, cmd=`ffmpeg -hide_banner -ss %d -t %d -i %s %s`, err=%v\n",
					lineNo, startTime, duration, t.inputFile, destination)
			}
		}
	}
}
