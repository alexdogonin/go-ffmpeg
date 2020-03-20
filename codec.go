package ffmpeg

//#include <libavcodec/avcodec.h>
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

type CodecID C.enum_AVCodecID

type Codec C.struct_AVCodec

func CodecByName(name string) (*Codec, error) {
	codec := (*Codec)(C.avcodec_find_encoder_by_name(C.CString(name)))

	if codec == nil {
		return nil, errors.New("init codec error")
	}

	return codec, nil
}

func CodecByID(codecID CodecID) (*Codec, error) {
	codec := (*Codec)(C.avcodec_find_encoder((C.enum_AVCodecID)(codecID)))
	if codec == nil {
		return nil, fmt.Errorf("couldn't find encoder for codec id %d", codecID)
	}

	return codec, nil
}

func (c *Codec) ctype() *C.struct_AVCodec {
	return (*C.struct_AVCodec)(unsafe.Pointer(c))
}
