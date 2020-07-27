package ffmpeg

//#include <libavutil/avutil.h>
import "C"

const (
	MediaTypeUnknown    = int(C.AVMEDIA_TYPE_UNKNOWN)
	MediaTypeVideo      = int(C.AVMEDIA_TYPE_VIDEO)
	MediaTypeAudio      = int(C.AVMEDIA_TYPE_AUDIO)
	MediaTypeData       = int(C.AVMEDIA_TYPE_DATA)
	MediaTypeSubtitle   = int(C.AVMEDIA_TYPE_SUBTITLE)
	MediaTypeAttachment = int(C.AVMEDIA_TYPE_ATTACHMENT)
	MediaTypeNB         = int(C.AVMEDIA_TYPE_NB)
)

type MediaType int

func (mt MediaType) ctype() C.enum_AVMediaType {
	return (C.enum_AVMediaType)(mt)
}
