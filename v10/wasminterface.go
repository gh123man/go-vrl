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

//go:embed target/wasm32-wasi/release/vrl_bridge.wasm
var wasmBytes []byte

type WasmInterface struct {
	ctx     context.Context
	mod     api.Module
	runtime wazero.Runtime
}

type WasmString struct {
	ptr    uint64
	length uint64
}

func NewWasmInterface(ctx context.Context) *WasmInterface {
	// Create a new WebAssembly Runtime.
	r := wazero.NewRuntime(ctx)

	wasi_snapshot_preview1.MustInstantiate(ctx, r)
	mod, err := r.Instantiate(ctx, wasmBytes)
	if err != nil {
		log.Panicln(err)
	}

	// TODO check for expected exports, return error if they're not present

	return &WasmInterface{
		ctx: ctx, mod: mod, runtime: r,
	}
}

func (wr *WasmInterface) Compile(program string) (*Program, error) {
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

func (wr *WasmInterface) NewRuntime() (*Runtime, error) {
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

func (wr *WasmInterface) newWasmString(input string) *WasmString {
	inputSize := uint64(len(input))
	ptr := wr.allocate(inputSize)

	// The pointer is a linear memory offset, which is where we write the input string.
	if !wr.mod.Memory().Write(uint32(ptr), []byte(input)) {
		log.Panicf("Memory.Write(%d, %d) out of range of memory size %d",
			ptr, len(input), wr.mod.Memory().Size())
	}

	ws := WasmString{
		ptr:    ptr,
		length: uint64(len(input)),
	}

	runtime.SetFinalizer(&ws, func(ws *WasmString) { wr.deallocate(ws.ptr, ws.length) })

	return &ws
}

func (wr *WasmInterface) allocate(numBytes uint64) uint64 {
	allocate := wr.mod.ExportedFunction("allocate")

	results, err := allocate.Call(wr.ctx, numBytes)
	if err != nil {
		log.Panicln(err)
	}

	return results[0]
}

func (wr *WasmInterface) deallocate(ptr uint64, size uint64) {
	deallocate := wr.mod.ExportedFunction("deallocate")

	deallocate.Call(wr.ctx, ptr, size)
}

func (wr *WasmInterface) Close() {
	wr.runtime.Close(wr.ctx)
}

func unpackUInt64(val uint64) (uint32, uint32) {
	return uint32(val >> 32), uint32(val)
}
