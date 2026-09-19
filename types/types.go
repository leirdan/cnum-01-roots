package types

type Function func(float64) float64
type Interval []float64
type RootIntervalList []Interval

type Problem struct {
	Function Function
	Start    float64
	End      float64
	Id       uint8
	H        float64
}

type RootMethod func(inst Problem, epsilon float64) (float64, uint16)

type NamedMethod struct {
	Name   string
	Method RootMethod
}

type RaceResult struct {
	Root    float64
	Counter uint16
	Name    string
}

func (self Problem) IsolateRoots() RootIntervalList {
	return isolateRec(self.Start, self.End, self.Function, self.H)
}

func (self Problem) Deriv(x float64) float64 {
	return deriv(x, self.Function)
}

func isolateRec(a, b float64, f Function, step float64) RootIntervalList {
	var collection = RootIntervalList{}
	var curr float64 = a

	for curr < b {
		sequence := nextInterval(curr, b, f, step)
		curr = sequence[len(sequence)-1] // atualiza pro elemento que parou
		n := len(sequence)
		if n == 1 {
			continue
		}
		xmin := sequence[n-2]
		xmax := sequence[n-1]
		// tem pelo menos 1 raiz; se não tiver, ignora o intervalo
		if f(xmin)*f(xmax) < 0 {
			d1 := deriv(xmin, f)
			d2 := deriv(xmax, f)
			if d1*d2 >= 0 { // raiz única, boa
				collection = append(collection, Interval{xmin, xmax})
			} else {
				subintervals := isolateRec(xmin, xmax, f, step)
				collection = append(collection, subintervals...)
			}
		}
	}

	return collection
}

func deriv(x float64, f Function) float64 {
	const h float64 = 1e-2
	return (f(x+h) - f(x)) / h
}
func nextInterval(a, b float64, f Function, step float64) Interval {
	var slice Interval
	curr := a
	for curr < b && f(a)*f(curr) >= 0 {
		slice = append(slice, curr)
		curr += step
	}

	if curr <= b {
		slice = append(slice, curr)
	} else {
		slice = append(slice, b)
	}

	return slice
}
