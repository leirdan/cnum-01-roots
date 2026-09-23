package roots

import (
	"cnum/types"
	"math"
)

// Função que executa o algoritmo da Falsa Posição para encontrar a raiz de uma função.
// O candidato inicial à raiz é calculado usando a função "fp_candidate", que realiza uma interpolação linear entre os extremos.
// Cada iteração calcula o novo candidato a raiz e ajusta o intervalo e os valores da função de acordo com a mudança de sinal.
// A iteração avança enquanto a precisão da função ou a distância entre as imagens dos extremos do intervalo não for cumprida.
//
// Input: Um Problem e uma taxa de precisão
//
// Output: Raiz encontrada e quantidade de iterações
func FalsePosition(inst types.Problem, epsilon float64) (float64, uint16) {
	a_ := inst.Start
	b_ := inst.End
	fa_ := inst.Function(a_)
	fb_ := inst.Function(b_)
	root := fp_candidate(a_, b_, fa_, fb_)
	froot := inst.Function(root)
	var counter uint16 = 1

	for math.Abs(fb_-fa_) > epsilon && math.Abs(froot) > epsilon {
		if fa_*froot < 0 {
			b_ = root
			fb_ = froot
		} else {
			a_ = root
			fa_ = froot
		}

		root = fp_candidate(a_, b_, fa_, fb_)
		froot = inst.Function(root)

		counter++
	}

	return root, counter
}

func fp_candidate(a, b, fa, fb float64) float64 {
	return (a*fb - b*fa) / (fb - fa)
}
