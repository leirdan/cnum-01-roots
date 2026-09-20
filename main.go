package main

import (
	"cnum/roots"
	"cnum/types"
	"fmt"
	"math"
	"time"
)

func main() {
	var instances []types.Problem

	const PRECISION float64 = 0.00000001

	// f(x) = 2x^4 + 4x^3 + 3x^2 + 10x - 15
	instances = append(instances, types.Problem{
		Function: func(x float64) float64 {
			return ((2 * math.Pow(float64(x), 4)) + (4 * math.Pow(float64(x), 3)) + (3 * math.Pow(float64(x), 2)) + (10 * x) - 15)
		},
		GFunction: func(x float64) float64 {
			return (15 - (2 * math.Pow(x, 4)) - (4 * math.Pow(x, 3)) - (3 * math.Pow(x, 2))) / 10.0
		},
		Start: 0, End: 3, Id: 1, H: 0.6})

	// f(x) = x^5 - 2x^4 + 9x^3 - 22x^2 - 4x + 24
	instances = append(instances, types.Problem{
		Function: func(x float64) float64 {
			return math.Pow(x, 5) - (2*math.Pow(float64(x), 4) - (9*math.Pow(float64(x), 3) + (22*math.Pow(float64(x), 2) + (4 * x) - 24)))
		},
		GFunction: func(x float64) float64 {
			return (math.Pow(x, 5) - 2*math.Pow(x, 4) + 9*math.Pow(x, 3) - 22*math.Pow(x, 2) + 24) / 4.0
		},
		Start: 0, End: 5, Id: 2, H: 0.7})

	// f(x) = 5x^3 + x^2 - e^(1-2x) + cos(x) + 20
	instances = append(instances, types.Problem{
		Function: func(x float64) float64 {
			return 5*math.Pow(x, 3) + math.Pow(x, 2) - math.Exp(1-2*x) + math.Cos(x) + 20
		},
		GFunction: func(x float64) float64 {
			return math.Sqrt(math.Exp(1-2*x) - 5*math.Pow(x, 3) - math.Cos(x) - 20)
		},
		Start: -5, End: 5, Id: 3, H: 0.5})

	// f(x) = x*sin(x) + 4
	instances = append(instances, types.Problem{
		Function: func(x float64) float64 {
			return (math.Sin(x) * x) + 4
		},
		GFunction: func(x float64) float64 {
			sinX := math.Sin(x)
			if sinX == 0 {
				return 0.1
			}
			return -4.0 / sinX
		},
		Start: 1, End: 5, Id: 4, H: 0.5})

	methods := []types.NamedMethod{
		{Name: "Bissecção", Method: roots.Bisection},
		{Name: "Falsa Posição", Method: roots.FalsePosition},
		{
			Name: "Newton",
			Method: func(inst types.Problem, epsilon float64) (float64, uint16) {
				return roots.Newton(inst, epsilon, math.MaxUint16)
			},
		},
		{
			Name: "Secante", Method: func(inst types.Problem, epsilon float64) (float64, uint16) {
				return roots.Secant(inst, epsilon, math.MaxUint16)
			},
		},
	}

	for _, currentProblem := range instances {
		fmt.Printf("\n=========================================================================\n")
		fmt.Printf(" FUNÇÃO ID: %d \n", currentProblem.Id)
		fmt.Printf(" Intervalo Inicial: [%.2f, %.2f] | H: %.2f\n", currentProblem.Start, currentProblem.End, currentProblem.H)
		fmt.Printf("=========================================================================\n\n")

		rootsInterval := currentProblem.IsolateRoots()
		fmt.Println("Intervalos com raízes encontrados:")
		if len(rootsInterval) == 0 {
			fmt.Println("  Nenhum intervalo encontrado com raiz.")
			continue
		}

		for k, r := range rootsInterval {
			fmt.Printf("  %d: [%.12f, %.12f]\n", k, r[0], r[len(r)-1])
		}
		fmt.Println("-------------------------------------------------------------------------")

		fmt.Println("[EXECUÇÕES TRADICIONAIS]")
		for _, m := range methods {
			var iterations uint16 = 0
			foundRoots := []float64{}

			startTime := time.Now()

			for _, r := range rootsInterval {
				intervalProblem := currentProblem
				intervalProblem.Start = r[0]
				intervalProblem.End = r[len(r)-1]

				root, its := m.Method(intervalProblem, PRECISION)

				foundRoots = append(foundRoots, root)
				iterations += its
			}

			time := time.Since(startTime)

			fmt.Printf("- %-15s | Tempo: %10v | Iterações (soma): %3d\n", m.Name, time, iterations)
			fmt.Printf("  Raízes: ")
			for i, rz := range foundRoots {
				fmt.Printf("%.12f", rz)
				if i < len(foundRoots)-1 {
					fmt.Print(", ")
				}
			}
			fmt.Println("")
		}

		fmt.Println("-------------------------------------------------------------------------")
		fmt.Println("[CORRIDA MALUCA]")

		start := time.Now()
		winners := []types.RaceResult{}

		for _, r := range rootsInterval {
			intervalProblem := currentProblem
			intervalProblem.Start = r[0]
			intervalProblem.End = r[len(r)-1]

			winnerRoot, winnerCounter := roots.Chaotic(intervalProblem, PRECISION)
			winners = append(winners, types.RaceResult{Root: winnerRoot, Counter: winnerCounter})
		}

		time := time.Since(start)
		fmt.Printf("- Tempo Total de Execução: %10v \n", time)
		for i, v := range winners {
			fmt.Printf("- Intervalo %d: Raiz: %.12f | Iterações: %d\n", i, v.Root, v.Counter)
		}
	}
}
