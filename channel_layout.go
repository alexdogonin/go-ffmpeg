package ffmpeg

import "C"

type Channel uint64

const (
	ChannelFrontLeft           Channel = 0x00000001
	ChannelFrontRight          Channel = 0x00000002
	ChannelFrontCenter         Channel = 0x00000004
	ChannelLowFrequency        Channel = 0x00000008
	ChannelBackLeft            Channel = 0x00000010
	ChannelBackRight           Channel = 0x00000020
	ChannelFrontLeftOfCenter   Channel = 0x00000040
	ChannelFrontRightOfCenter  Channel = 0x00000080
	ChannelBackCenter          Channel = 0x00000100
	ChannelSideLeft            Channel = 0x00000200
	ChannelSideRight           Channel = 0x00000400
	ChannelTopCenter           Channel = 0x00000800
	ChannelTopFrontLeft        Channel = 0x00001000
	ChannelTopFrontCenter      Channel = 0x00002000
	ChannelTopFrontRight       Channel = 0x00004000
	ChannelTopBackLeft         Channel = 0x00008000
	ChannelTopBackCenter       Channel = 0x00010000
	ChannelTopBackRight        Channel = 0x00020000
	ChannelStereoLeft          Channel = 0x20000000 ///< Stereo downmix.
	ChannelStereoRight         Channel = 0x40000000 ///< See ChannelStereoLeft.
	ChannelWideLeft            Channel = 0x0000000080000000
	ChannelWideRight           Channel = 0x0000000100000000
	ChannelSurroundDirectLeft  Channel = 0x0000000200000000
	ChannelSurroundDirectRight Channel = 0x0000000400000000
	ChannelLowFrequency2       Channel = 0x0000000800000000
)

type ChannelLayout uint64

