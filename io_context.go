package main

/*

#include <libavformat/avio.h>
#include "io_context.h"

extern int Read(void *opaque, unsigned char *buff,  int size);
extern int Write(void *opaque, unsigned char *buff, int size);

typedef struct io_context {
	int (*read)(void *opaque, uint8_t *buf, int buf_size);
	int (*write)(void *opaque, uint8_t *buf, int buf_size);

} io_context;
*/
import "C"
import (
	"io"
	"unsafe"
)

type ioContextOpaque struct {
	r  io.Reader
	wr io.Writer
}

type IOContext struct {
	// opaque *ioContextOpaque
	c *C.struct_AVIOContext
}

func NewIOContext(source io.Reader, destination io.Writer) *IOContext {
	opaque := &ioContextOpaque{source, destination}

	const writable = 1
	c := C.avio_alloc_context(
		(*C.uchar)(C.av_malloc(1024)),
		1024,
		writable,
		unsafe.Pointer(opaque),
		(*[0]byte)(C.Read),
		(*[0]byte)(C.Write),
		nil,
	)
	return &IOContext{
		c: c,
	}
}

func (context *IOContext) Release() {
	C.avio_context_free(&context.c)
}
