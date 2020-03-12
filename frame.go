package ffmpeg

//#include <libavcodec/avcodec.h>
//#include <libavutil/imgutils.h>
import "C"

import (
	"errors"
	"fmt"
	"image"
	"unsafe"
)

type Frame C.struct_AVFrame

func NewFrame(width, height int, pixFmt PixelFormat) (*Frame, error) {
	frame := (*Frame)(unsafe.Pointer(C.av_frame_alloc()))

	frame.format = C.int(pixFmt)
	frame.width = C.int(width)
	frame.height = C.int(height)
	if ret := C.av_frame_get_buffer(frame.ctype(), C.int(1) /*alignment*/); ret < 0 {
		frame.Release()
		return nil, fmt.Errorf("Error allocating avframe buffer. Err: %v", ret)
	}

	return frame, nil
}

func (frame *Frame) SetPts(pts int) {
	frame.ctype().pts = C.long(pts)
}

func (frame *Frame) Release() {
	C.av_frame_free((**C.struct_AVFrame)(unsafe.Pointer(&frame)))
}

func (frame *Frame) MakeWritable() error {
	if 0 != int(C.av_frame_make_writable(frame.ctype())) {
		return errors.New("make writable error")
	}

	return nil
}

func (frame *Frame) FillRGBA(img *image.RGBA) error {
	if !img.Bounds().Eq(image.Rect(0, 0, int(frame.ctype().width), int(frame.ctype().height))) {
		return errors.New("image resolution not equal frame resolution")
	}

	ok := 0 <= int(C.av_image_fill_arrays(
		&(frame.ctype().data[0]),
		&(frame.ctype().linesize[0]),
		(*C.uint8_t)(&img.Pix[0]),
		C.AV_PIX_FMT_RGBA,
		frame.ctype().width,
		frame.ctype().height,
		1,
	))

	if !ok {
		return errors.New("Could not fill raw picture buffer")
	}

	return nil
}

func (frame *Frame) ctype() *C.struct_AVFrame {
	return (*C.struct_AVFrame)(unsafe.Pointer(frame))
}
