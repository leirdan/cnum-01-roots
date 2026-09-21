package roots

import (
	"cnum/types"
	"math"
)

// FixedPoint encontra uma raiz de inst.Function pelo método do Ponto Fixo.
// Em vez de exigir uma função de iteração g(x) fornecida manualmente, constrói
// automaticamente g(x) = x - f(x)/λ, com λ estimado por pf_lambda a partir da
// própria f, garantindo convergência (|g'(x)| < 1) no intervalo isolado.
//
// Entrada:
//   - inst: problema já isolado, usando inst.Function, inst.Start e inst.End
//     (inst.Start/inst.End definem o intervalo [a,b] onde a raiz foi isolada)
//   - epsilon: precisão desejada, critério de parada |f(root)| <= epsilon
//   - kmax: número máximo de iterações permitido
//
// Saída: a raiz aproximada e a quantidade de iterações realizadas até
// satisfazer o critério de parada (ou até atingir kmax).
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

// pf_candidate calcula a próxima aproximação x_{k+1} = x_k - f(x_k)/λ da
// iteração de ponto fixo.
//
// Entrada: x (aproximação atual), inst (para acessar inst.Function) e lambda
// (constante estimada por pf_lambda).
//
// Saída: a próxima aproximação x_{k+1}.
func pf_candidate(x float64, inst types.Problem, lambda float64) float64 {
	return x - inst.Function(x)/lambda
}

// pf_lambda estima a constante λ usada para montar g(x) = x - f(x)/λ.
// Amostra a derivada de f em vários pontos do intervalo [Start, End] e toma o
// maior valor absoluto encontrado, com o mesmo sinal da derivada no ponto
// médio e uma margem de segurança de 5%. Isso garante |λ| > max|f'(x)| no
// intervalo, condição suficiente para que |g'(x)| = |1 - f'(x)/λ| < 1 e a
// iteração de ponto fixo convirja.
//
// Entrada: inst (usa inst.Start, inst.End e inst.Deriv).
//
// Saída: o valor estimado de λ.
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
