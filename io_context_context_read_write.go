package ffmpeg

import "C"
import (
	"log"
	"unsafe"
)

//export Read
func Read(opaque unsafe.Pointer, buff unsafe.Pointer, size int) int {
	id := *(*int)(opaque)
	b := (*[bufferSize]byte)(buff)
	p, _ := contexts.Load(id)

	context := p.(*IOContext)
	n, err := context.r.Read((b)[:])
	if err != nil {
		log.Println(err)
		return -1
	}

	return n
}

//export Write
func Write(opaque unsafe.Pointer, buff unsafe.Pointer, size int) int {
	id := *(*int)(opaque)
	b := (*[bufferSize]byte)(buff)
	p, _ := contexts.Load(id)

	context := p.(*IOContext)
	n, err := context.w.Write((b)[:])
	if err != nil {
		log.Println(err)
		return -1
	}

	return n
}
