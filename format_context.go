package ffmpeg

/*
#include <libavformat/avformat.h>
#include <libavutil/opt.h>
#include <libavutil/error.h>
#include <libavutil/avstring.h>

char errbuf[AV_ERROR_MAX_STRING_SIZE] = {0};

char *av_err(int errnum) {
	av_strerror(errnum, errbuf, AV_ERROR_MAX_STRING_SIZE);
    return errbuf;
}
*/
import "C"

import (
	"fmt"
	"unsafe"
)

const (
	AVFMT_FLAG_CUSTOM_IO = 0x0080
)

type FormatContext C.struct_AVFormatContext

func NewOutputFormatContext(filename string, oFormat *OutputFormat) (*FormatContext, error) {
	context := (*FormatContext)(unsafe.Pointer(C.avformat_alloc_context()))

	C.av_strlcpy(&(context.ctype().filename[0]), C.CString(filename), C.ulong(len(filename)+1))
	context.url = C.av_strdup(C.CString(filename))

	context.oformat = oFormat.ctype()

	context.priv_data = nil
	if context.oformat.priv_data_size != C.int(0) {
		context.priv_data = C.av_mallocz(C.ulong(context.oformat.priv_data_size))

		if context.oformat.priv_class != nil {
			*((**C.struct_AVClass)(context.priv_data)) = context.oformat.priv_class

			C.av_opt_set_defaults(context.priv_data)
		}
	}

	return context, nil
}

func NewInputFormatContext(input *IOContext, iFormat *InputFormat) (*FormatContext, error) {

	context := (*FormatContext)(unsafe.Pointer(C.avformat_alloc_context()))

	context.ctype().iformat = iFormat.ctype()
	context.ctype().pb = input.ctype()
	context.ctype().flags |= AVFMT_FLAG_CUSTOM_IO

	// C.av_strlcpy(&(context.ctype().filename[0]), C.CString(filename), C.ulong(len(filename)+1))
	// context.ctype().url = C.av_strdup(C.CString(filename))

	context.ctype().duration = C.int64(0x8000000000000000)
	context.ctype().start_time = C.int64(0x8000000000000000)

	context.priv_data = nil
	if context.iformat.priv_data_size != C.int(0) {
		context.priv_data = C.av_mallocz(C.ulong(context.oformat.priv_data_size))

		if context.oformat.priv_class != nil {
			*((**C.struct_AVClass)(context.priv_data)) = context.oformat.priv_class

			C.av_opt_set_defaults(context.priv_data)
		}
	}

	//call static int update_stream_avctx(AVFormatContext *s)
}

func (context *FormatContext) StreamExists() bool {
	return C.avformat_find_stream_info(context.ctype(), nil) >= 0
}

func (context *FormatContext) Release() {
	C.avformat_free_context(context.ctype())
}

func (context *FormatContext) DumpFormat() {
	C.av_dump_format(context.ctype(), 0, &(context.ctype().filename[0]), 1)
}

func (context *FormatContext) OpenIO() error {
	if (context.ctype().oformat.flags & C.AVFMT_NOFILE) != 0 {
		return nil
	}

	ret := C.avio_open(&(context.ctype().pb), &(context.ctype().url), C.AVIO_FLAG_WRITE)
	if ret < 0 {
		return fmt.Errorf("open %q error, %s", context.ctype().filename, C.av_err(C.int(ret)))
	}

	return nil
}

func (context *FormatContext) CloseOutput() {
	C.avio_closep(&(context.ctype().pb))
}

func (context *FormatContext) WriteHeader(opts map[string]string) error {
	// var opt C.struct_AVDictionary
	// p := &opt

	// for k, v := range opts {
	// 	C.av_dict_set(&p, C.CString(k), C.CString(v), 0)
	// }

	ret := C.avformat_write_header(context.ctype(), nil)
	if int(ret) < 0 {
		return fmt.Errorf("write header error, %s", C.av_err(ret))
	}

	return nil
}

func (context *FormatContext) WriteTrailer() error {
	ret := C.av_write_trailer(context.ctype())

	if int(ret) < 0 {
		return fmt.Errorf("write trailer error, %s", C.av_err(ret))
	}

	return nil
}

func (context *FormatContext) WritePacket(packet *Packet) error {
	ret := C.av_interleaved_write_frame(context.ctype(), packet.ctype())
	if int(ret) < 0 {
		return fmt.Errorf("write packet error, %s", C.av_err(ret))
	}

	return nil
}

func (context *FormatContext) Streams() []*Stream {
	streams := make([]*Stream, context.ctype().nb_streams)
	for i := 0; i < len(streams); i++ {
		streams[i] = context.ctype().streams[C.Int(i)]
	}

	return streams
}

func (context *FormatContext) ctype() *C.struct_AVFormatContext {
	return (*C.struct_AVFormatContext)(unsafe.Pointer(context))
}
