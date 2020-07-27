package ffmpeg

//#include <libavutil/pixfmt.h>
import "C"

type ChromaLocation int

func (l ChromaLocation) ctype() C.enum_AVChromaLocation {
	return (C.enum_AVChromaLocation)(l)
}
