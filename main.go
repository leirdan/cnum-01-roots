package main

import (
	"cnum/roots"
	"fmt"
)

func main() {
	fmt.Println("Hello, world!")
	root := roots.Bisection(0)
	fmt.Println(root)
}
