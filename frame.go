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

var pixelFormatSubsampleRatio = map[PixelFormat]image.YCbCrSubsampleRatio{
	YUV420P:  image.YCbCrSubsampleRatio420,
	YUVJ444P: image.YCbCrSubsampleRatio444,
	YUV444P:  image.YCbCrSubsampleRatio444,
}

type VideoFrame C.struct_AVFrame

func NewVideoFrame(width, height int, pixFmt PixelFormat) (*VideoFrame, error) {
	frame := (*VideoFrame)(unsafe.Pointer(C.av_frame_alloc()))

	frame.format = C.int(pixFmt)
	frame.width = C.int(width)
	frame.height = C.int(height)
	if ret := C.av_frame_get_buffer(frame.ctype(), C.int(1) /*alignment*/); ret < 0 {
		frame.Release()
		return nil, fmt.Errorf("Error allocating avframe buffer. Err: %v", ret)
	}

	return frame, nil
}

func (frame *VideoFrame) SetPts(pts int) {
	frame.ctype().pts = C.long(pts)
}

func (frame *VideoFrame) Release() {
	C.av_frame_free((**C.struct_AVFrame)(unsafe.Pointer(&frame)))
}

func (frame *VideoFrame) MakeWritable() error {
	if 0 != int(C.av_frame_make_writable(frame.ctype())) {
		return errors.New("make writable error")
	}

	return nil
}

func (frame *VideoFrame) FillRGBA(img *image.RGBA) error {
	if PixelFormat(frame.format) != RGBA {
		return errors.New("pixel format of frame is not RGBA")
	}

	if !img.Bounds().Eq(image.Rect(0, 0, int(frame.ctype().width), int(frame.ctype().height))) {
		return errors.New("image resolution not equal frame resolution")
	}

	ok := 0 <= int(C.av_image_fill_arrays(
		&(frame.data[0]),
		&(frame.linesize[0]),
		(*C.uint8_t)(&img.Pix[0]),
		int32(frame.format),
		frame.width,
		frame.height,
		1,
	))

	if !ok {
		return errors.New("Could not fill raw picture buffer")
	}

	return nil
}

func (frame *VideoFrame) FillYCbCr(img *image.YCbCr) error {
	ratio, ok := pixelFormatSubsampleRatio[PixelFormat(frame.format)]
	if !ok {
		return errors.New("pixel format of frame is not YCbCr")
	}

	if ratio != img.SubsampleRatio {
		return errors.New("subsample ration of frame not equal image")
	}

	if !img.Bounds().Eq(image.Rect(0, 0, int(frame.ctype().width), int(frame.ctype().height))) {
		return errors.New("image resolution not equal frame resolution")
	}

	data := make([]uint8, len(img.Y)+len(img.Cb)+len(img.Cr))
	copy(data[:len(img.Y)], img.Y)
	copy(data[len(img.Y):len(img.Y)+len(img.Cb)], img.Cb)
	copy(data[len(img.Y)+len(img.Cb):], img.Cr)

	ok = 0 <= int(C.av_image_fill_arrays(
		&(frame.data[0]),
		&(frame.linesize[0]),
		(*C.uint8_t)(&(data[0])),
		int32(frame.format),
		frame.width,
		frame.height,
		1,
	))

	if !ok {
		return errors.New("Could not fill raw picture buffer")
	}

	return nil
}

func (frame *VideoFrame) ctype() *C.struct_AVFrame {
	return (*C.struct_AVFrame)(unsafe.Pointer(frame))
}
