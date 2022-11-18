
# Go VRL

Experimental go bindings for [Vector Remap Language](https://vector.dev/docs/reference/vrl/)

> Vector Remap Language (VRL) is an expression-oriented language designed for transforming observability data (logs and metrics) in a safe and performant manner. It features a simple syntax and a rich set of built-in functions tailored specifically to observability use cases.

## Usage

### Building and importing

Not quite ready yet. It's difficult to distribute a go module that depends on an external build system (However I am open to suggestions)

### Example

```go
program, err := Compile(`replace(string!(.), "go", "rust")`)
if err != nil {
    log.Panicln(err)
}

runtime := NewRuntime()
res, err := runtime.resolve(program, "hello go")
if err != nil {
    log.Panicln(err)
}

fmt.Println(res)
```

## What works

- Compiling VRL programs (and handling errors)
- Initializing the VRL runtime including:
  - `resolve` (run) the compled program
  - `clear` 
  - `is_empty`

## What doesn't work/missing bindings

- secrets
- metadata 
- timezone
- environment configuration 
- non-string typed inputs