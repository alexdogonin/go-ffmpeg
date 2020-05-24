package ffmpeg

//#include <libavutil/frame.h>
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

type AudioFrame C.struct_Frame

func NewAudioFrame(samples int, sampleFmt SampleFormat, channelLayout ChannelLayout) (*AudioFrame, error) {
	f := C.av_frame_alloc()
	f.nb_samples = C.int(samples)
	f.format = C.int(sampleFmt.ctype())
	f.channel_layout = channelLayout.ctype()

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

func (frame *AudioFrame) Write(data []byte) (int, error) {
	if int(frame.ctype().linesize[0]) < len(data) {
		return 0, errors.New("frame buffer less than writable data")
	}

	C.memcpy(unsafe.Pointer(frame.ctype().data[0]), unsafe.Pointer(&(data[0])), C.ulong(len(data)))

	return len(data), nil
}

func (frame *AudioFrame) SetPts(pts int) {
	frame.ctype().pts = C.long(pts)
}

func (frame *AudioFrame) ctype() *C.struct_AVFrame {
	return (*C.struct_AVFrame)(unsafe.Pointer(frame))
}
