package main

import (
	"cnum/roots"
	"cnum/types"
	"fmt"
	"math"
)

func deriv(x float32, f types.Function) float32 {
	const h float32 = 1e-2
	return (f(x+h) - f(x)) / h
}

func nextInterval(a, b float32, f types.Function, step float32) types.Interval {
	var slice types.Interval
	curr := a
	for curr < b && f(a)*f(curr) >= 0 {
		slice = append(slice, curr)
		curr += step
	}

	if curr <= b {
		slice = append(slice, curr)
	} else {
		slice = append(slice, b)
	}

	return slice
}

func isolate(a, b float32, f types.Function, step float32) types.RootIntervalList {
	var collection = types.RootIntervalList{}
	var curr float32 = a

	for curr < b {
		sequence := nextInterval(curr, b, f, step)
		curr = sequence[len(sequence)-1] // atualiza pro elemento que parou
		n := len(sequence)
		if n == 1 {
			continue
		}
		xmin := sequence[n-2]
		xmax := sequence[n-1]
		// tem pelo menos 1 raiz; se não tiver, ignora o intervalo
		if f(xmin)*f(xmax) < 0 {
			d1 := deriv(xmin, f)
			d2 := deriv(xmax, f)
			if d1*d2 >= 0 { // raiz única, boa
				collection = append(collection, types.Interval{xmin, xmax})
			} else {
				subintervals := isolate(xmin, xmax, f, step)
				collection = append(collection, subintervals...)
			}
		}
	}

	return collection
}

func main() {
	fmt.Println("Hello, world!")
	var f types.Function = func(x float32) float32 {
		return float32(math.Pow(float64(x), 3)) - 9*x + 5
	}

	rootsInterval := isolate(-4, 10, f, 0.6)
	for k, r := range rootsInterval {
		fmt.Printf("Intervalo %d: [", k)
		for i, el := range r {
			fmt.Printf("%.8f", el)
			if i != len(r)-1 {
				fmt.Print(", ")
			}
		}
		fmt.Println("]")
		result, counter := roots.Bisection(r[0], r[len(r)-1], 0.6, 0.000001, f)
		fmt.Printf("Raiz do intervalo: %.8f. Total de iterações: %d.\n", result, counter)
	}

}

/*
 f(-4) = -23
 f(-3) = 5
 f(-2) = 15
 f(-1) = 13
 f(0) = 5
 f(1) = -3
 f(2) = -5
 f(3) = 5
 f(4) = 33

 Intervalos: [-4, -3], [0, 1], [2, 3]
*/
