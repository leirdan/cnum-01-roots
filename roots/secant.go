package roots

import (
	"cnum/types"
	"math"
)

// Função que executa o algoritmo da Secante para encontrar a raiz de uma função.
// Dois chutes iniciais são calculados usando o ponto médio e a falsa posição dentro do intervalo do problema.
// Cada iteração calcula o novo candidato usando a função "s_candidate", que aproxima a derivada por uma secante.
// A iteração avança enquanto o critério da precisão da função, distância entre os pontos ou o limite máximo de iterações não for cumprido.
//
// Input: Um Problem, uma taxa de precisão e um limite máximo de iterações
//
// Output: Raiz encontrada e quantidade de iterações
func Secant(inst types.Problem, epsilon float64, kmax uint16) (float64, uint16) {
	x1 := b_candidate(inst.Start, inst.End)
	x2 := fp_candidate(inst.Start, inst.End, inst.Function(inst.Start), inst.Function(inst.End))
	fx1 := inst.Function(x1)
	fx2 := inst.Function(x2)
	var counter uint16 = 0

	for ; math.Abs(fx2) > epsilon && math.Abs(x2-x1) > epsilon && counter < kmax; counter++ {
		x_ := s_candidate(x1, x2, fx1, fx2)
		x1 = x2
		fx1 = fx2
		x2 = x_
		fx2 = inst.Function(x2)
	}

	return x2, counter
}

func s_candidate(x1, x2, fx1, fx2 float64) float64 {
	return x2 - (fx2 * (x2 - x1) / (fx2 - fx1))
}
