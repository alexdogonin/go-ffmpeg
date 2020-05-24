package main

import (
	"image"
	"log"

	"os"

	"github.com/alexdogonin/go-ffmpeg"
)

func main() {
	const (
		framerate   = 25
		codecID     = ffmpeg.CodecIDMPEG4
		resultName  = "result.avi"
		videoWidth  = 704
		videoHeight = 576
		durationSec = 8
	)

	f, err := os.Create(resultName)
	if err != nil {
		log.Fatalf("create result file error, %s", err)
	}
	defer f.Close()

	codec, err := ffmpeg.CodecByID(codecID)
	if err != nil {
		log.Fatalf("find codec %d error, %s", codecID, err)
	}

	codecCtx, err := ffmpeg.NewVideoCodecContext(codec, videoWidth, videoHeight, ffmpeg.YUV420P,
		ffmpeg.WithGopSize(10),
	)
	if err != nil {
		log.Fatalf("initialize codec context error, %s", err)
	}
	defer codecCtx.Release()

	frame, err := ffmpeg.NewVideoFrame(videoWidth, videoHeight, ffmpeg.YUV420P)
	if err != nil {
		log.Fatalf("initialize frame error, %s", err)
	}
	defer frame.Release()

	frameImg := image.NewYCbCr(image.Rect(0, 0, videoWidth, videoHeight), image.YCbCrSubsampleRatio420)

	packet, err := ffmpeg.NewPacket()
	if err != nil {
		log.Fatalf("initialize packet error, %s", err)
	}
	defer packet.Release()

	totalFrames := framerate * durationSec
	for i := 0; i < totalFrames; i++ {
		fillFakeImg(frameImg, i)

		if err = frame.FillYCbCr(frameImg); err != nil {
			log.Fatalf("fill frame error, %s", err)
		}

		frame.SetPts(i)

		if err = codecCtx.SendFrame(frame); err != nil {
			log.Fatalf("send frame to encoding error, %s", err)
		}

		for codecCtx.ReceivePacket(packet) {
			if _, err := f.Write(packet.Data()); err != nil {
				log.Fatalf("write file error, %v", err)
			}

			packet.Unref()
		}

		if err = codecCtx.Err(); err != nil {
			log.Fatalf("receive packet error, %s", err)
		}
	}

	// flush
	if err = codecCtx.SendFrame(nil); err != nil {
		log.Fatalf("send frame to encoding error, %s", err)
	}

	for codecCtx.ReceivePacket(packet) {
		if _, err := f.Write(packet.Data()); err != nil {
			log.Fatalf("write file error, %v", err)
		}

		packet.Unref()
	}

	if err = codecCtx.Err(); err != nil {
		log.Fatalf("receive packet error, %s", err)
	}
}

func fillFakeImg(img *image.YCbCr, i int) {
	/* prepare a dummy image */
	/* Y */
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			img.Y[y*img.YStride+x] = uint8(x + y + i*3)
		}
	}

	/* Cb and Cr */
	for y := 0; y < img.Bounds().Max.Y/2; y++ {
		for x := 0; x < img.Bounds().Max.X/2; x++ {
			img.Cb[y*img.CStride+x] = uint8(128 + y + i*2)
			img.Cr[y*img.CStride+x] = uint8(64 + x + i*5)
		}
	}
}
