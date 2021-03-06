package ffmpeg

//#include <libavutil/samplefmt.h>
import "C"

type SampleFormat int

const (
	SampleFormatNone SampleFormat = iota - 1
	SampleFormatU8                ///< unsigned 8 bits
	SampleFormatS16               ///< signed 16 bits
	SampleFormatS32               ///< signed 32 bits
	SampleFormatFLT               ///< float
	SampleFormatDBL               ///< double
	SampleFormatU8P               ///< unsigned 8 bits, planar
	SampleFormatS16P              ///< signed 16 bits, planar
	SampleFormatS32P              ///< signed 32 bits, planar
	SampleFormatFLTP              ///< float, planar
	SampleFormatDBLP              ///< double, planar
	SampleFormatS64               ///< signed 64 bits
	SampleFormatS64P              ///< signed 64 bits, planar
	SampleFormatNB                ///< Number of sample formats. DO NOT USE if linking dynamically
)

func (sformat SampleFormat) ctype() C.enum_AVSampleFormat {
	return (C.enum_AVSampleFormat)(sformat)
}
