package ffmpeg

//#include <libavutil/pixfmt.h>
import "C"

type ColorSpace int

const (
	ColorSpace_Unspecified = 2
)

func (c ColorSpace) ctype() C.enum_AVColorSpace {
	return (C.enum_AVColorSpace)(c)
}
