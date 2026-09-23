package roots

import (
	"cnum/types"
	"math"
	"math/rand/v2"
)

// Função que executa o algoritmo de Newton-Raphson para encontrar a raiz de uma função.
// O candidato inicial à raiz é escolhido aleatoriamente dentro do intervalo.
// Cada iteração calcula o novo candidato usando a função "n_candidate", que utiliza a derivada da função no ponto atual.
// A iteração avança enquanto a precisão ou o limite máximo de iterações não for cumprido.
//
// Input: Um Problem, uma taxa de precisão e um limite máximo de iterações
//
// Output: Raiz encontrada e quantidade de iterações
func Newton(inst types.Problem, epsilon float64, kmax uint16) (float64, uint16) {
	var a = inst.Start
	var b = inst.End
	var root float64
	var counter uint16 = 1
	for root = a + rand.Float64()*(b-a); math.Abs(inst.Function(root)) > epsilon && counter < kmax; root = n_candidate(root, inst) {
		counter++
	}

	return root, counter
}

func n_candidate(x float64, inst types.Problem) float64 {
	return x - (inst.Function(x) / inst.Deriv(x))
}
