package ffmpeg

//#include <libavformat/avformat.h>
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

type OutputFormat C.struct_AVOutputFormat

func GuessFileNameFormat(filename string) (*OutputFormat, error) {
	format := (*OutputFormat)(unsafe.Pointer(C.av_guess_format(nil, C.CString(filename), nil)))

	if format == nil {
		return nil, errors.New(fmt.Sprintf("format not for filename %q", filename))
	}

	return format, nil
}

func GuessFormat(name string) (*OutputFormat, error) {
	format := (*OutputFormat)(unsafe.Pointer(C.av_guess_format(C.CString(name), nil, nil)))

	if format == nil {
		return nil, errors.New(fmt.Sprintf("format %q not found", name))
	}

	return format, nil
}

func (format *OutputFormat) VideoCodec() CodecID {
	return CodecID(format.ctype().video_codec)
}

func (format *OutputFormat) AudioCodec() CodecID {
	return CodecID(format.ctype().audio_codec)
}

func (format *OutputFormat) ctype() *C.struct_AVOutputFormat {
	return (*C.struct_AVOutputFormat)(unsafe.Pointer(format))
}

type InputFormat C.struct_AVInputFormat

func ProbeFormat(ioCtx *IOContext) (*InputFormat, error) {
	var inputFmt *InputFormat

	ret := C.av_probe_input_buffer2(ioCtx.ctype(), (**C.struct_AVInputFormat)(unsafe.Pointer(&inputFmt)), (*C.char)(unsafe.Pointer(&[0]byte{})), nil, 0, 0)
	if int(ret) < 0 {
		return nil, fmt.Errorf("probe input error, %s", Error(ret))
	}

	return inputFmt, nil
}

func (format *InputFormat) ctype() *C.struct_AVInputFormat {
	return (*C.struct_AVInputFormat)(unsafe.Pointer(format))
}
