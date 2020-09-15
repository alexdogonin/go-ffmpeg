package ffmpeg

//#include <libavcodec/avcodec.h>
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

var (
	EAGAIN = errors.New("EAGAIN")
	EOF    = errors.New("EOF")
)

type VideoCodecContext struct {
	context *C.struct_AVCodecContext

	err error
}

/*NewCodecContext creates new codec context and try to open codec

	Recomended bitrates (https://support.google.com/youtube/answer/1722171?hl=en):
	Type       |   Video Bitrate,    | Video Bitrate,
	           | Standard Frame Rate | High Frame Rate
               |   (24, 25, 30)      |  (48, 50, 60)
	-----------------------------------------------------
	2160p (4k) | 35-45 Mbps          | 53-68 Mbps
	1440p (2k) | 16 Mbps             | 24 Mbps
	1080p      | 8 Mbps              | 12 Mbps
	720p       | 5 Mbps              | 7.5 Mbps
	480p       | 2.5 Mbps            | 4 Mbps
	360p       | 1 Mbps              | 1.5 Mbps
*/
func NewVideoCodecContext(codec *Codec, width, height int, pixFmt PixelFormat, opts ...VideoCodecContextOpt) (*VideoCodecContext, error) {
	context := C.avcodec_alloc_context3((*C.struct_AVCodec)(codec))

	const defaultFramerate = 25

	context.width = C.int(width)
	context.height = C.int(height)
	context.time_base = C.AVRational{C.int(1), C.int(defaultFramerate)}
	context.framerate = C.AVRational{C.int(defaultFramerate), C.int(1)}
	context.gop_size = C.int(10)
	context.pix_fmt = pixFmt.ctype()

	context.bit_rate = C.long(calculateBitrate(int(context.height), int(context.framerate.num)))

	codecCtx := &VideoCodecContext{
		context: context,
	}

	for _, opt := range opts {
		opt(codecCtx)
	}

	if int(context.framerate.num) != defaultFramerate {
		context.bit_rate = C.long(calculateBitrate(int(context.height), int(context.framerate.num)))
	}

	if ok := int(C.avcodec_open2(context, (*C.struct_AVCodec)(codec), nil)) == 0; !ok {
		return nil, errors.New("codec open error")
	}

	return codecCtx, nil
}

func (context *VideoCodecContext) Release() {
	C.avcodec_free_context(&context.context)
}

func (context *VideoCodecContext) SendFrame(frame *VideoFrame) error {
	if context.err != nil {
		return context.err
	}

	if int(C.avcodec_send_frame(context.context, frame.ctype())) != 0 {
		return errors.New("send frame error")
	}

	return nil
}

func (context *VideoCodecContext) ReceivePacket(dest *Packet) bool {
	if context.err != nil {
		return false
	}

	ret := int(C.avcodec_receive_packet(context.context, dest.ctype()))
	if ret == -int(C.EAGAIN) || ret == int(C.AVERROR_EOF) {
		return false
	}

	if ret < 0 {
		context.err = fmt.Errorf("error during encoding (code = %q)", ret)
		return false
	}

	return true
}

func (context *VideoCodecContext) Err() error {
	return context.err
}

func (context *VideoCodecContext) CodecParameters() *CodecParameters {
	parms := &CodecParameters{
		ColorPrimaries:     ColorPrimaries_Unspecified,
		ColorTrc:           ColorTransferCharacteristic_Unspecified,
		ColorSpace:         ColorSpace_Unspecified,
		SampleAspectRatio:  Rational{0, 1},
		CodecType:          MediaType(context.context.codec_type),
		CodecID:            CodecID(context.context.codec_id),
		CodecTag:           uint32(context.context.codec_tag),
		Bitrate:            int(context.context.bit_rate),
		BitsPerCodedSample: int(context.context.bits_per_coded_sample),
		BitsPerRawSample:   int(context.context.bits_per_raw_sample),
		Profile:            int(context.context.profile),
		Level:              int(context.context.level),
		Format:             -1,
	}

	switch MediaType(parms.CodecType) {
	case MediaTypeVideo:
		parms.Format = int(context.context.pix_fmt)
		parms.Width = int(context.context.width)
		parms.Height = int(context.context.height)
		parms.FieldOrder = FieldOrder(context.context.field_order)
		parms.ColorRange = ColorRange(context.context.color_range)
		parms.ColorPrimaries = ColorPrimaries(context.context.color_primaries)
		parms.ColorTrc = ColorTransferCharacteristic(context.context.color_trc)
		parms.ColorSpace = ColorSpace(context.context.colorspace)
		parms.ChromaLocation = ChromaLocation(context.context.chroma_sample_location)
		parms.SampleAspectRatio = Rational{int(context.context.sample_aspect_ratio.num), int(context.context.sample_aspect_ratio.den)}
		parms.VideoDelay = int(context.context.has_b_frames)
	case MediaTypeAudio:
		parms.Format = int(context.context.sample_fmt)
		parms.ChannelLayout = ChannelLayout(context.context.channel_layout)
		parms.Channels = int(context.context.channels)
		parms.SampleRate = int(context.context.sample_rate)
		parms.BlockAlign = int(context.context.block_align)
		parms.FrameSize = int(context.context.frame_size)
		parms.InitialPadding = int(context.context.initial_padding)
		parms.TrailingPadding = int(context.context.trailing_padding)
		parms.SeekPreroll = int(context.context.seek_preroll)
	case MediaTypeSubtitle:
		parms.Width = int(context.context.width)
		parms.Height = int(context.context.height)
	}

	if context.context.extradata != nil {
		parms.ExtraData = C.GoBytes(unsafe.Pointer(context.context.extradata), context.context.extradata_size)
	}

	return parms
}

func calculateBitrate(height, framerate int) int {
	switch framerate {
	case 24, 25, 30:
		switch height {
		case 2160:
			return 40 * 1024 * 1024
		case 1440:
			return 16 * 1024 * 1024
		case 1080:
			return 8 * 1024 * 1024
		case 720:
			return 5 * 1024 * 1024
		case 480:
			return 2.5 * 1024 * 1024
		case 360:
			return 1 * 1024 * 1024
		}
	case 48, 50, 60:
		switch height {
		case 2160:
			return 60 * 1024 * 1024
		case 1440:
			return 24 * 1024 * 1024
		case 1080:
			return 12 * 1024 * 1024
		case 720:
			return 7.5 * 1024 * 1024
		case 480:
			return 4 * 1024 * 1024
		case 360:
			return 1.5 * 1024 * 1024
		}
	}

	width := height * 16 / 9
	coef := 6

	return width * height * framerate / coef
}
