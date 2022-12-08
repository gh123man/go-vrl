package govrl

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"runtime"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

func unpackUInt64(val uint64) (uint32, uint32) {
	return uint32(val >> 32), uint32(val)
}

type Program struct {
	ptr  uint32
	wasm *VrlInterface
}

type Runtime struct {
	ptr  uint32
	wasm *VrlInterface
}

func (r *Runtime) Resolve(program *Program, input string) (string, error) {
	runtimeResolveFunc := r.wasm.mod.ExportedFunction("runtime_resolve")

	inputWasmString := r.wasm.newWasmString(input)

	results, err := runtimeResolveFunc.Call(r.wasm.ctx, uint64(r.ptr), uint64(program.ptr), uint64(inputWasmString.ptr), uint64(inputWasmString.length))
	if err != nil {
		return "", nil
	}

	resultPtr, resultLength := unpackUInt64(results[0])
	resultStringBytes, ok := r.wasm.mod.Memory().Read(r.wasm.ctx, resultPtr, resultLength)
	if !ok {
		log.Panicf("Memory.Read(%d, %d) out of range of memory size %d",
			resultPtr, resultLength, r.wasm.mod.Memory().Size(r.wasm.ctx))
	}
	res := string(resultStringBytes)
	return res, nil
}

func (r *Runtime) Clear() {
	// TODO
}

type VrlInterface struct {
	ctx     context.Context
	mod     api.Module
	runtime wazero.Runtime
}

type WasmString struct {
	ptr    uint64
	length uint64
}

//go:embed target/wasm32-wasi/release/vrl_bridge.wasm
var wasmBytes []byte

func NewVrlInterface(ctx context.Context) *VrlInterface {
	// Create a new WebAssembly Runtime.
	r := wazero.NewRuntime(ctx)

	wasi_snapshot_preview1.MustInstantiate(ctx, r)
	mod, err := r.InstantiateModuleFromBinary(ctx, wasmBytes)
	if err != nil {
		log.Panicln(err)
	}

	// TODO check for expected exports, return error if they're not present

	return &VrlInterface{
		ctx: ctx, mod: mod, runtime: r,
	}
}

func (wr *VrlInterface) Compile(program string) (*Program, error) {
	compileVrlFunc := wr.mod.ExportedFunction("compile_vrl")

	ws := wr.newWasmString(program)

	results, err := compileVrlFunc.Call(wr.ctx, ws.ptr, ws.length)
	if err != nil {
		return nil, err
	}

	programPtr := uint32(results[0])
	if programPtr == 0 {
		return nil, fmt.Errorf("Unknown error from compile_vrl rust-side")
	}

	return &Program{programPtr, wr}, nil
}

func (wr *VrlInterface) NewRuntime() (*Runtime, error) {
	newRuntimeFunc := wr.mod.ExportedFunction("new_runtime")

	results, err := newRuntimeFunc.Call(wr.ctx)

	if err != nil {
		return nil, nil
	}

	runtimePtr := uint32(results[0])
	if runtimePtr == 0 {
		return nil, nil
	}

	return &Runtime{runtimePtr, wr}, nil
}

// Helpers

func (wr *VrlInterface) newWasmString(input string) *WasmString {
	inputSize := uint64(len(input))
	ptr := wr.allocate(inputSize)

	// The pointer is a linear memory offset, which is where we write the input string.
	if !wr.mod.Memory().Write(wr.ctx, uint32(ptr), []byte(input)) {
		log.Panicf("Memory.Write(%d, %d) out of range of memory size %d",
			ptr, len(input), wr.mod.Memory().Size(wr.ctx))
	}

	ws := WasmString{
		ptr:    ptr,
		length: uint64(len(input)),
	}

	runtime.SetFinalizer(&ws, func(ws *WasmString) { wr.deallocate(ws.ptr, ws.length) })

	return &ws
}

func (wr *VrlInterface) allocate(numBytes uint64) uint64 {
	allocate := wr.mod.ExportedFunction("allocate")

	results, err := allocate.Call(wr.ctx, numBytes)
	if err != nil {
		log.Panicln(err)
	}

	return results[0]
}

func (wr *VrlInterface) deallocate(ptr uint64, size uint64) {
	deallocate := wr.mod.ExportedFunction("deallocate")

	deallocate.Call(wr.ctx, ptr, size)
}

func (wr *VrlInterface) Close() {
	wr.runtime.Close(wr.ctx)
}
