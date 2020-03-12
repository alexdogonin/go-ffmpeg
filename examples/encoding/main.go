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

	codec, err := ffmpeg.CodecByName(codecName)
	if err != nil {
		log.Fatalf("find codec %q error, %s", err)
	}

	codecCtx, err := ffmpeg.NewCodecContext(codec, videoWidth, videoHeight, ffmpeg.YUV420P)
	if err != nil {
		log.Fatalf("initialize codec context error, %s", err)
	}
	defer codecCtx.Release()

	frame, err := ffmpeg.NewFrame(videoWidth, videoHeight, ffmpeg.YUV420P)
	if err != nil {
		log.Fatalf("initialize frame error, %s", err)
	}
	defer frame.Release()
	frameImg := image.NewRGBA(image.Rect(0, 0, videoWidth, videoHeight))

	packet, err := ffmpeg.NewPacket()
	if err != nil {
		log.Fatalf("initialize packet error, %s", err)
	}
	defer packet.Release()

	totalFrames := framerate * durationSec
	for i := 0; i < totalFrames; i++ {
		// draw to frame
		// . . .

		if err = codecCtx.SendFrame(frame); err != nil {
			log.Fatalf("send frame to encoding error, %s", err)
		}

		for {
			if err = codecCtx.ReceivePacket(packet); err != nil {
				log.Fatalf("receive packet error, %s", err)
			}

		}
	}
}
