package ffmpeg

//#include <libavcodec/avcodec.h>
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

var (
	EAGAIN = errors.New("EAGAIN")
	EOF    = errors.New("EOF")
)

type VideoCodecContext struct {
	context *C.struct_AVCodecContext

	err error
}

/*NewCodecContext creates new codec context and try to open codec

	Recomended bitrates (https://support.google.com/youtube/answer/1722171?hl=en):
	Type       |   Video Bitrate,    | Video Bitrate,
	           | Standard Frame Rate | High Frame Rate
               |   (24, 25, 30)      |  (48, 50, 60)
	-----------------------------------------------------
	2160p (4k) | 35-45 Mbps          | 53-68 Mbps
	1440p (2k) | 16 Mbps             | 24 Mbps
	1080p      | 8 Mbps              | 12 Mbps
	720p       | 5 Mbps              | 7.5 Mbps
	480p       | 2.5 Mbps            | 4 Mbps
	360p       | 1 Mbps              | 1.5 Mbps
*/
func NewVideoCodecContext(codec *Codec, width, height int, pixFmt PixelFormat, opts ...VideoCodecContextOpt) (*VideoCodecContext, error) {
	context := C.avcodec_alloc_context3((*C.struct_AVCodec)(codec))

	const defaultFramerate = 25

	context.width = C.int(width)
	context.height = C.int(height)
	context.time_base = C.AVRational{C.int(1), C.int(defaultFramerate)}
	context.framerate = C.AVRational{C.int(defaultFramerate), C.int(1)}
	context.gop_size = C.int(10)
	context.pix_fmt = pixFmt.ctype()

	context.bit_rate = C.long(calculateBitrate(int(context.height), int(context.framerate.num)))

	codecCtx := &VideoCodecContext{
		context: context,
	}

	for _, opt := range opts {
		opt(codecCtx)
	}

	if int(context.framerate.num) != defaultFramerate {
		context.bit_rate = C.long(calculateBitrate(int(context.height), int(context.framerate.num)))
	}

	if ok := int(C.avcodec_open2(context, (*C.struct_AVCodec)(codec), nil)) == 0; !ok {
		return nil, errors.New("codec open error")
	}

	return codecCtx, nil
}

func (context *VideoCodecContext) Release() {
	C.avcodec_free_context(&context.context)
}

func (context *VideoCodecContext) SendFrame(frame *VideoFrame) error {
	if context.err != nil {
		return context.err
	}

	if int(C.avcodec_send_frame(context.context, frame.ctype())) != 0 {
		return errors.New("send frame error")
	}

	return nil
}

func (context *VideoCodecContext) ReceivePacket(dest *Packet) bool {
	if context.err != nil {
		return false
	}

	ret := int(C.avcodec_receive_packet(context.context, dest.ctype()))
	if ret == -int(C.EAGAIN) || ret == int(C.AVERROR_EOF) {
		return false
	}

	if ret < 0 {
		context.err = fmt.Errorf("error during encoding (code = %q)", ret)
		return false
	}

	return true
}

func (context *VideoCodecContext) Err() error {
	return context.err
}

func (context *VideoCodecContext) CodecParameters() *CodecParameters {
	parms := &CodecParameters{}

	C.avcodec_parameters_from_context((*C.struct_AVCodecParameters)(unsafe.Pointer(parms)), context.context)

	return parms
}

func calculateBitrate(height, framerate int) int {
	switch framerate {
	case 24, 25, 30:
		switch height {
		case 2160:
			return 40 * 1024 * 1024
		case 1440:
			return 16 * 1024 * 1024
		case 1080:
			return 8 * 1024 * 1024
		case 720:
			return 5 * 1024 * 1024
		case 480:
			return 2.5 * 1024 * 1024
		case 360:
			return 1 * 1024 * 1024
		}
	case 48, 50, 60:
		switch height {
		case 2160:
			return 60 * 1024 * 1024
		case 1440:
			return 24 * 1024 * 1024
		case 1080:
			return 12 * 1024 * 1024
		case 720:
			return 7.5 * 1024 * 1024
		case 480:
			return 4 * 1024 * 1024
		case 360:
			return 1.5 * 1024 * 1024
		}
	}

	width := height * 16 / 9
	coef := 6

	return width * height * framerate / coef
}
