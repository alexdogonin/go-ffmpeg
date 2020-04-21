package ffmpeg

import "C"

type AudioCodecContextOpt func(*AudioCodecContext)

func WithAudioBitrate(bitrate int) AudioCodecContextOpt {
	return func(context *AudioCodecContext) {
		context.ctype().bit_rate = C.long(bitrate)
	}
}

func WithSampleRate(sampleRate int) AudioCodecContextOpt {
	return func(context *AudioCodecContext) {
		context.ctype().sample_rate = C.int(sampleRate)
	}
}

func WithChannelLayout(layout int) AudioCodecContextOpt {
	return func(context *AudioCodecContext) {
		context.ctype().channel_layout = C.ulong(layout)
	}
}

func WithChannels(channels int) AudioCodecContextOpt {
	return func(context *AudioCodecContext) {
		context.ctype().channels = C.int(channels)
	}
}

func WithSampleFormat(sampleFmt SampleFormat) AudioCodecContextOpt {
	return func(context *AudioCodecContext) {
		context.ctype().sample_fmt = sampleFmt.ctype()
	}
}
