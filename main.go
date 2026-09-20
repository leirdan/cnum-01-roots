// package main

// import (
// 	"cnum/roots"
// 	"cnum/types"
// 	"fmt"
// 	"math"
// )

// func main() {
// 	var instances []types.Problem
// 	instances = append(instances, types.Problem{
// 		Function: func(x float64) float64 {
// 			return 2 * math.Pow(float64(x), 4)
// 		}, Start: 0, End: 3, Id: 1, H: 0.6})
// 	instances = append(instances, types.Problem{
// 		Function: func(x float64) float64 {
// 			return math.Pow(x, 3) - 9*x + 5
// 		}, Start: -4, End: 4, Id: 1, H: 0.6})

// 	// TODO: organizar melhor a impressão
// 	rootsInterval := instances[len(instances)-1].IsolateRoots()
// 	for k, r := range rootsInterval {
// 		fmt.Printf("Intervalo %d: [", k)
// 		for i, el := range r {
// 			fmt.Printf("%.8f", el)
// 			if i != len(r)-1 {
// 				fmt.Print(", ")
// 			}
// 		}
// 		fmt.Println("]")
// 		result, counter := roots.Bisection(instances[len(instances)-1], 0.000001)
// 		fmt.Printf("[Bissecção] Raiz do intervalo: %.8f. Total de iterações: %d.\n", result, counter)
// 		result, counter = roots.FalsePosition(instances[len(instances)-1], 0.000001)
// 		fmt.Printf("[Falsa Posição] Raiz do intervalo: %.8f. Total de iterações: %d.\n", result, counter)
// 		result, counter = roots.Newton(instances[len(instances)-1], 0.000001, math.MaxUint16)
// 		fmt.Printf("[Newton] Raiz do intervalo: %.8f. Total de iterações: %d.\n", result, counter)
// 	}

// }

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

	// Problema ID 1
	instances = append(instances, types.Problem{
		Function: func(x float64) float64 {
			return ((2 * math.Pow(float64(x), 4)) + (4 * math.Pow(float64(x), 3)) + (3 * math.Pow(float64(x), 2)) + (10 * x) - 15)
		}, Start: 0, End: 3, Id: 1, H: 0.6})

	// Problema ID 2
	instances = append(instances, types.Problem{
		Function: func(x float64) float64 {
			return math.Pow(x, 5) - (2*math.Pow(float64(x), 4) - (9*math.Pow(float64(x), 3) + (22*math.Pow(float64(x), 2) + (4 * x) - 24)))
		}, Start: 0, End: 5, Id: 2, H: 0.7})

	instances = append(instances, types.Problem{
		Function: func(x float64) float64 {
			return (math.Sin(x) * x) + 4
		}, Start: 1, End: 5, Id: 4, H: 0.5})

	// Declaramos a lista de métodos FORA do loop para não recriar a cada função
	methods := []types.NamedMethod{
		{Name: "Bissecção", Method: roots.Bisection},
		{Name: "Falsa Posição", Method: roots.FalsePosition},
		{
			Name: "Newton",
			Method: func(inst types.Problem, epsilon float64) (float64, uint16) {
				return roots.Newton(inst, epsilon, math.MaxUint16)
			},
		},
	}

	// Tipo auxiliar declarado fora do loop
	type ResultadoCorrida struct {
		Name string
		Root float64
		Iter uint16
	}

	// O loop principal itera sobre CADA função que você adicionou em 'instances'
	for idx, currentProblem := range instances {

		fmt.Printf("\n=========================================================================\n")
		fmt.Printf(" TESTANDO FUNÇÃO ID: %d (Índice no Array: %d)\n", currentProblem.Id, idx)
		fmt.Printf(" Intervalo Inicial: [%.2f, %.2f] | Passo (H): %.2f\n", currentProblem.Start, currentProblem.End, currentProblem.H)
		fmt.Printf("=========================================================================\n\n")

		rootsInterval := currentProblem.IsolateRoots()

		fmt.Println("Intervalos encontrados:")
		if len(rootsInterval) == 0 {
			fmt.Println("  Nenhum intervalo encontrado com raiz.")
			continue // Pula para a próxima função se não achou raízes
		}

		for k, r := range rootsInterval {
			fmt.Printf("  %d: [%.8f, %.8f]\n", k, r[0], r[len(r)-1])
		}
		fmt.Println("-------------------------------------------------------------------------")

		fmt.Println("Tempo Total de Execução por Método:")
		for _, m := range methods {
			var totalIteracoes uint16 = 0
			raizesEncontradas := []float64{}

			inicioMetodo := time.Now()

			for _, r := range rootsInterval {
				intervalProblem := currentProblem
				intervalProblem.Start = r[0]
				intervalProblem.End = r[len(r)-1]

				raiz, iteracoes := m.Method(intervalProblem, 0.000001)

				raizesEncontradas = append(raizesEncontradas, raiz)
				totalIteracoes += iteracoes
			}

			tempoMetodo := time.Since(inicioMetodo)

			fmt.Printf("- %-15s | Tempo: %10v | Iterações (soma): %3d\n", m.Name, tempoMetodo, totalIteracoes)
			fmt.Printf("  Raízes: ")
			for i, rz := range raizesEncontradas {
				fmt.Printf("%.8f", rz)
				if i < len(raizesEncontradas)-1 {
					fmt.Print(", ")
				}
			}
			fmt.Println("\n")
		}

		fmt.Println("-------------------------------------------------------------------------")
		fmt.Println("Corrida Paralela (Goroutines):")

		inicioCorrida := time.Now()
		vencedores := []ResultadoCorrida{}

		for _, r := range rootsInterval {
			intervalProblem := currentProblem
			intervalProblem.Start = r[0]
			intervalProblem.End = r[len(r)-1]

			winnerRoot, winnerCounter, winnerName := roots.ParallelRace(methods, intervalProblem, 0.000001)
			vencedores = append(vencedores, ResultadoCorrida{winnerName, winnerRoot, winnerCounter})
		}

		tempoCorrida := time.Since(inicioCorrida)

		fmt.Printf("- Tempo total da corrida para todos os intervalos: %v\n\n", tempoCorrida)

		for i, v := range vencedores {
			fmt.Printf("  Intervalo %d: Vencedor: %-13s | Raiz: %.8f | Iterações: %d\n", i, v.Name, v.Root, v.Iter)
		}
	}
}
