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
	defer ioContext.Release()

	inputFormat, err := ffmpeg.ProbeFormat(ioContext)
	if err != nil {
		log.Fatalf("probe format error, %s", err)
	}

	formatContext, err := ffmpeg.NewInputFormatContext(ioContext, inputFormat)
	if err != nil {
		log.Fatalf("create format context error, %s", err)
	}
	// defer formatContext.Close() - duplicated ioContext.Release
	defer formatContext.Release()

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

		if audioStream == nil && p.CodecType == ffmpeg.MediaTypeAudio {
			audioStream = s
			continue
		}

		if videoStream == nil && p.CodecType == ffmpeg.MediaTypeVideo {
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

	codecParms := videoStream.CodecParameters()
	videoCodec, err := ffmpeg.CodecByID(codecParms.CodecID)
	if err != nil {
		log.Fatal(err)
	}

	videoCodecContext, err := ffmpeg.NewVideoCodecContext(
		videoCodec,
		codecParms.Width,
		codecParms.Height,
		ffmpeg.PixelFormat(codecParms.Format),
	)

	if err != nil {
		log.Fatal(err)
	}
	defer videoCodecContext.Release()

	// codecParms = audioStream.CodecParameters()
	// audioCodec, err := ffmpeg.CodecByID(codecParms.CodecID)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// audioCodecContext, err := ffmpeg.NewAudioCodecContext(
	// 	audioCodec,
	// 	codecParms.Bitrate,
	// 	codecParms.SampleRate,
	// 	ffmpeg.SampleFormat(codecParms.Format),
	// 	codecParms.ChannelLayout,
	// )
	// defer audioCodecContext.Release()

	frame, err := ffmpeg.NewVideoFrame(codecParms.Width, codecParms.Height, ffmpeg.PixelFormat(codecParms.Format))
	if err != nil {
		log.Fatal(err)
	}
	defer frame.Release()

	formatContext.ReadPacket()
}
