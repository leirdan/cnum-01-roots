package main

import (
	"cnum/roots"
	"fmt"
)

func deriv(x float32, f func(float32) float32) float32 {
	const h float32 = 10 ^ -2
	return (f(x+h) - f(x)) / h
}

func interval(a float32, b float32, f func(float32) float32) map[int8][]float32 {
	var collection = map[int8][]float32{}
	var slice_accum []float32
	var idx int8 = 1
	var curr float32 = a

	for curr < b {
		for ; f(a)*f(curr) >= 0; curr += 2 {
			slice_accum = append(slice_accum, curr)
			collection[idx] = slice_accum
			if curr >= b {
				break
			}
		}
		idx++
		a = curr
	}

	// TODO: usar corolário com derivada para verificar se há possibilidade de refinar ainda mais os intervalos?
	// for key, sequence := range collection {
	// 	n := len(sequence)
	// 	if n >= 3 {
	// 		idx1 := n / 3
	// 		idx2 := (2 * n) / 3
	// 		x1 :=
	// 	}
	// }

	return collection
}

func main() {
	fmt.Println("Hello, world!")
	root := roots.Bisection(0)
	fmt.Println(root)
}
