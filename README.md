
# Go VRL

Experimental go bindings for [Vector Remap Language](https://vector.dev/docs/reference/vrl/)

> Vector Remap Language (VRL) is an expression-oriented language designed for transforming observability data (logs and metrics) in a safe and performant manner. It features a simple syntax and a rich set of built-in functions tailored specifically to observability use cases.

## Versions
There are two major versions of this module and consumers must choose which is a
better fit for their use case.

They aim to support as similar an interface as possible, with the key
distinction being how VRL programs are executed.

- **V5** uses `cgo` to interface with a custom library built from VRL. This has
  better performance with the main downside being that it relies on `cgo`, which
  some applications may not care for.
- **V10** uses `wasm` to execute VRL. It performs worse, on the order of 2-3 times
  slower, however VRL is quite efficient so this still offers relatively good
  absolute performance.

## Usage

### Feature Support

|                           | V5                  | V10 |
|-------------------------- | ------------------- | --- |
| Compiling a VRL Program   | ✅                  | ✅  |
| Running a VRL Program     | ✅                  | ✅  |
| VRL Runtime "Basic"\* API | ✅                  | ✅  |
| Environment Kinds         | 'Byte' and 'Object' | ❌  |
| Secrets                   | ❌                  | ❌  |
| Metadata                  | ❌                  | ❌  |
| Timezones                 | ❌                  | ❌  |
| Requires CGO              | ✅                  | ❌  |
| Full VRL stdlib support   | ✅                  | ❌\* |


\* "Basic" API currently means:
- `compile`
- `resolve` (run) the compiled program
- `clear`
- `is_empty`

\* WASM supports almost most of VRL's stdlib functions, the unsupported ones can
be found [with this GH issues
search](https://github.com/vectordotdev/vector/issues?q=is%3Aopen+is%3Aissue+label%3A%22vrl%3A+playground%22+wasm+compatible)

### Building and importing

Not quite ready yet. It's difficult to distribute a go module that depends on an external build system, we have some ideas though.

To use this repo as-is, its required to manually compile the rust dependency.
For V5: `cd v5; cargo build --release; cd example/; go run .`
For V10: `cd v10; cargo build --target wasm32-wasi --release; cd example/; go run .`

### Examples

#### V5

```go
program, err := govrl.CompileWithExternal(`replace(., "go", "rust")`, govrl.GetExternalEnv(govrl.Bytes, govrl.Bytes))
if err != nil {
    log.Panicln(err)
}

runtime := govrl.NewRuntime()
res, err := runtime.Resolve(program, "hello go")
if err != nil {
    log.Panicln(err)
}

fmt.Println(res)
```

```bash
$ go run .
hello rust
```

[see `./v5/example/main.go` for more examples](./v5/example/main.go)

#### V10

```go
package main

import (
	"context"
	"fmt"
	"log"

	govrl "github.com/gh123man/go-vrl/v10"
)

func main() {
	simpleDefault()
}

func simpleDefault() {
	ctx := context.Background()
	wasmInterface := govrl.NewWasmInterface(ctx)
	program, err := wasmInterface.Compile(`
	. = parse_json!(string!(.))
	del(.foo)

	.timestamp = now()

	http_status_code = parse_int!(.http_status)
	del(.http_status)

	if http_status_code >= 200 && http_status_code <= 299 {
		.status = "success"
	} else {
		.status = "error"
	}
	.
	`)

	if err != nil {
		log.Panicln(err)
		return
	}

	runtime, err := wasmInterface.NewRuntime()
	if err != nil {
		log.Panicln(err)
	}

	res, err := runtime.Resolve(program, `{
		"message": "Hello VRL",
		"foo": "delete me",
		"http_status": "200"
	}
	`)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res)
	runtime.Clear()
}
```

```bash
$ go run .
{ "message": "Hello VRL", "status": "success", "timestamp": t'2022-01-01T00:00:00Z' }
```

[see `./v10/example/main.go` for more examples](./v10/example/main.go)

