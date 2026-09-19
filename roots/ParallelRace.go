package roots

import (
	"cnum/types"
)

func ParallelRace(methods []types.NamedMethod, inst types.Problem, epsilon float64) (float64, uint16, string) {

	ch := make(chan types.RaceResult, len(methods))

	for _, m := range methods {
		go func(name string, method types.RootMethod) {
			r, c := method(inst, epsilon)
			ch <- types.RaceResult{
				Root:    r,
				Counter: c,
				Name:    name,
			}
		}(m.Name, m.Method)
	}

	winner := <-ch

	return winner.Root, winner.Counter, winner.Name
}
