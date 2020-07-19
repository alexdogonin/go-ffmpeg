package main

import (
	"log"

	"github.com/alexdogonin/go-ffmpeg"
)

func main() {
	const (
		inputFilename              = "input.mpeg"
		outputAudioFilename        = "result.mp3"
		outputImagePatternFilename = "result%d.jpeg"
	)

	ioContext, err := ffmpeg.NewIOContext(inputFilename)
	if err != nil {
		log.Fatalf("create io context error, %s", err)
	}

	inputFormat, err := ffmpeg.ProbeFormat(ioContext)
	if err != nil {
		log.Fatalf("probe format error, %s", err)
	}

	formatContext, err := ffmpeg.NewInputFormatContext(ioContext, inputFormat)
	if err != nil {
		log.Fatalf("create format context error, %s", err)
	}

	if !formatContext.StreamExists() {
		log.Fatal("streams not found")
	}

	
}
