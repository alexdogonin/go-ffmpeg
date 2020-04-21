package main

import (
	"log"

	"time"

	"image"

	"github.com/alexdogonin/go-ffmpeg"
)

func main() {
	const (
		pixelFormat     = ffmpeg.YUV420P
		width           = 640
		height          = 480
		fileName        = "result.avi"
		duration        = 10 * time.Second
		defaulFramerate = 25
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

	codecContext, err := ffmpeg.NewVideoCodecContext(codec, width, height, pixelFormat)
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

	frame, err := ffmpeg.NewVideoFrame(width, height, pixelFormat)
	if err != nil {
		log.Fatalf("init frame error, %s", err)
	}
	defer frame.Release()

	formatContext.DumpFormat()

	if err = formatContext.OpenOutput(); err != nil {
		log.Fatalf("open output error, %s", err)
	}
	defer formatContext.CloseOutput()

	if err = formatContext.WriteHeader(nil); err != nil {
		log.Fatalf("write header error, %s", err)
	}

	packet, err := ffmpeg.NewPacket()
	if err != nil {
		log.Fatalf("init packet error, %s", err)
	}
	defer packet.Release()

	img := image.NewYCbCr(image.Rect(0, 0, width, height), image.YCbCrSubsampleRatio420)
	// write frame here
	framesCount := int(duration.Seconds()) * defaulFramerate
	for i := 0; i < framesCount; i++ {
		fillFakeImg(img, i)

		if err = frame.FillYCbCr(img); err != nil {
			log.Fatalf("fill image error, %s", err)
		}

		frame.SetPts(i)

		if err = codecContext.SendFrame(frame); err != nil {
			log.Fatalf("encode frame error, %s", err)
		}

		if err = codecContext.ReceivePacket(packet); err != nil {
			log.Fatalf("receive packet error, %s", err)
		}

		if err = formatContext.WritePacket(packet); err != nil {
			log.Fatalf("write packet to output error, %s", err)
		}

		packet.Unref()
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
