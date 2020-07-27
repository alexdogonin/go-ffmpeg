package ffmpeg

//#include <libavutil/pixfmt.h>
import "C"

type ColorSpace int

func (c ColorSpace) ctype() C.enum_AVColorSpace {
	return (C.enum_AVColorSpace)(c)
}
