package main

/*
#cgo pkg-config: libavcodec
#cgo pkg-config: libswscale
#cgo pkg-config: libavformat
#cgo pkg-config: libavutil

#include <libavformat/avio.h>
void do_read(AVIOContext *context) {
	context->read_packet(context->opaque, context->buffer, 0);
}

*/
import "C"
import (
	"fmt"
	"strings"
	"unsafe"
)

func main() {
	r := strings.NewReader("hello!")
	context := NewIOContext(r, nil)
	defer context.Release()
	ioCtx := context.ctype()
	b := (*[20]byte)(unsafe.Pointer(ioCtx.buffer))
	fmt.Printf("buffer init state: %s\n", b)
	// ((ioCtx.read_packet)(C.void (*)(void*, int8_t*, int)))(ioCtx.opaque, ioCtx.buffer, C.int(0))
	C.do_read(ioCtx)
	fmt.Printf("buffer mod state: %s\n", b)

}
