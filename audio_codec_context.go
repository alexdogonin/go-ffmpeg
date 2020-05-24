package ffmpeg

//#include <libavcodec/avcodec.h>
import "C"

import (
	"errors"
	"fmt"
	"math/bits"
	"unsafe"
)

type AudioCodecContext struct {
	context *C.struct_AVCodecContext

	err error
}

func NewAudioCodecContext(codec *Codec, bitrate int, rate int, fmt SampleFormat, chLayout ChannelLayout) (*AudioCodecContext, error) {
	c := C.avcodec_alloc_context3((*C.struct_AVCodec)(codec))

	c.bit_rate = C.long(bitrate)
	c.channels = C.int(bits.OnesCount64(uint64(chLayout)))
	c.channel_layout = chLayout.ctype()
	c.sample_rate = C.int(rate)
	c.sample_fmt = fmt.ctype()

	context := &AudioCodecContext{
		context: c,
	}
	// for _, opt := range opts {
	// 	opt(context)
	// }

	if ok := int(C.avcodec_open2(c, (*C.struct_AVCodec)(codec), nil)) == 0; !ok {
		return nil, errors.New("codec open error")
	}

	return context, nil
}

func (context *AudioCodecContext) Release() {
	C.avcodec_free_context(&context.context)
}

func (context *AudioCodecContext) SendFrame(frame *AudioFrame) error {
	if context.err != nil {
		return context.err
	}

	if int(C.avcodec_send_frame(context.context, frame.ctype())) != 0 {
		return errors.New("send frame error")
	}

	return nil
}

func (context *AudioCodecContext) ReceivePacket(dest *Packet) bool {
	if context.err != nil {
		return false
	}

	ret := int(C.avcodec_receive_packet(context.context, dest.ctype()))
	if ret == -int(C.EAGAIN) || ret == int(C.AVERROR_EOF) {
		return false
	}

	if ret < 0 {
		context.err = fmt.Errorf("error during encoding (code = %q)", ret)
		return false
	}

	return true
}

func (context *AudioCodecContext) Err() error {
	return context.err
}

func (context *AudioCodecContext) SamplesPerFrame() int {
	return int(context.context.frame_size)
}

func (context *AudioCodecContext) CodecParameters() *CodecParameters {
	parms := &CodecParameters{}

	C.avcodec_parameters_from_context((*C.struct_AVCodecParameters)(unsafe.Pointer(parms)), context.context)

	return parms
}

func (context *AudioCodecContext) ctype() *C.struct_AVCodecContext {
	return (*C.struct_AVCodecContext)(unsafe.Pointer(context))
}
