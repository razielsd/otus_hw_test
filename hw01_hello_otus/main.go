package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	var message string = "Hello, OTUS!"
	fmt.Println(stringutil.Reverse(message))
}
