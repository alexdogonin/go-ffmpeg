package ffmpeg

//#include <libavutil/pixfmt.h>
import "C"

const (
	RGBA     PixelFormat = C.AV_PIX_FMT_RGBA
	YUV420P  PixelFormat = C.AV_PIX_FMT_YUV420P
	YUVJ444P PixelFormat = C.AV_PIX_FMT_YUVJ444P
	YUV444P  PixelFormat = C.AV_PIX_FMT_YUV444P
)

type PixelFormat C.enum_AVPixelFormat

func (fmt PixelFormat) ctype() C.enum_AVPixelFormat {
	return (C.enum_AVPixelFormat)(fmt)
}
