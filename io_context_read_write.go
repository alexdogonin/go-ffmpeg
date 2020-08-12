package main

import "C"
import "unsafe"

//export Read
func Read(opaque unsafe.Pointer, buff unsafe.Pointer, size int) int {
	p := (*ioContextOpaque)(opaque)
	_ = p
	return 0
}

//export Write
func Write(opaque unsafe.Pointer, buff unsafe.Pointer, size int) int { return 0 }
