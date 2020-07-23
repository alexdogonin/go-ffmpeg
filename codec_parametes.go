package ffmpeg

//#include <libavcodec/avcodec.h>
import "C"
import "unsafe"

type CodecParameters C.struct_AVCodecParameters

func (c *CodecParameters) CodecID() CodecID {
	return CodecID(c.ctype().codec_id)
}

func (c *CodecParameters) CodecType() MediaType {
	return MediaType(c.ctype().codec_type)
}

func (c *CodecParameters) Width() int {
	return int(c.ctype().width)
}

func (c *CodecParameters) Height() int {
	return int(c.ctype().height)
}

func (c *CodecParameters) Format() int {
	return int(c.ctype().format)
}

func (c *CodecParameters) Bitrate() int {
	return int(c.ctype().bit_rate)
}

func (parms *CodecParameters) ctype() *C.struct_AVCodecParameters {
	return (*C.struct_AVCodecParameters)(unsafe.Pointer(parms))
}
