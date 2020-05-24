package ffmpeg

import "C"

type Rational struct {
	Num,
	Den int
}

func (rational Rational) ctype() C.struct_AVRational {
	return C.struct_AVRational{C.int(rational.Num), C.int(rational.Den)}
}
