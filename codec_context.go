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

type CodecContext C.struct_AVCodecContext

type CodecContextOpt func(codecCtx *CodecContext)

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
func NewCodecContext(codec *Codec, opts ...CodecContextOpt) (*CodecContext, error) {
	context := C.avcodec_alloc_context3((*C.struct_AVCodec)(codec))

	context.width = C.int(1280)
	context.height = C.int(720)
	context.time_base = C.AVRational{C.int(1), C.int(25)}
	context.framerate = C.AVRational{C.int(25), C.int(1)}
	context.gop_size = C.int(10)
	context.pix_fmt = YUV420P.ctype()

	context.bit_rate = C.long(-1)

	codecCtx := (*CodecContext)(unsafe.Pointer(context))
	for _, opt := range opts {
		opt(codecCtx)
	}

	if int(context.bit_rate) == -1 {
		context.bit_rate = C.long(calculateBitrate(int(context.height), int(context.framerate.num)))
	}

	if ok := int(C.avcodec_open2(context, (*C.struct_AVCodec)(codec), nil)) == 0; !ok {
		return nil, errors.New("codec open error")
	}

	return codecCtx, nil
}

func (context *CodecContext) Release() {
	C.avcodec_free_context((**C.struct_AVCodecContext)(unsafe.Pointer(&context)))
}

func (context *CodecContext) SendFrame(frame *Frame) error {
	if int(C.avcodec_send_frame(context.ctype(), frame.ctype())) != 0 {
		return errors.New("send frame error")
	}

	return nil
}

func (context *CodecContext) ReceivePacket(dest *Packet) error {
	ret := int(C.avcodec_receive_packet(context.ctype(), dest.ctype()))
	if ret == -int(C.EAGAIN) {
		return EAGAIN
	}

	if ret == int(C.AVERROR_EOF) {
		return EOF
	}

	if ret < 0 {
		return fmt.Errorf("error during encoding (code = %q)", ret)
	}

	return nil
}

func (context *CodecContext) ctype() *C.struct_AVCodecContext {
	return (*C.struct_AVCodecContext)(unsafe.Pointer(context))
}

func WithBitrate(bitrate int) CodecContextOpt {
	return func(codecCtx *CodecContext) {
		codecCtx.ctype().bit_rate = C.long(bitrate)
	}
}

func WithResolution(width, height int) CodecContextOpt {
	return func(codecCtx *CodecContext) {
		codecCtx.ctype().width = C.int(width)
		codecCtx.ctype().height = C.int(height)
	}
}

func WithFramerate(framerate int) CodecContextOpt {
	return func(codecCtx *CodecContext) {
		codecCtx.ctype().time_base = C.AVRational{C.int(1), C.int(framerate)}
		codecCtx.ctype().framerate = C.AVRational{C.int(framerate), C.int(1)}
	}
}

func WithPixFmt(pixFmt PixelFormat) CodecContextOpt {
	return func(codecCtx *CodecContext) {
		codecCtx.ctype().pix_fmt = pixFmt.ctype()
	}
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
