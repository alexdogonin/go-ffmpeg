package main

import (
	"log"
	"math"
	"os"
	"time"

	"github.com/alexdogonin/go-ffmpeg"
)

func main() {
	const (
		filename       = "result.mp3"
		codecID        = ffmpeg.CodecIDMp3
		bitrate        = 64000
		sampleFmt      = ffmpeg.SampleFormatS32P
		sampleRate     = 44100
		channels       = 2
		channelsLayout = ffmpeg.ChannelLayoutStereo
		duration       = 5 * time.Second
	)

	codec, err := ffmpeg.CodecByID(codecID)
	if err != nil {
		log.Fatal("init codec error")
	}

	context, err := ffmpeg.NewAudioCodecContext(codec, bitrate, sampleRate, sampleFmt, channelsLayout)
	if err != nil {
		log.Fatal("init context error")
	}
	defer context.Release()

	frame, err := ffmpeg.NewAudioFrame(context.SamplesPerFrame(), sampleFmt, channelsLayout)
	if err != nil {
		log.Fatal("init frame error")
	}
	defer frame.Release()

	packet, err := ffmpeg.NewPacket()
	if err != nil {
		log.Fatal("init packet error")
	}
	defer packet.Release()

	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	frames := int(duration.Seconds()) * sampleRate / context.SamplesPerFrame()

	t := 0.
	tIncr := 2 * math.Pi * 440.0 / sampleRate

	for i := 0; i < frames; i++ {
		// frame.MakeWritable()

		data := make([]byte, 0)
		for j := 0; j < context.SamplesPerFrame(); j++ {
			left := int32(math.Sin(t) * 100)
			data = append(data, uint8(left>>24), uint8(left>>16), uint8(left>>8), uint8(left))
			// data = append(data, uint8(left>>8), uint8(left))

			// right := left
			// data = append(data, uint8(right>>24), uint8(right>>16), uint8(right>>8), uint8(right))
			// data = append(data, uint8(right>>8), uint8(right))

			t += tIncr
		}

		if _, err = frame.Write(data); err != nil {
			log.Fatal(err)
		}

		if err = context.SendFrame(frame); err != nil {
			log.Fatal(err)
		}

		for {
			if err = context.ReceivePacket(packet); err != nil {
				if err == ffmpeg.EAGAIN || err == ffmpeg.EOF {
					break
				}

				log.Fatalf("receive packet error, %s", err)
			}

			if _, err := f.Write(packet.Data()); err != nil {
				log.Fatalf("write file error, %v", err)
			}

			packet.Unref()
		}

	}
}
