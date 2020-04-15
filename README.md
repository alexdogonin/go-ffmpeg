# go-ffmpeg
Go binding to FFmpeg

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
