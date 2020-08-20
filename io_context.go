package ffmpeg

/*

#include <libavformat/avio.h>

extern int Read(void *opaque, unsigned char *buff,  int size);
extern int Write(void *opaque, unsigned char *buff, int size);
*/
import "C"
import (
	"io"
	"sync"
	"sync/atomic"
	"unsafe"
)

const (
	bufferSize = 1024
)

var contexts = sync.Map{}
var contextsCounter int32

type IOContext struct {
	id int
	r  io.Reader
	w  io.Writer
	c  *C.struct_AVIOContext
}

func NewIOContext(source io.Reader, destination io.Writer) (*IOContext, error) {
	id := int(atomic.AddInt32(&contextsCounter, 1))

	context := &IOContext{
		id: id,
		r:  source,
		w:  destination,
	}

	const writable = 1
	c := C.avio_alloc_context(
		(*C.uchar)(C.av_malloc(bufferSize)),
		bufferSize,
		writable,
		unsafe.Pointer(&context.id),
		(*[0]byte)(C.Read),
		(*[0]byte)(C.Write),
		nil,
	)

	context.c = c

	contexts.Store(context.id, context)

	return context, nil
}

func (context *IOContext) ctype() *C.struct_AVIOContext {
	return context.c
}

func (context *IOContext) Release() {
	C.avio_context_free(&context.c)

	contexts.Delete(context.id)
}
