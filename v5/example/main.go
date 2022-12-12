package main

import (
	"fmt"
	"log"

	govrl "github.com/gh123man/go-vrl/v5"
)

func main() {
	bytesEnv()
	simpleDefault()
}

func bytesEnv() {
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
}

func simpleDefault() {
	program, err := govrl.Compile(`
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
		fmt.Println(err)
		return
	}

	runtime := govrl.NewRuntime()

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
