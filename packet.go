package ffmpeg

//#include <libavcodec/avcodec.h>
import "C"

import (
	"errors"
	"unsafe"
)

type Packet C.struct_AVPacket

func NewPacket() (*Packet, error) {
	packet := (*Packet)(C.av_packet_alloc())

	if packet == nil {
		return nil, errors.New("create packet error")
	}

	return packet, nil
}

func (packet *Packet) Release() {
	C.av_packet_free((**C.struct_AVPacket)(unsafe.Pointer(&packet)))
}

func (packet *Packet) Data() []byte {
	return C.GoBytes(unsafe.Pointer(packet.ctype().data), C.int(packet.ctype().size))
}

func (packet *Packet) Unref() {
	C.av_packet_unref(packet.ctype())
}

func (packet *Packet) ctype() *C.struct_AVPacket {
	return (*C.struct_AVPacket)(unsafe.Pointer(packet))
}
