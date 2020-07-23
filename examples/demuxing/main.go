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

	streams := formatContext.Streams()

	var videoStream, audioStream *ffmpeg.Stream
	for _, s := range streams {
		if videoStream != nil && audioStream != nil {
			break
		}

		p := s.CodecParameters()

		if audioStream == nil && p.CodecType() == ffmpeg.MediaTypeAudio {
			audioStream = s
			continue
		}

		if videoStream == nil && p.CodecType() == ffmpeg.MediaTypeVideo {
			videoStream = s
			continue
		}
	}

	if audioStream == nil {
		log.Fatal("audio stream not found")
	}

	if videoStream == nil {
		log.Fatal("video stream not found")
	}

	videoCodec, err := ffmpeg.CodecByID(videoStream.CodecParameters().CodecID())
	if err != nil {
		log.Fatal(err)
	}

	audioCodec, err := ffmpeg.CodecByID(audioStream.CodecParameters().CodecID())
	if err != nil {
		log.Fatal(err)
	}

	codecParms := videoStream.CodecParameters()
	videoCodecContext, err := ffmpeg.NewVideoCodecContext(
		videoCodec,
		codecParms.Width(),
		codecParms.Height(),
		ffmpeg.PixelFormat(codecParms.Format()),
	)

	if err != nil {
		log.Fatal(err)
	}

	codecParms = audioStream.CodecParameters()
	audioCodecContext, err := ffmpeg.NewAudioCodecContext(audioCodec, )
}
