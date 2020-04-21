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
		codecName   = "h263p"
		resultName  = "result.avi"
		videoWidth  = 640
		videoHeight = 320
		durationSec = 8
	)

	f, err := os.Create(resultName)
	if err != nil {
		log.Fatalf("create result file error, %s", err)
	}
	defer f.Close()

	codec, err := ffmpeg.CodecByName(codecName)
	if err != nil {
		log.Fatalf("find codec %q error, %s", err)
	}

	codecCtx, err := ffmpeg.NewVideoCodecContext(codec, videoWidth, videoHeight, ffmpeg.YUV420P,
		ffmpeg.WithGopSize(10),
	)
	if err != nil {
		log.Fatalf("initialize codec context error, %s", err)
	}
	defer codecCtx.Release()

	frame, err := ffmpeg.NewFrame(videoWidth, videoHeight, ffmpeg.YUV420P)
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
		/* prepare a dummy image */
		/* Y */
		for y := 0; y < frameImg.Bounds().Max.Y; y++ {
			for x := 0; x < frameImg.Bounds().Max.X; x++ {
				frameImg.Y[y*frameImg.YStride+x] = uint8(x + y + i*3)
			}
		}

		/* Cb and Cr */
		for y := 0; y < frameImg.Bounds().Max.Y/2; y++ {
			for x := 0; x < frameImg.Bounds().Max.X/2; x++ {
				frameImg.Cb[y*frameImg.CStride+x] = uint8(128 + y + i*2)
				frameImg.Cr[y*frameImg.CStride+x] = uint8(64 + x + i*5)
			}
		}

		// draw to frame
		// . . .
		if err = frame.FillYCbCr(frameImg); err != nil {
			log.Fatalf("fill frame error, %s", err)
		}

		frame.SetPts(i)

		if err = codecCtx.SendFrame(frame); err != nil {
			log.Fatalf("send frame to encoding error, %s", err)
		}

		for {
			if err = codecCtx.ReceivePacket(packet); err != nil {
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
