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
	stream.time_base = C.struct_AVRational{1, 25}

	return stream, nil
}

func (s *Stream) SetCodecParameters(parms *CodecParameters) {
	*(s.ctype().codecpar) = *(parms.ctype())
}

func (s *Stream) ctype() *C.struct_AVStream {
	return (*C.struct_AVStream)(unsafe.Pointer(s))
}
