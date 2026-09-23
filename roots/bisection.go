package roots

import (
	"cnum/types"
	"math"
)

// Função que executa o algoritmo da Bissecção para encontrar a raiz de uma função.
// O candidato inicial à raiz é calculado usando a função "b_candidate", que retorna o ponto médio do intervalo.
// Cada iteração calcula o novo candidato a raiz e ajusta o intervalo de acordo com a mudança de sinal entre raiz e extremos.
// A iteração avança enquanto o critério da precisão ou da distância dos extremos do intervalo não for cumprido.
//
// Input: Um Problem e uma taxa de precisão
//
// Output: Raiz encontrada e quantidade de iterações
func Bisection(inst types.Problem, epsilon float64) (float64, uint16) {
	a_ := inst.Start
	b_ := inst.End
	var root float64
	var counter uint16 = 1

	for root = b_candidate(a_, b_); math.Abs(inst.Function(root)) > epsilon && (b_-a_) > epsilon; root = b_candidate(a_, b_) {
		if inst.Function(a_)*inst.Function(root) < 0 {
			b_ = root
		} else if inst.Function(b_)*inst.Function(root) < 0 {
			a_ = root
		}
		counter++
	}

	return root, counter
}

func b_candidate(a, b float64) float64 {
	return (a + b) / 2
}
