package ffmpeg

//#include <libavcodec/avcodec.h>
import "C"

const (
	YUV420P PixelFormat = C.AV_PIX_FMT_YUV420P
	RGBA    PixelFormat = C.AV_PIX_FMT_RGBA
)

type PixelFormat C.enum_AVPixelFormat

func (fmt PixelFormat) ctype() C.enum_AVPixelFormat {
	return (C.enum_AVPixelFormat)(fmt)
}
