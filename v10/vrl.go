package govrl

import (
	_ "embed"
	"log"
)

type Program struct {
	ptr  uint32
	wasm *WasmInterface
}

type Runtime struct {
	ptr  uint32
	wasm *WasmInterface
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

func (r *Runtime) Clear() error {
	runtimeClearFunc := r.wasm.mod.ExportedFunction("runtime_clear")

	results, err := runtimeClearFunc.Call(r.wasm.ctx, uint64(r.ptr))
	return err
}

func (r *Runtime) IsEmpty() bool {
	runtimeIsEmptyFunc := r.wasm.mod.ExportedFunction("runtime_is_empty")

	results, err := runtimeIsEmptyFunc.Call(r.wasm.ctx, uint64(r.ptr))
	return err
}
