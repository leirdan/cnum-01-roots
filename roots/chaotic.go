package roots

import (
	"cnum/types"
	"math"
	"math/rand/v2"
)

func Chaotic(inst types.Problem, epsilon float64) (float64, uint16) {
	var sharedRoot float64 = (inst.Start + inst.End) / 2.0

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
				x := sharedRoot

				fx := inst.Function(x)
				fxh := inst.Function(x + h)
				deriv := (fxh - fx) / h

				if deriv != 0 {
					sharedRoot = x - (fx / deriv)
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
				xCurr := sharedRoot
				fxCurr := inst.Function(xCurr)
				fxPrev := inst.Function(xPrev)

				if fxCurr != fxPrev {
					next := (xPrev*fxCurr - xCurr*fxPrev) / (fxCurr - fxPrev)
					sharedRoot = next
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
				sharedRoot = inst.GFunction(sharedRoot)
			}
		}
	}()

	// Goroutine observadora
	go func() {
		var iterations uint16 = 0
		for {
			iterations++

			snapshot := sharedRoot

			if math.IsNaN(snapshot) || math.IsInf(snapshot, 0) {
				sharedRoot = inst.Start + rand.Float64()*(inst.End-inst.Start)
				continue
			}

			if math.Abs(inst.Function(snapshot)) < epsilon {
				chResult <- types.ResultStr{Root: snapshot, Iter: iterations}
				return
			}
		}
	}()

	final := <-chResult
	close(done)

	return final.Root, final.Iter
}
