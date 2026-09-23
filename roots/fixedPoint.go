package roots

import (
	"cnum/types"
	"math"
)

// Função que executa o algoritmo do Ponto Fixo para encontrar a raiz de uma função.
// Estima a constante lambda via "pf_lambda" para construir uma função de iteração convergente g(x) = x - f(x)/λ.
// Cada iteração calcula o novo candidato a raiz usando a função "pf_candidate" a partir do candidato atual.
// A iteração avança enquanto a precisão ou o limite máximo de iterações não for cumprido.
//
// Input: Um Problem, uma taxa de precisão e um limite máximo de iterações
//
// Output: Raiz encontrada e quantidade de iterações
func FixedPoint(inst types.Problem, epsilon float64, kmax uint16) (float64, uint16) {
	lambda := pf_lambda(inst)
	root := b_candidate(inst.Start, inst.End)
	var counter uint16 = 1

	for math.Abs(inst.Function(root)) > epsilon && counter < kmax {
		root = pf_candidate(root, inst, lambda)
		counter++
	}

	return root, counter
}

// pf_candidate calcula o próximo candidato x_{k+1} para a iteração do ponto fixo
//
// Input: Candidato atual x, um Problem e a constante lambda
//
// Output: Próximo candidato a raiz
func pf_candidate(x float64, inst types.Problem, lambda float64) float64 {
	return x - inst.Function(x)/lambda
}

// pf_lambda estima a constante lambda para garantir a convergência da iteração.
// Calcula a derivada da função em múltiplos pontos do intervalo para determinar o maior valor absoluto com margem de segurança
//
// Input: Um Problem
//
// Output: O valor estimado de lambda
func pf_lambda(inst types.Problem) float64 {
	const samples = 20

	mid := b_candidate(inst.Start, inst.End)
	sign := math.Copysign(1, inst.Deriv(mid))
	maxAbs := math.Abs(inst.Deriv(inst.Start))
	step := (inst.End - inst.Start) / samples

	for i := 1; i <= samples; i++ {
		x := inst.Start + float64(i)*step
		d := math.Abs(inst.Deriv(x))
		if d > maxAbs {
			maxAbs = d
		}
	}

	const margem = 1.05
	return sign * maxAbs * margem
}
