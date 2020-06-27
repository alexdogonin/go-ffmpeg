package ffmpeg

//#include <libavformat/avio.h>
//#include <libavutil/error.h>
import "C"
import (
	"fmt"
	"unsafe"
)

type IOContext C.struct_AVIOContext

func NewIOContext(filename string) (*IOContext, error) {
	var context *IOContext

	ret := C.avio_open(&(context.ctype()), C.String(filename), C.AVIO_FLAG_WRITE)
	if ret < 0 {
		return nil, fmt.Errorf("open %q error, %s", filename, C.av_err(C.int(ret)))
	}

	return context, nil
}

func (context *IOContext) Release() {
	C.avio_close(context.ctype())
}

func (context *IOContext) ctype() *C.struct_AVIOContext {
	return (*C.struct_AVIOContext)(unsafe.Pointer(context))
}
