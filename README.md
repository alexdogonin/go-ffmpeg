# go-ffmpeg
Go binding to FFmpeg

# Installation

## Install ffpeg libraries

#### required versions  
libavcodec >= 58  
libavformat >= 58  
libavutil >= 56  
libswresample >= 3  
libswscale >= 5  

### Solus 
eopkg it pkg-config ffmpeg-devel

### Ubuntu 
apt install pkg-config libavcodec-dev libavformat-dev libavutil-dev libswscale-dev libswresample-dev

### Source code 
git clone https://github.com/FFmpeg/FFmpeg.git ffmpeg
cd ffmpeg
./configure
make
make install

# TODO
- [ ] add minimal tools number sufficient for this works:
  - [ ] encoding - generate video from image and audio track
  - [ ] decoding - generate number of images from video
  - [ ] rescaling video
  - [ ] convert video formats
- [ ] refactor code, reducing type dependencies among themselves
- [ ] restructure code to use go types instead FFmpeg types everywheare it possible (e.g. AVIOContext -> io.Writer)
- [ ] add comprehensive documentation
- [ ] (maybe never) rewrite code with pure go, rejection of cgo 
