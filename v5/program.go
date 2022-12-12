package govrl

//#include <stdio.h>
//#include <stdlib.h>
//#include <string.h>
//#include <vrl.h>
//#cgo LDFLAGS: -lvrl_bridge -lm
//#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/deps/darwin_x86_64/release -framework CoreFoundation
//#cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/deps/darwin_arm64/release -framework CoreFoundation
//#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/deps/linux_x86_64/release -ldl
//#cgo linux,arm64 LDFLAGS: -L${SRCDIR}/deps/linux_arm64/release -ldl
import "C"
import "unsafe"

type Program = RustPointer
type ExternalEnv = RustPointer

type Kind int

const (
	Bytes Kind = iota
	Object
)

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

func CompileWithExternal(str string, env *ExternalEnv) (*Program, error) {
	cs := C.CString(str)
	defer C.free(unsafe.Pointer(cs))
	result := C.compile_vrl_with_external(cs, env.p)
	if result.error != nil {
		defer C.free(unsafe.Pointer(result.error))
		return nil, VrlError{str: C.GoString(result.error)}
	}
	program := &Program{}
	program.own(result.value)
	return program, nil
}

func GetDefaultExternalEnv() *ExternalEnv {
	e := &ExternalEnv{}
	e.own(C.external_env_default())
	return e
}

func GetExternalEnv(target Kind, metadata Kind) *ExternalEnv {
	e := &ExternalEnv{}

	var targetKind unsafe.Pointer
	var metadataKind unsafe.Pointer

	switch target {
	case Bytes:
		targetKind = C.kind_bytes()
	case Object:
		targetKind = C.kind_object()
	}

	switch metadata {
	case Bytes:
		metadataKind = C.kind_bytes()
	case Object:
		targetKind = C.kind_object()
	}

	e.own(C.external_env(targetKind, metadataKind))
	return e
}
