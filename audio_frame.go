package ffmpeg

//#include <libavutil/frame.h>
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

type AudioFrame C.struct_Frame

func NewAudioFrame(size int, sampleFmt SampleFormat, channelLayout uint64) (*AudioFrame, error) {
	f := C.av_frame_alloc()
	f.nb_samples = C.int(size)
	f.format = C.int(sampleFmt.ctype())
	f.channel_layout = C.ulong(channelLayout)

	frame := (*AudioFrame)(unsafe.Pointer(f))
	if ret := C.av_frame_get_buffer(frame.ctype(), C.int(1) /*alignment*/); ret < 0 {
		frame.Release()
		return nil, fmt.Errorf("Error allocating avframe buffer. Err: %v", ret)
	}

	return frame, nil
}

func (frame *AudioFrame) Release() {
	C.av_frame_free((**C.struct_AVFrame)(unsafe.Pointer(&frame)))
}

func (frame *AudioFrame) MakeWritable() error {
	if 0 != int(C.av_frame_make_writable(frame.ctype())) {
		return errors.New("make writable error")
	}

	return nil
}

func (frame *AudioFrame) ctype() *C.struct_AVFrame {
	return (*C.struct_AVFrame)(unsafe.Pointer(frame))
}
