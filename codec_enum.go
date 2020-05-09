package ffmpeg

//#include <libavcodec/avcodec.h>
import "C"

type CodecID C.enum_AVCodecID

const (
	CodecIDMp2 CodecID = C.AV_CODEC_ID_MP2
	CodecIDMp3 CodecID = C.AV_CODEC_ID_MP3
	CodecIDAAC CodecID = C.AV_CODEC_ID_AAC
)