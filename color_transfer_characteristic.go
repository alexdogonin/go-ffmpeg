package ffmpeg

//#include <libavutil/pixfmt.h>
import "C"

type ColorTransferCharacteristic int

func (c ColorTransferCharacteristic) ctype() C.enum_AVColorTransferCharacteristic {
	return (C.enum_AVColorTransferCharacteristic)(c)
}
