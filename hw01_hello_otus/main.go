package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	phrase := "Hello, OTUS!"

	fmt.Println(reverse.String(phrase))
}
