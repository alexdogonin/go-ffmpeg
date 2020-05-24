package ffmpeg

//#include <libavcodec/avcodec.h>
import "C"

type VideoCodecContextOpt func(*VideoCodecContext)

func WithVideoBitrate(bitrate int) VideoCodecContextOpt {
	return func(codecCtx *VideoCodecContext) {
		codecCtx.context.bit_rate = C.long(bitrate)
	}
}

func WithResolution(width, height int) VideoCodecContextOpt {
	return func(codecCtx *VideoCodecContext) {
		codecCtx.context.width = C.int(width)
		codecCtx.context.height = C.int(height)
	}
}

func WithFramerate(framerate int) VideoCodecContextOpt {
	return func(codecCtx *VideoCodecContext) {
		codecCtx.context.time_base = C.AVRational{C.int(1), C.int(framerate)}
		codecCtx.context.framerate = C.AVRational{C.int(framerate), C.int(1)}
	}
}

func WithPixFmt(pixFmt PixelFormat) VideoCodecContextOpt {
	return func(codecCtx *VideoCodecContext) {
		codecCtx.context.pix_fmt = pixFmt.ctype()
	}
}

func WithGopSize(gopSize int) VideoCodecContextOpt {
	return func(codecCtx *VideoCodecContext) {
		codecCtx.context.gop_size = C.int(gopSize)
	}
}
