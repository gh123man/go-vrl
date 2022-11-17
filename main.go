package main

import "fmt"

func main() {
	program, err := Compile(`
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

	runtime := NewRuntime()

	res, err := runtime.resolve(program, `{
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

	runtime.clear()
}
