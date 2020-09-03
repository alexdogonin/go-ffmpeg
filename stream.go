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
		ExtraData:          C.GoBytes(unsafe.Pointer(s.ctype().codecpar.extradata), s.ctype().codecpar.extradata_size),
		CodecType:          MediaType(s.ctype().codecpar.codec_type),
		CodecID:            CodecID(s.ctype().codecpar.codec_id),
		CodecTag:           uint32(s.ctype().codecpar.codec_tag),
		Format:             int(s.ctype().codecpar.format),
		Bitrate:            int(s.ctype().codecpar.bit_rate),
		BitsPerCodedSample: int(s.ctype().codecpar.bits_per_coded_sample),
		BitsPerRawSample:   int(s.ctype().codecpar.bits_per_raw_sample),
		Profile:            int(s.ctype().codecpar.profile),
		Level:              int(s.ctype().codecpar.level),
		Width:              int(s.ctype().codecpar.width),
		Height:             int(s.ctype().codecpar.height),
		SampleAspectRatio:  Rational{int(s.ctype().codecpar.sample_aspect_ratio.num), int(s.ctype().codecpar.sample_aspect_ratio.den)},
		FieldOrder:         FieldOrder(s.ctype().codecpar.field_order),
		ColorRange:         ColorRange(s.ctype().codecpar.color_range),
		ColorPrimaries:     ColorPrimaries(s.ctype().codecpar.color_primaries),
		ColorTrc:           ColorTransferCharacteristic(s.ctype().codecpar.color_trc),
		ColorSpace:         ColorSpace(s.ctype().codecpar.color_space),
		ChromaLocation:     ChromaLocation(s.ctype().codecpar.chroma_location),
		VideoDelay:         int(s.ctype().codecpar.video_delay),
		ChannelLayout:      ChannelLayout(s.ctype().codecpar.channel_layout),
		Channels:           int(s.ctype().codecpar.channels),
		SampleRate:         int(s.ctype().codecpar.sample_rate),
		BlockAlign:         int(s.ctype().codecpar.block_align),
		FrameSize:          int(s.ctype().codecpar.frame_size),
		InitialPadding:     int(s.ctype().codecpar.initial_padding),
		TrailingPadding:    int(s.ctype().codecpar.trailing_padding),
		SeekPrerol:         int(s.ctype().codecpar.seek_preroll),
	}

	return &p
}

func (s *Stream) TimeBase() Rational {
	return Rational{int(s.ctype().time_base.num), int(s.ctype().time_base.den)}
}

func (s *Stream) ctype() *C.struct_AVStream {
	return (*C.struct_AVStream)(unsafe.Pointer(s))
}
