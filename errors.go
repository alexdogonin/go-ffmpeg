package ffmpeg

//#include <libavutil/error.h>
import "C"
import (
	"unsafe"
)

const (
	bufferSize = 64
)

func Error(errCode C.int) error {
	return err(errCode)
}

type err int

func (err err) Error() string {
	var buff [bufferSize]byte

	ret := C.av_strerror(C.int(err), (*C.char)(unsafe.Pointer(&buff[0])), C.ulong(bufferSize))
	if C.int(ret) != 0 {
		return "unknown error"
	}

	return string(buff[:])
}
