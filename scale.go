package ffmpeg

//#include <libswscale/swscale.h>
import "C"

import (
	"errors"
	"unsafe"
)

type ScaleContext C.struct_SwsContext

func NewScaleContext(srcWidth, srcHeight int, srcPixFmt PixelFormat, dstWidth, dstHeight int, dstPixFmt PixelFormat) (*ScaleContext, error) {
	convertContext := (*ScaleContext)(C.sws_getContext(
		C.int(srcWidth),
		C.int(srcHeight),
		RGBA.ctype(),
		C.int(dstWidth),
		C.int(dstHeight),
		YUV420P.ctype(),
		C.SWS_GAUSS,
		(*C.struct_SwsFilter)(nil),
		(*C.struct_SwsFilter)(nil),
		(*C.double)(nil),
	))

	if convertContext == nil {
		return nil, errors.New("create convetr context error")
	}

	return convertContext, nil
}

func (scaleContext *ScaleContext) Release() {
	C.sws_freeContext(scaleContext.ctype())
}

func (scaleContext *ScaleContext) Scale(src, dst *Frame) {
	C.sws_scale(
		scaleContext.ctype(),
		&(src.ctype().data[0]),
		&(src.ctype().linesize[0]),
		C.int(0),
		src.ctype().height,
		&(dst.ctype().data[0]),
		&(dst.ctype().linesize[0]),
	)
}

func (scaleContext *ScaleContext) ctype() *C.struct_SwsContext {
	return (*C.struct_SwsContext)(unsafe.Pointer(scaleContext))
}
