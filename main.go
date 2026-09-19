package main

import (
	"cnum/roots"
	"cnum/types"
	"fmt"
	"math"
)

func main() {
	var instances []types.Problem
	instances = append(instances, types.Problem{
		Function: func(x float64) float64 {
			return 2 * math.Pow(float64(x), 4)
		}, Start: 0, End: 3, Id: 1, H: 0.6})
	instances = append(instances, types.Problem{
		Function: func(x float64) float64 {
			return math.Pow(x, 3) - 9*x + 5
		}, Start: -4, End: 4, Id: 1, H: 0.6})

	// TODO: organizar melhor a impressão
	rootsInterval := instances[len(instances)-1].IsolateRoots()
	for k, r := range rootsInterval {
		fmt.Printf("Intervalo %d: [", k)
		for i, el := range r {
			fmt.Printf("%.8f", el)
			if i != len(r)-1 {
				fmt.Print(", ")
			}
		}
		fmt.Println("]")
		result, counter := roots.Bisection(instances[len(instances)-1], 0.000001)
		fmt.Printf("[Bissecção] Raiz do intervalo: %.8f. Total de iterações: %d.\n", result, counter)
		result, counter = roots.FalsePosition(instances[len(instances)-1], 0.000001)
		fmt.Printf("[Falsa Posição] Raiz do intervalo: %.8f. Total de iterações: %d.\n", result, counter)
		result, counter = roots.Newton(instances[len(instances)-1], 0.000001, math.MaxUint16)
		fmt.Printf("[Newton] Raiz do intervalo: %.8f. Total de iterações: %d.\n", result, counter)
	}

}
