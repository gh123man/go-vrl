package main

import "fmt"

func main() {
	program, err := Compile(`
	. = replace(string!(.), r'\b\w{4}', "rust", 1)
	`)
	if err != nil {
		fmt.Println(err)
		return
	}

	runtime := NewRuntime()

	res, err := runtime.resolve(program, "hello world")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res)
}
