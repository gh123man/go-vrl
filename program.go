package main

//#include <stdio.h>
//#include <stdlib.h>
//#include <string.h>
//#include <vrl.h>
//#cgo LDFLAGS: -L${SRCDIR}/target/release -Wl,-rpath,${SRCDIR}/target/release -lvrl_bridge -lm -ldl
import "C"
import "unsafe"

type Program = RustPointer

func Compile(str string) (*Program, error) {
	cs := C.CString(str)
	defer C.free(unsafe.Pointer(cs))
	result := C.compile_vrl(cs)
	if result.error != nil {
		defer C.free(unsafe.Pointer(result.error))
		return nil, VrlError{str: C.GoString(result.error)}
	}
	program := &Program{}
	program.own(result.value)
	return program, nil
}
