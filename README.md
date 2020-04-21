
# go-ffmpeg

Go binding to FFmpeg

## Installation

### Install ffpeg libraries

#### Solus

```bash
eopkg it pkg-config ffmpeg-devel
```

#### Ubuntu

```bash
apt install pkg-config libavcodec-dev libavformat-dev libavutil-dev libswscale-dev libswresample-dev
```

#### Source code

```bash
git clone https://github.com/FFmpeg/FFmpeg.git ffmpeg
cd ffmpeg
./configure
make
make install
```

#### Required versions

libavcodec >= 58  
libavformat >= 58  
libavutil >= 56  
libswresample >= 3  
libswscale >= 5  

## go-ffmpeg installation

```bash
go get github.com/alexdogonin/go-ffmpeg
```

## Usage

see [examples](https://github.com/alexdogonin/go-ffmpeg/tree/master/examples) directory

## TODO

- [ ] add minimal tools number sufficient for this works:
  - [ ] encoding - generate video from image and audio track
  - [ ] decoding - generate number of images from video
  - [ ] rescaling video
  - [ ] resampling audio
  - [ ] convert formats
- [ ] refactor code, reducing type dependencies among themselves
- [ ] restructure code to use go types instead FFmpeg types everywheare it possible (e.g. AVIOContext -> io.Writer)
- [ ] add comprehensive documentation
- [ ] (maybe never) rewrite code with pure go, rejection of cgo
