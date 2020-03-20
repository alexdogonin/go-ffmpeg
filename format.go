package ffmpeg

//#import <libavformat/avformat.h>
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

func (format *OutputFormat) ctype() *C.struct_AVOutputFormat {
	return (*C.struct_AVOutputFormat)(unsafe.Pointer(format))
}
