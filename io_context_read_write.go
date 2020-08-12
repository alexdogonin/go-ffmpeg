package main

import "C"
import "unsafe"

//export Read
func Read(opaque unsafe.Pointer, buff unsafe.Pointer, size int) int {
	ind := *(*int)(opaque)
	b := (*[1024]byte)(buff)

	p := opaques[ind]

	p.r.Read((b)[:])
	return 0
}

//export Write
func Write(opaque unsafe.Pointer, buff unsafe.Pointer, size int) int { return 0 }
