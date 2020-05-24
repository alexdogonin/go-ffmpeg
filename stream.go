package ffmpeg

//#include <libavformat/avformat.h>
import "C"

import (
	"errors"
	"unsafe"
)

type Stream C.struct_AVStream

// NewStream creates new stream and add them to formatContext.streams
func NewStream(formatContext *FormatContext) (*Stream, error) {
	stream := (*Stream)(unsafe.Pointer(C.avformat_new_stream(formatContext.ctype(), nil)))

	if stream == nil {
		return nil, errors.New("create stream error")
	}

	stream.id = C.int(formatContext.nb_streams - 1)

	return stream, nil
}

func (s *Stream) SetCodecParameters(parms *CodecParameters) {
	*(s.ctype().codecpar) = *(parms.ctype())
}

func (s *Stream) TimeBase() Rational {
	return Rational{int(s.ctype().time_base.num), int(s.ctype().time_base.den)}
}

func (s *Stream) ctype() *C.struct_AVStream {
	return (*C.struct_AVStream)(unsafe.Pointer(s))
}
