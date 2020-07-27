package ffmpeg

//#include <libavcodec/avcodec.h>
import "C"
import "unsafe"

type CodecParameters struct {
	// General type of the encoded data.
	CodecType MediaType
	// Specific type of the encoded data (the codec used).
	CodecID CodecID
	// Additional information about the codec (corresponds to the AVI FOURCC).
	CodecTag uint32
	// Extra binary data needed for initializing the decoder, codec-dependent.
	// Must be allocated with av_malloc() and will be freed by
	// avcodec_parameters_free(). The allocated size of extradata must be at
	// least extradata_size + AV_INPUT_BUFFER_PADDING_SIZE, with the padding
	// bytes zeroed.
	ExtraData []byte
	// - video: the pixel format, the value corresponds to enum AVPixelFormat.
	// - audio: the sample format, the value corresponds to enum AVSampleFormat.
	Format int
	// The average bitrate of the encoded data (in bits per second).
	Bitrate int
	// The number of bits per sample in the codedwords.
	// This is basically the bitrate per sample. It is mandatory for a bunch of
	// formats to actually decode them. It's the number of bits for one sample in
	// the actual coded bitstream.
	// This could be for example 4 for ADPCM
	// For PCM formats this matches bits_per_raw_sample
	// Can be 0
	BitsPerCodedSample int
	// This is the number of valid bits in each output sample. If the
	// sample format has more bits, the least significant bits are additional
	// padding bits, which are always 0. Use right shifts to reduce the sample
	// to its actual size. For example, audio formats with 24 bit samples will
	// have bits_per_raw_sample set to 24, and format set to AV_SAMPLE_FMT_S32.
	// To get the original sample use "(int32_t)sample >> 8"."
	// For ADPCM this might be 12 or 16 or similar
	// Can be 0
	BitsPerRawSample int
	// Codec-specific bitstream restrictions that the stream conforms to.
	Profile int
	Level   int
	// Video only. The dimensions of the video frame in pixels.
	Width, Height int
	// Video only. The aspect ratio (width / height) which a single pixel
	// should have when displayed.
	// When the aspect ratio is unknown / undefined, the numerator should be
	// set to 0 (the denominator may have any value).
	SampleAspectRatio Rational
	// Video only. The order of the fields in interlaced video.
	FieldOrder FieldOrder
	// Video only. Additional colorspace characteristics.
	ColorRange     ColorRange
	ColorPrimaries ColorPrimaries
	ColorTrc       ColorTransferCharacteristic
	ColorSpace     ColorSpace
	ChromaLocation ChromaLocation
	// Video only. Number of delayed frames.
	VideoDelay int
	// Audio only. The channel layout bitmask. May be 0 if the channel layout is
	// unknown or unspecified, otherwise the number of bits set must be equal to
	// the channels field.
	ChannelLayout ChannelLayout
	// Audio only. The number of audio channels.
	Channels int
	// Audio only. The number of audio samples per second.
	SampleRate int
	// Audio only. The number of bytes per coded audio frame, required by some
	// formats.
	// Corresponds to nBlockAlign in WAVEFORMATEX.
	BlockAlign int
	// Audio only. Audio frame size, if known. Required by some formats to be static.
	FrameSize int
	// Audio only. The amount of padding (in samples) inserted by the encoder at
	// the beginning of the audio. I.e. this number of leading decoded samples
	// must be discarded by the caller to get the original audio without leading
	// padding.
	InitialPadding int
	// Audio only. The amount of padding (in samples) appended by the encoder to
	// the end of the audio. I.e. this number of decoded samples must be
	// discarded by the caller from the end of the stream to get the original
	// audio without any trailing padding.
	TrailingPadding int
	// Audio only. Number of samples to skip after a discontinuity.
	SeekPrerol int
}

func (parms *CodecParameters) ctype() *C.struct_AVCodecParameters {
	p := C.struct_AVCodecParameters{
		codec_type:            parms.CodecType.ctype(),
		codec_id:              parms.CodecID.ctype(),
		codec_tag:             C.uint32_t(parms.CodecTag),
		extradata:             C.CBytes(unsafe.Pointer(&parms.ExtraData)),
		extradata_size:        C.int(len(parms.ExtraData)),
		format:                C.int(parms.Format),
		bit_rate:              C.int64_t(parms.Bitrate),
		bits_per_coded_sample: C.int(parms.BitsPerCodedSample),
		bits_per_raw_sample:   C.int(parms.BitsPerRawSample),
		profile:               C.int(parms.Profile),
		level:                 C.int(parms.Level),
		width:                 C.int(parms.Width),
		height:                C.int(parms.Height),
		sample_aspect_ratio:   parms.SampleAspectRatio.ctype(),
		field_order:           parms.FieldOrder.ctype(),
		color_range:           parms.ColorRange.ctype(),
		color_primaries:       parms.ColorPrimaries.ctype(),
		color_trc:             parms.ColorTrc.ctype(),
		color_space:           parms.ColorSpace.ctype(),
		chroma_location:       parms.ChromaLocation.ctype(),
		video_delay:           C.int(parms.VideoDelay),
		channel_layout:        parms.ChannelLayout.ctype(),
		channels:              C.int(parms.Channels),
		sample_rate:           C.int(parms.SampleRate),
		block_align:           C.int(parms.BlockAlign),
		frame_size:            C.int(parms.FrameSize),
		initial_padding:       C.int(parms.InitialPadding),
		trailing_padding:      C.int(parms.TrailingPadding),
		seek_preroll:          C.int(parms.SeekPrerol),
	}

	return &p
}
