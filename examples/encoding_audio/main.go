package main

import (
	"encoding/binary"
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
		channelsLayout = ffmpeg.ChannelLayoutMono
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

	data := make([]byte, 0)

	for i := 0; i < frames; i++ {
		data = fillFakeFrame(context.SamplesPerFrame(), i, sampleRate, data[:0])

		if _, err = frame.Write(data); err != nil {
			log.Fatal(err)
		}

		if err = context.SendFrame(frame); err != nil {
			log.Fatal(err)
		}

		for context.ReceivePacket(packet) {
			if _, err := f.Write(packet.Data()); err != nil {
				log.Fatalf("write file error, %v", err)
			}

			packet.Unref()
		}

		if err = context.Err(); err != nil {
			log.Fatalf("receive packet error, %s", err)
		}

	}

	// flush
	if err = context.SendFrame(nil); err != nil {
		log.Fatal(err)
	}

	for context.ReceivePacket(packet) {
		if _, err := f.Write(packet.Data()); err != nil {
			log.Fatalf("write file error, %v", err)
		}

		packet.Unref()
	}

	if err = context.Err(); err != nil {
		log.Fatalf("receive packet error, %s", err)
	}
}

func fillFakeFrame(samples, i int, sampleRate int, dest []byte) []byte {
	timePerSample := 2 * math.Pi * 440 / float64(sampleRate)

	for j := 0; j < samples; j++ {
		t := timePerSample * float64(i*samples+j)

		const maxLevel = math.MaxUint32 / 2
		level := uint32(math.Sin(t) * maxLevel)

		const formatBytes = 4
		bytes := make([]byte, formatBytes)

		binary.LittleEndian.PutUint32(bytes, level)

		dest = append(dest, bytes...)
	}

	return dest
}
