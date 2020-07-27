package ffmpeg

//#include <libavformat/avformat.h>
import "C"

import (
	"errors"
	"unsafe"
)

type Stream C.struct_AVStream

// NewStream creates new stream and add them to formatContext.streams
func NewStream(formatContext *FormatContext) (*Stream, error) {
	stream := (*Stream)(unsafe.Pointer(C.avformat_new_stream(formatContext.ctype(), nil)))

	if stream == nil {
		return nil, errors.New("create stream error")
	}

	stream.id = C.int(formatContext.nb_streams - 1)

	return stream, nil
}

func (s *Stream) SetCodecParameters(parms *CodecParameters) {
	*(s.ctype().codecpar) = *(parms.ctype())
}

func (s *Stream) CodecParameters() *CodecParameters {
	p := CodecParameters{
		// (*C.uint8_t)(unsafe.Pointer(&ExtraData[0])): extradata,
		// C.int(len(ExtraData)):                       extradata_size,
		ExtraData:          C.GoBytes(s.ctype().codecpar.extradata, s.ctype().codecpar.extradata_size),
		CodecType:          s.ctype().codecpar.codec_type,
		CodecID:            s.ctype().codecpar.codec_id,
		CodecTag:           s.ctype().codecpar.codec_tag,
		Format:             s.ctype().codecpar.format,
		Bitrate:            s.ctype().codecpar.bit_rate,
		BitsPerCodedSample: s.ctype().codecpar.bits_per_coded_sample,
		BitsPerRawSample:   s.ctype().codecpar.bits_per_raw_sample,
		Profile:            s.ctype().codecpar.profile,
		Level:              s.ctype().codecpar.level,
		Width:              s.ctype().codecpar.width,
		Height:             s.ctype().codecpar.height,
		SampleAspectRatio:  s.ctype().codecpar.sample_aspect_ratio,
		FieldOrder:         s.ctype().codecpar.field_order,
		ColorRange:         s.ctype().codecpar.color_range,
		ColorPrimaries:     s.ctype().codecpar.color_primaries,
		ColorTrc:           s.ctype().codecpar.color_trc,
		ColorSpace:         s.ctype().codecpar.color_space,
		ChromaLocation:     s.ctype().codecpar.chroma_location,
		VideoDelay:         s.ctype().codecpar.video_delay,
		ChannelLayout:      s.ctype().codecpar.channel_layout,
		Channels:           s.ctype().codecpar.channels,
		SampleRate:         s.ctype().codecpar.sample_rate,
		BlockAlign:         s.ctype().codecpar.block_align,
		FrameSize:          s.ctype().codecpar.frame_size,
		InitialPadding:     s.ctype().codecpar.initial_padding,
		TrailingPadding:    s.ctype().codecpar.trailing_padding,
		SeekPrerol:         s.ctype().codecpar.seek_preroll,
	}

	return &p
}

func (s *Stream) TimeBase() Rational {
	return Rational{int(s.ctype().time_base.num), int(s.ctype().time_base.den)}
}

func (s *Stream) ctype() *C.struct_AVStream {
	return (*C.struct_AVStream)(unsafe.Pointer(s))
}
