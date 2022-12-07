
# Go VRL

Experimental go bindings for [Vector Remap Language](https://vector.dev/docs/reference/vrl/)

> Vector Remap Language (VRL) is an expression-oriented language designed for transforming observability data (logs and metrics) in a safe and performant manner. It features a simple syntax and a rich set of built-in functions tailored specifically to observability use cases.

## Usage

### Building and importing

Not quite ready yet. It's difficult to distribute a go module that depends on an external build system (However I am open to suggestions)

To use this repo as-is. `./run.sh` to build and run `main.go`

### Example

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
"hello rust"
```

[see `./example/main.go` for more examples](./example/main.go)

## What works

- Compiling VRL programs (and handling errors)
  - Supports bytes and object external environment kinds
- Initializing the VRL runtime including:
  - `resolve` (run) the compled program
  - `clear`
  - `is_empty`

## What doesn't work/missing bindings

- secrets
- metadata
- timezone
- environment configuration (partially implemented)
- most input types (other than bytes and object)
