package main

import (
	"image"
	"log"

	"github.com/alexdogonin/go-ffmpeg"
)

func main() {
	const (
		framerate   = 25
		codecName   = "h263p"
		resultName  = "result.avi"
		videoWidth  = 640
		videoHeight = 320
		durationSec = 8
	)

	frameImg := image.NewRGBA(image.Rect(0, 0, videoWidth, videoHeight))

	codec, err := ffmpeg.CodecByName(codecName)
	if err != nil {
		log.Fatalf("find codec %q error, %s", err)
	}

	codecCtx, err := ffmpeg.NewCodecContext(codec, videoWidth, videoHeight)
	if err != nil {
		log.Fatalf("initialize codec context error, %s", err)
	}
	defer codecodecCtx.Release()

	frame, err := ffmpeg.NewFrame(videoWidth, videoHeight, ffmpeg.YUV420P)
}
