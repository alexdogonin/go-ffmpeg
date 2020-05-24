package main

//#cgo pkg-config: libavcodec
//#include <libavformat/avformat.h>
//#include <libavcodec/avcodec.h>
//#include <libavcodec/avcodec.h>
import "C"

import (
	"encoding/binary"
	"image"
	"log"
	"math"
	"time"

	"github.com/alexdogonin/go-ffmpeg"
)

func main() {
	const (
		filename       = "result.avi"
		pixelFormat    = ffmpeg.YUV420P
		subsampleRatio = image.YCbCrSubsampleRatio420
		width          = 640
		height         = 480
		framerate      = 25
		bitrate        = 64000
		sampleFormat   = ffmpeg.SampleFormatS32P
		sampleRate     = 44100
		channelsLayout = ffmpeg.ChannelLayoutMono
		duration       = 5 * time.Second
	)

	format, err := ffmpeg.GuessFormat("avi")
	if err != nil {
		log.Fatal(err)
	}

	formatContext, err := ffmpeg.NewFormatContext(filename, format)
	if err != nil {
		log.Fatal(err)
	}
	defer formatContext.Release()

	videoCodec, err := ffmpeg.CodecByID(format.VideoCodec())
	if err != nil {
		log.Fatal(err)
	}

	audioCodec, err := ffmpeg.CodecByID(format.AudioCodec())
	if err != nil {
		log.Fatal(err)
	}

	videoCodecContext, err := ffmpeg.NewVideoCodecContext(videoCodec, width, height, pixelFormat)
	if err != nil {
		log.Fatal(err)
	}
	defer videoCodecContext.Release()

	audioCodecContext, err := ffmpeg.NewAudioCodecContext(audioCodec, bitrate, sampleRate, sampleFormat, channelsLayout)
	if err != nil {
		log.Fatal(err)
	}
	defer audioCodecContext.Release()

	videoStream, err := ffmpeg.NewStream(formatContext)
	if err != nil {
		log.Fatal(err)
	}
	codecParms := videoCodecContext.CodecParameters()
	videoStream.SetCodecParameters(codecParms)

	audioStream, err := ffmpeg.NewStream(formatContext)
	if err != nil {
		log.Fatal(err)
	}
	codecParms = audioCodecContext.CodecParameters()
	audioStream.SetCodecParameters(codecParms)

	if err = formatContext.OpenOutput(); err != nil {
		log.Fatalf("open output error, %s", err)
	}
	defer formatContext.CloseOutput()

	packet, err := ffmpeg.NewPacket()
	if err != nil {
		log.Fatal(err)
	}
	defer packet.Release()

	videoFrame, err := ffmpeg.NewVideoFrame(width, height, pixelFormat)
	if err != nil {
		log.Fatal(err)
	}
	defer videoFrame.Release()

	audioFrame, err := ffmpeg.NewAudioFrame(audioCodecContext.SamplesPerFrame(), sampleFormat, channelsLayout)
	if err != nil {
		log.Fatal(err)
	}
	defer audioFrame.Release()

	if err = formatContext.WriteHeader(nil); err != nil {
		log.Fatalf("write header error, %s", err)
	}

	formatContext.DumpFormat()

	videoFramesCount := int(duration.Seconds()) * framerate
	img := image.NewYCbCr(image.Rect(0, 0, width, height), subsampleRatio)

	packet.SetStream(0)

	for i := 0; i < videoFramesCount; i++ {
		fillFakeImg(img, i)

		if err = videoFrame.FillYCbCr(img); err != nil {
			log.Fatalf("fill image error, %s", err)
		}

		videoFrame.SetPts(i)

		if err = videoCodecContext.SendFrame(videoFrame); err != nil {
			log.Fatalf("encode frame error, %s", err)
		}

		for videoCodecContext.ReceivePacket(packet) {
			packet.RescaleTimestamp(ffmpeg.Rational{1, 25}, videoStream.TimeBase())

			if err = formatContext.WritePacket(packet); err != nil {
				log.Fatalf("write packet to output error, %s", err)
			}

			packet.Unref()
		}

		if err = videoCodecContext.Err(); err != nil {
			log.Fatalf("receive packet error, %s", err)
		}

	}

	if err = videoCodecContext.SendFrame(nil); err != nil {
		log.Fatalf("encode frame error, %s", err)
	}

	for videoCodecContext.ReceivePacket(packet) {
		packet.RescaleTimestamp(ffmpeg.Rational{1, 25}, videoStream.TimeBase())

		if err = formatContext.WritePacket(packet); err != nil {
			log.Fatalf("write packet to output error, %s", err)
		}

		packet.Unref()
	}

	if err = videoCodecContext.Err(); err != nil {
		log.Fatalf("receive packet error, %s", err)
	}

	audioFramesCount := int(duration.Seconds()) * sampleRate / audioCodecContext.SamplesPerFrame()
	timePerSample := 2 * math.Pi * 440.0 / sampleRate

	for i := 0; i < audioFramesCount; i++ {
		data := make([]byte, 0)

		for j := 0; j < audioCodecContext.SamplesPerFrame(); j++ {
			t := timePerSample * float64(i*audioCodecContext.SamplesPerFrame()+j)

			const maxLevel = math.MaxUint32 / 2
			level := uint32(math.Sin(t) * maxLevel)

			const formatBytes = 4
			bytes := make([]byte, formatBytes)

			binary.LittleEndian.PutUint32(bytes, level)

			data = append(data, bytes...)
		}

		if _, err = audioFrame.Write(data); err != nil {
			log.Fatal(err)
		}

		audioFrame.SetPts(i * audioCodecContext.SamplesPerFrame())

		if err = audioCodecContext.SendFrame(audioFrame); err != nil {
			log.Fatal(err)
		}

		for audioCodecContext.ReceivePacket(packet) {
			packet.SetStream(1)

			packet.RescaleTimestamp(ffmpeg.Rational{1, 44100}, audioStream.TimeBase())

			if err = formatContext.WritePacket(packet); err != nil {
				log.Fatalf("write packet to output error, %s", err)
			}

			packet.Unref()
		}

		if err = audioCodecContext.Err(); err != nil {
			log.Fatalf("receive packet error, %s", err)
		}
	}

	if err = audioCodecContext.SendFrame(nil); err != nil {
		log.Fatal(err)
	}

	for audioCodecContext.ReceivePacket(packet) {
		packet.SetStream(1)

		packet.RescaleTimestamp(ffmpeg.Rational{1, 44100}, audioStream.TimeBase())

		if err = formatContext.WritePacket(packet); err != nil {
			log.Fatalf("write packet to output error, %s", err)
		}

		packet.Unref()
	}

	if err = audioCodecContext.Err(); err != nil {
		log.Fatalf("receive packet error, %s", err)
	}

	if err = formatContext.WriteTrailer(); err != nil {
		log.Fatalf("write trailer error, %s", err)
	}

}

func fillFakeImg(img *image.YCbCr, frameInd int) {
	/* prepare a dummy image */
	/* Y */
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			img.Y[y*img.YStride+x] = uint8(x + y + frameInd*3)
		}
	}

	/* Cb and Cr */
	for y := 0; y < img.Bounds().Max.Y/2; y++ {
		for x := 0; x < img.Bounds().Max.X/2; x++ {
			img.Cb[y*img.CStride+x] = uint8(128 + y + frameInd*2)
			img.Cr[y*img.CStride+x] = uint8(64 + x + frameInd*5)
		}
	}
}
