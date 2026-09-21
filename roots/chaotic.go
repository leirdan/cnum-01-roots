package roots

import (
	"cnum/types"
	"math"
	"math/rand/v2"
	"sync"
)

// Chaotic implementa a estratégia proposta no item 1.3: uma "corrida caótica"
// onde Newton, Secante e Ponto Fixo disputam a mesma raiz aproximada,
// guardada em uma única variável compartilhada (sharedRoot). Cada método
// lê o valor mais recente escrito por qualquer um dos outros e calcula sua
// própria próxima aproximação a partir dele — não roda isoladamente com sua
// própria sequência, mas "surfa" no progresso feito pelos demais.
//
// sharedRoot é protegida por um sync.Mutex (via read/write) porque as quatro
// goroutines abaixo (Newton, Secante, Ponto Fixo e a observadora) leem e
// escrevem nela concorrentemente; sem essa exclusão mútua o acesso é uma
// data race real (confirmado com `go run -race`), podendo corromper o valor
// lido/escrito.
//
// A goroutine observadora também garante que o valor aceito como resultado
// esteja dentro do intervalo isolado [inst.Start, inst.End]: como os métodos
// escrevem na mesma variável sem coordenação, uma iteração de Newton ou do
// Ponto Fixo pode eventualmente "empurrar" sharedRoot para fora do intervalo
// isolado e convergir para outra raiz de f, fora do subintervalo que se
// pediu para refinar. Se isso acontecer, o valor é descartado e a busca
// reinicia dentro do intervalo, em vez de aceitar uma raiz de fora dele.
//
// Entrada:
//   - inst: problema já isolado (usa inst.Start/inst.End como limites do
//     intervalo válido, inst.Function e, se disponível, inst.GFunction para
//     a goroutine do Ponto Fixo)
//   - epsilon: precisão desejada, critério de parada |f(root)| < epsilon
//   - kmax: número máximo de iterações da goroutine observadora
//
// Saída: a raiz aproximada (sempre dentro de [inst.Start, inst.End]) e o
// número de iterações da goroutine observadora até aceitá-la.
func Chaotic(inst types.Problem, epsilon float64, kmax uint16) (float64, uint16) {
	var mu sync.Mutex
	sharedRoot := (inst.Start + inst.End) / 2.0

	read := func() float64 {
		mu.Lock()
		defer mu.Unlock()
		return sharedRoot
	}
	write := func(v float64) {
		mu.Lock()
		sharedRoot = v
		mu.Unlock()
	}

	done := make(chan struct{})
	chResult := make(chan types.ResultStr)

	// Goroutine do Newton
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

	// Goroutine do Secante
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

	// Goroutine do Ponto Fixo
	go func() {
		if inst.GFunction == nil {
			return
		}
		for {
			select {
			case <-done:
				return
			default:
				write(inst.GFunction(read()))
			}
		}
	}()

	// Goroutine observadora
	go func() {
		for i := uint16(1); i < kmax; i++ {
			snapshot := read()

			foraDoIntervalo := snapshot < inst.Start || snapshot > inst.End
			if math.IsNaN(snapshot) || math.IsInf(snapshot, 0) || foraDoIntervalo {
				// alguma goroutine divergiu ou "fugiu" do intervalo isolado;
				// descarta e reinicia a busca dentro de [Start, End]
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
