package ffmpeg

//#include <libavcodec/avcodec.h>
import "C"
import "unsafe"

type CodecParameters C.struct_AVCodecParameters

func (parms *CodecParameters) ctype() *C.struct_AVCodecParameters {
	return (*C.struct_AVCodecParameters)(unsafe.Pointer(parms))
}
