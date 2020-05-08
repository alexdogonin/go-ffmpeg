package ffmpeg

//#include <libavcodec/avcodec.h>
import "C"

import (
	"errors"
	"fmt"
	"math/bits"
	"unsafe"
)

type AudioCodecContext C.struct_AVCodecContext

func NewAudioCodecContext(codec *Codec, bitrate int, rate int, fmt SampleFormat, chLayout ChannelLayout) (*AudioCodecContext, error) {
	c := C.avcodec_alloc_context3((*C.struct_AVCodec)(codec))

	c.bitrate = C.int(bitrate)
	c.channels = bits.OnesCount64(chLayout)
	c.channel_layout = chLayout.ctype()

	context := (*AudioCodecContext)(unsafe.Pointer(c))
	for _, opt := range opts {
		opt(context)
	}

	if ok := int(C.avcodec_open2(c, (*C.struct_AVCodec)(codec), nil)) == 0; !ok {
		return nil, errors.New("codec open error")
	}

	return context, nil
}

func (context *AudioCodecContext) Release() {
	C.avcodec_free_context((**C.struct_AVCodecContext)(unsafe.Pointer(&context)))
}

func (context *AudioCodecContext) ctype() *C.struct_AVCodecContext {
	return (*C.struct_AVCodecContext)(unsafe.Pointer(context))
}
