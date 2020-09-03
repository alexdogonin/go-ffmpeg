package ffmpeg

//#include <libavutil/avutil.h>
import "C"

const (
	MediaTypeUnknown    = MediaType(C.AVMEDIA_TYPE_UNKNOWN)
	MediaTypeVideo      = MediaType(C.AVMEDIA_TYPE_VIDEO)
	MediaTypeAudio      = MediaType(C.AVMEDIA_TYPE_AUDIO)
	MediaTypeData       = MediaType(C.AVMEDIA_TYPE_DATA)
	MediaTypeSubtitle   = MediaType(C.AVMEDIA_TYPE_SUBTITLE)
	MediaTypeAttachment = MediaType(C.AVMEDIA_TYPE_ATTACHMENT)
	MediaTypeNB         = MediaType(C.AVMEDIA_TYPE_NB)
)

type MediaType int

func (mt MediaType) ctype() C.enum_AVMediaType {
	return (C.enum_AVMediaType)(mt)
}
