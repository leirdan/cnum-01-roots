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
}

type ResultStr struct {
	Root float64
	Iter uint16
}

// Função da struct Problem que executa o isolamento das raízes, fase 1 do método numérico.
//
// Input: N/A
//
// Output: Lista dos intervalos da função, onde cada um tem somente 1 raiz
func (self Problem) IsolateRoots() RootIntervalList {
	return isolateRec(self.Start, self.End, self.Function, self.H)
}

// Função da struct Problem que calcula a derivada de forma analítica.
//
// Input: Um ponto x qualquer
//
// Output: f'(x), o valor da derivada no ponto
func (self Problem) Deriv(x float64) float64 {
	return deriv(x, self.Function)
}

// Função recursiva utilitária, aplica o Teorema de Bolzano e o corolário derivado
// para construir os menores intervalos possíveis com raízes. Em caso de cumprir com
// o teorema e corolário, o intervalo é adicionado à lista; caso contrário, ele é particionado novamente
//
// Input: início e final do intervalo, a função analisada e a distância (step) entre os pontos
//
// Output: Lista dos intervalos da função, onde cada um tem somente 1 raiz
func isolateRec(a, b float64, f Function, step float64) RootIntervalList {
	var collection = RootIntervalList{}
	var curr float64 = a

	for curr < b {
		sequence := nextInterval(curr, b, f, step)
		curr = sequence[len(sequence)-1]
		n := len(sequence)
		if n == 1 {
			continue
		}
		xmin := sequence[n-2]
		xmax := sequence[n-1]
		if f(xmin)*f(xmax) < 0 {
			d1 := deriv(xmin, f)
			d2 := deriv(xmax, f)
			if d1*d2 >= 0 {
				collection = append(collection, Interval{xmin, xmax})
			} else {
				subintervals := isolateRec(xmin, xmax, f, step)
				collection = append(collection, subintervals...)
			}
		}
	}

	return collection
}

// Função que calcula a derivada de forma analítica com h = 0.01
//
// Input: O ponto x pertencente ao domínio e a função a ser aplicada
//
// Output: O valor da derivada naquele ponto
func deriv(x float64, f Function) float64 {
	const h float64 = 1e-2
	return (f(x+h) - f(x)) / h
}

// Função que gera um intervalo com ao menos 1 raiz. A função caminha de um ponto a
// verificando a cada passo se houve mudança de sinal, montando o intervalo em caso positivo
// ou avançando caso contrário, até o limite máximo do ponto b.
//
// Input: início e final do intervalo, a função analisada e a distância (step) entre os pontos
//
// Output: um intervalo com raiz
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
