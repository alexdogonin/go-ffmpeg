package ffmpeg

//#include <libavutil/pixfmt.h>
import "C"

type ColorPrimaries int

const (
	ColorPrimaries_Unspecified = 2
)

func (cp ColorPrimaries) ctype() C.enum_AVColorPrimaries {
	return (C.enum_AVColorPrimaries)(cp)
}
