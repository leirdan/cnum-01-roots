package roots

import (
	"cnum/types"
	"math"
	"math/rand/v2"
)

// Função que executa uma competição caótica entre os métodos de Newton, Secante e Ponto Fixo para encontrar a raiz de uma função.
// Três goroutines concorrentes compartilham e atualizam uma mesma variável de candidato sem sincronização,
// enquanto uma goroutine observadora monitora o estado atual.
// A goroutine observadora valida se a aproximação está dentro do intervalo e cumpre a precisão esperada,
// reiniciando o candidato aleatoriamente caso o valor divirja ou saia do intervalo.
// A execução é encerrada assim que a precisão for satisfeita ou o limite de iterações da goroutine observadora for atingido.
//
// Input: Um Problem, uma taxa de precisão e um limite máximo de iterações
//
// Output: Raiz encontrada e quantidade de iterações até a aceitação
func Chaotic(inst types.Problem, epsilon float64, kmax uint16) (float64, uint16) {
	sharedRoot := (inst.Start + inst.End) / 2.0

	read := func() float64 {
		return sharedRoot
	}
	write := func(v float64) {
		sharedRoot = v
	}

	done := make(chan struct{})
	chResult := make(chan types.ResultStr)

	// Newton goroutine
	go func() {
		h := 0.001
		for {
			select {
			case <-done:
				return
			default:
				x := read()

				fx := inst.Function(x)
				fxh := inst.Function(x + h)
				deriv := (fxh - fx) / h

				if deriv != 0 {
					write(x - (fx / deriv))
				}
			}
		}
	}()

	// Secant goroutine
	go func() {
		xPrev := inst.Start
		for {
			select {
			case <-done:
				return
			default:
				xCurr := read()
				fxCurr := inst.Function(xCurr)
				fxPrev := inst.Function(xPrev)

				if fxCurr != fxPrev {
					next := (xPrev*fxCurr - xCurr*fxPrev) / (fxCurr - fxPrev)
					write(next)
					xPrev = xCurr
				}
			}
		}
	}()

	// Fixed Point goroutine
	go func() {
		lambda := pf_lambda(inst)
		for {
			select {
			case <-done:
				return
			default:
				write(pf_candidate(read(), inst, lambda))
			}
		}
	}()

	// Watcher goroutine
	go func() {
		for i := uint16(1); i < kmax; i++ {
			snapshot := read()

			foraDoIntervalo := snapshot < inst.Start || snapshot > inst.End
			if math.IsNaN(snapshot) || math.IsInf(snapshot, 0) || foraDoIntervalo {
				write(inst.Start + rand.Float64()*(inst.End-inst.Start))
				continue
			}
			if math.Abs(inst.Function(snapshot)) < epsilon {
				chResult <- types.ResultStr{Root: snapshot, Iter: i}
				return
			}
		}
	}()

	final := <-chResult
	close(done)

	return final.Root, final.Iter
}
