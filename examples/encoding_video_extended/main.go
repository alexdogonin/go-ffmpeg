package main

import (
	"log"

	"github.com/alexdogonin/go-ffmpeg"
)

func main() {
	const (
		pixelFormat = ffmpeg.YUV420P
		width       = 640
		height      = 480
		fileName    = "result.avi"
	)

	format, err := ffmpeg.GuessFileNameFormat(fileName)
	if err != nil {
		log.Fatalf("init format error, %s", err)
	}

	formatContext, err := ffmpeg.NewFormatContext(fileName, format)
	if err != nil {
		log.Fatalf("init format context error, %s", err)
	}
	defer formatContext.Release()

	codec, err := ffmpeg.CodecByID(format.VideoCodec())
	if err != nil {
		log.Fatalf("find encodec error, %s", err)
	}

	codecContext, err := ffmpeg.NewCodecContext(codec, width, height, pixelFormat)
	if err != nil {
		log.Fatalf("init codec context error")
	}
	defer codecContext.Release()

	stream, err := ffmpeg.NewStream(formatContext)
	if err != nil {
		log.Fatalf("init stream error, %s", err)
	}

	codecParms := codecContext.CodecParameters()
	stream.SetCodecParameters(codecParms)

	frame, err := ffmpeg.NewFrame(width, height, pixelFormat)
	if err != nil {
		log.Fatalf("init frame error, %s", err)
	}

	formatContext.DumpFormat()

	if err = formatContext.OpenOutput(); err != nil {
		log.Fatalf("open output error, %s", err)
	}

	

}
