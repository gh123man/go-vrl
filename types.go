package main

//#include <stdio.h>
//#include <stdlib.h>
//#include <string.h>
//#include <vrl.h>
//#cgo LDFLAGS: -L${SRCDIR}/target/release -Wl,-rpath,${SRCDIR}/target/release -lvrl_bridge -lm -ldl
import "C"
import (
	"runtime"
	"unsafe"
)

type RustPointer struct {
	p unsafe.Pointer
}

func (rp *RustPointer) own(p unsafe.Pointer) {
	rp.p = p
	runtime.SetFinalizer(rp, free)
}

func free(rp *RustPointer) {
	C.free(unsafe.Pointer(rp.p))
}

type VrlError struct{ str string }

func (e VrlError) Error() string {
	return e.str
}
