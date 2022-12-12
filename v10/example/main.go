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
