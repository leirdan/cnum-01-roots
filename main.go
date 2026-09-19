package main

import (
	"fmt"
	"math"
)

type Interval []float32
type RootIntervalList []Interval

func deriv(x float32, f func(float32) float32) float32 {
	const h float32 = 1e-2
	return (f(x+h) - f(x)) / h
}

func nextInterval(a float32, b float32, f func(float32) float32, step float32) Interval {
	var slice Interval
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

func isolate(a float32, b float32, f func(float32) float32, step float32) RootIntervalList {
	var collection = RootIntervalList{}
	var curr float32 = a

	for curr < b {
		var sequence Interval = nextInterval(curr, b, f, step)
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
				collection = append(collection, Interval{xmin, xmax})
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
	var f = func(x float32) float32 {
		return float32(math.Pow(float64(x), 3)) - 9*x + 5
	}

	rootsInterval := isolate(-4, 10, f, 0.01)
	for k, r := range rootsInterval {
		fmt.Printf("Intervalo %d: [", k)
		for i, el := range r {
			fmt.Printf("%.2f", el)
			if i != len(r)-1 {
				fmt.Print(", ")
			}
		}
		fmt.Println("]")
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