const (
	ChannelLayoutMono            ChannelLayout = ChannelLayout(ChannelFrontCenter)
	ChannelLayoutStereo          ChannelLayout = ChannelLayout(ChannelFrontLeft | ChannelFrontRight)
	ChannelLayout2Point1         ChannelLayout = ChannelLayout(ChannelFrontLeft | ChannelFrontRight | ChannelLowFrequency)
	ChannelLayout2_1             ChannelLayout = ChannelLayout(ChannelFrontLeft | ChannelFrontRight | ChannelBackCenter)
	ChannelLayoutSurround        ChannelLayout = ChannelLayout(ChannelFrontLeft | ChannelFrontRight | ChannelFrontCenter)
	ChannelLayout3Point1         ChannelLayout = ChannelLayout(ChannelFrontLeft | ChannelFrontRight | ChannelFrontCenter | ChannelLowFrequency)
	ChannelLayout4Point0         ChannelLayout = ChannelLayout(ChannelFrontLeft | ChannelFrontRight | ChannelFrontCenter | ChannelBackCenter)
	ChannelLayout4Point1         ChannelLayout = ChannelLayout(ChannelFrontLeft | ChannelFrontRight | ChannelFrontCenter | ChannelBackCenter | ChannelLowFrequency)
	ChannelLayout2_2             ChannelLayout = ChannelLayout(ChannelFrontLeft | ChannelFrontRight | ChannelSideLeft | ChannelSideRight)
	ChannelLayoutQuad            ChannelLayout = ChannelLayout(ChannelFrontLeft | ChannelFrontRight | ChannelBackLeft | ChannelBackRight)
	ChannelLayout5Point0         ChannelLayout = ChannelLayout(ChannelFrontLeft | ChannelFrontRight | ChannelFrontCenter | ChannelSideLeft | ChannelSideRight)
	ChannelLayout5Point1         ChannelLayout = ChannelLayout(ChannelFrontLeft | ChannelFrontRight | ChannelFrontCenter | ChannelSideLeft | ChannelSideRight | ChannelLowFrequency)
	ChannelLayout5Point0Back     ChannelLayout = ChannelLayout(ChannelFrontLeft | ChannelFrontRight | ChannelFrontCenter | ChannelBackLeft | ChannelBackRight)
	ChannelLayout5Point1Back     ChannelLayout = ChannelLayout(ChannelFrontLeft | ChannelFrontRight | ChannelFrontCenter | ChannelBackLeft | ChannelBackRight | ChannelLowFrequency)
	ChannelLayout6Point0         ChannelLayout = ChannelLayout(ChannelFrontLeft | ChannelFrontRight | ChannelFrontCenter | ChannelSideLeft | ChannelSideRight | ChannelBackCenter)
	ChannelLayout6Point0Front    ChannelLayout = ChannelLayout(ChannelFrontLeft | ChannelFrontRight | ChannelSideLeft | ChannelSideRight | ChannelFrontLeftOfCenter | ChannelFrontRightOfCenter)
	ChannelLayoutHexagonal       ChannelLayout = ChannelLayout(ChannelFrontLeft | ChannelFrontRight | ChannelFrontCenter | ChannelBackLeft | ChannelBackRight | ChannelBackCenter)
	ChannelLayout6Point1         ChannelLayout = ChannelLayout(ChannelFrontLeft | ChannelFrontRight | ChannelFrontCenter | ChannelSideLeft | ChannelSideRight | ChannelLowFrequency | ChannelBackCenter)
	ChannelLayout6Point1Back     ChannelLayout = ChannelLayout(ChannelFrontLeft | ChannelFrontRight | ChannelFrontCenter | ChannelBackLeft | ChannelBackRight | ChannelLowFrequency | ChannelBackCenter)
	ChannelLayout6Point1Front    ChannelLayout = ChannelLayout(ChannelFrontLeft | ChannelFrontRight | ChannelSideLeft | ChannelSideRight | ChannelFrontLeftOfCenter | ChannelFrontRightOfCenter | ChannelLowFrequency)
	ChannelLayout7Point0         ChannelLayout = ChannelLayout(ChannelFrontLeft | ChannelFrontRight | ChannelFrontCenter | ChannelSideLeft | ChannelSideRight | ChannelBackLeft | ChannelBackRight)
	ChannelLayout7Point0Front    ChannelLayout = ChannelLayout(ChannelFrontLeft | ChannelFrontRight | ChannelFrontCenter | ChannelSideLeft | ChannelSideRight | ChannelFrontLeftOfCenter | ChannelFrontRightOfCenter)
	ChannelLayout7Point1         ChannelLayout = ChannelLayout(ChannelFrontLeft | ChannelFrontRight | ChannelFrontCenter | ChannelSideLeft | ChannelSideRight | ChannelLowFrequency | ChannelBackLeft | ChannelBackRight)
	ChannelLayout7Point1Wide     ChannelLayout = ChannelLayout(ChannelFrontLeft | ChannelFrontRight | ChannelFrontCenter | ChannelSideLeft | ChannelSideRight | ChannelLowFrequency | ChannelFrontLeftOfCenter | ChannelFrontRightOfCenter)
	ChannelLayout7Point1WideBack ChannelLayout = ChannelLayout(ChannelFrontLeft | ChannelFrontRight | ChannelFrontCenter | ChannelBackLeft | ChannelBackRight | ChannelLowFrequency | ChannelFrontLeftOfCenter | ChannelFrontRightOfCenter)
	ChannelLayoutOctagonal       ChannelLayout = ChannelLayout(ChannelFrontLeft | ChannelFrontRight | ChannelFrontCenter | ChannelSideLeft | ChannelSideRight | ChannelBackLeft | ChannelBackCenter | ChannelBackRight)
	ChannelLayoutHexadecagonal   ChannelLayout = ChannelLayout(ChannelFrontLeft | ChannelFrontRight | ChannelFrontCenter | ChannelSideLeft | ChannelSideRight | ChannelBackLeft | ChannelBackCenter | ChannelBackRight | ChannelWideLeft | ChannelWideRight | ChannelTopBackLeft | ChannelTopBackRight | ChannelTopBackCenter | ChannelTopFrontCenter | ChannelTopFrontLeft | ChannelTopFrontRight)
	ChannelLayoutStereoDownmix   ChannelLayout = ChannelLayout(ChannelStereoLeft | ChannelStereoRight)
)

func (l ChannelLayout) ctype() C.ulong {
	return C.ulong(l)
}
