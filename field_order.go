package ffmpeg

//#include <libavcodec/avcodec.h>
import "C"

type FieldOrder int

func (fo FieldOrder) ctype() C.enum_AVFieldOrder {
	return (C.AVFieldOrder)(fo)
}
