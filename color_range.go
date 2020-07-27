package ffmpeg

//#include <libavutil/pixfmt.h>
import "C"

type ColorRange int

func (cr ColorRange) ctype() C.enum_AVColorRange {
	return (C.enum_AVColorRange)(cr)
}
